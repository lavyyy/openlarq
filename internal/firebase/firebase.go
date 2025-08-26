package firebase

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type FirebaseClient struct {
	conn          *websocket.Conn
	host          string
	namespace     string
	idToken       string
	requestID     int
	pending       map[int]chan QueryBody
	streamBuffer  []QueryData
	messageMutex  sync.Mutex
	userId        string
	autoReconnect bool
	stopChan      chan struct{}
}

func NewFirebaseClient(projectID, databaseURL string) (*FirebaseClient, error) {
	host := strings.TrimPrefix(strings.TrimPrefix(databaseURL, "https://"), "http://")
	host = strings.TrimSuffix(host, "/")
	host = strings.Split(host, "/")[0]

	namespace := projectID

	client := &FirebaseClient{
		host:          host,
		namespace:     namespace,
		requestID:     1,
		pending:       make(map[int]chan QueryBody),
		streamBuffer:  make([]QueryData, 0),
		autoReconnect: true,
		stopChan:      make(chan struct{}),
	}

	// connect immediately
	if err := client.Connect(); err != nil {
		return nil, err
	}

	return client, nil
}

func (c *FirebaseClient) Connect() error {
	u := url.URL{
		Scheme:   "wss",
		Host:     c.host,
		Path:     "/.ws",
		RawQuery: fmt.Sprintf("v=5&ns=%s", c.namespace),
	}

	log.Printf("Connecting to Firebase WebSocket: %s", u.String())
	log.Printf("Host: %s, Namespace: %s", c.host, c.namespace)

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Printf("Failed to connect to Firebase WebSocket: %v", err)
		return err
	}
	c.conn = conn

	log.Printf("Successfully connected to Firebase WebSocket")

	// start listening for messages
	go c.Listen()

	return nil
}

func (c *FirebaseClient) AuthenticateUser(idToken string) error {
	c.idToken = idToken

	log.Printf("Authenticating with Firebase using ID token (length: %d)", len(idToken))

	// make auth request
	authRes, err := c.sendRequest("auth", QueryBody{
		Cred: &c.idToken,
	})
	if err != nil {
		log.Printf("Failed to authenticate with Firebase: %v", err)
		return err
	}

	log.Printf("Authentication successful, response: %+v", authRes)

	// store user ID from auth response
	if data, ok := authRes.Data.(map[string]interface{}); ok {
		if auth, ok := data["auth"].(map[string]interface{}); ok {
			if userID, ok := auth["user_id"].(string); ok {
				c.userId = userID
				log.Printf("Extracted user ID: %s", userID)
			} else {
				log.Printf("user_id not found in auth response")
			}
		} else {
			log.Printf("auth field not found in response")
		}
	} else {
		log.Printf("No data field in auth response")
	}

	return nil
}

func (c *FirebaseClient) sendRequest(action string, requestData QueryBody) (QueryBody, error) {
	requestID := c.next()

	// create a channel to receive the response
	responseChan := make(chan QueryBody, 1)

	// register the response channel
	c.messageMutex.Lock()
	c.pending[requestID] = responseChan
	c.messageMutex.Unlock()

	msg := QueryMessage{
		Type: "d",
		Data: QueryData{
			RequestId: requestID,
			Action:    action,
			Body:      requestData,
		},
	}

	data, _ := json.Marshal(msg)

	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("Failed to write message to WebSocket: %v", err)
		c.messageMutex.Lock()
		delete(c.pending, requestID)
		c.messageMutex.Unlock()
		return QueryBody{}, fmt.Errorf("send error: %v", err)
	}

	// wait for the response
	select {
	case response := <-responseChan:
		return response, nil
	case <-time.After(30 * time.Second):
		log.Printf("Timeout waiting for response to request %d", requestID)
		c.messageMutex.Lock()
		delete(c.pending, requestID)
		c.messageMutex.Unlock()
		return QueryBody{}, fmt.Errorf("timeout waiting for response")
	}
}

func (c *FirebaseClient) processMessage(text string) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(text), &data); err != nil {
		log.Println("json decode error:", err)
		return
	}

	c.handleMessage(data)
}

func (c *FirebaseClient) Listen() {
	var chunkCount int
	var buffer []string

	for {
		select {
		case <-c.stopChan:
			return
		default:
			_, msg, err := c.conn.ReadMessage()
			if err != nil {
				log.Printf("WebSocket read error: %v", err)

				// Attempt to reconnect if enabled
				if c.autoReconnect {
					log.Println("Attempting to reconnect...")
					if err := c.reconnect(); err != nil {
						log.Printf("Reconnection failed: %v", err)
						time.Sleep(5 * time.Second) // Wait before next attempt
						continue
					}
					log.Println("Successfully reconnected")
					continue
				}

				return
			}

			text := string(msg)

			// case 1: first message is a number -> indicates chunk count
			if chunkCount == 0 {
				if n, err := strconv.Atoi(text); err == nil {
					chunkCount = n
					buffer = make([]string, 0, chunkCount)
					continue
				}
			}

			// case 2: in the middle of chunked message
			if chunkCount > 0 {
				buffer = append(buffer, text)
				if len(buffer) == chunkCount {
					// join and decode
					joined := strings.Join(buffer, "")
					c.processMessage(joined)
					// reset
					chunkCount = 0
					buffer = nil
				}
				continue
			}

			// case 3: normal single-frame message
			if strings.HasPrefix(text, "{") {
				c.processMessage(text)
				continue
			} else {
				log.Printf("Non-JSON message: %s\n", text)
			}
		}
	}
}

func (c *FirebaseClient) handleMessage(data map[string]interface{}) {
	msg := QueryMessage{
		Type: data["t"].(string),
	}

	if d, ok := data["d"].(map[string]interface{}); ok {
		// check for expired token error
		if b, ok := d["b"].(map[string]interface{}); ok {
			if s, ok := b["s"].(string); ok && s == "expired_token" {
				log.Printf("Detected expired token, attempting to reconnect...")
				if err := c.reconnect(); err != nil {
					log.Printf("Failed to reconnect after token expiration: %v", err)
				}
				return
			}
		}

		jsonData, err := json.Marshal(d)
		if err != nil {
			log.Printf("Error marshaling data: %v", err)
			return
		}
		if err := json.Unmarshal(jsonData, &msg.Data); err != nil {
			log.Printf("Error unmarshaling data: %v", err)
			return
		}
	}

	if msg.Type != "d" {
		log.Printf("Skipping non-data message type: %s", msg.Type)
		return
	}

	// subscription message
	if msg.Data.Action == "d" {
		// this is a streaming message without request ID
		c.messageMutex.Lock()
		c.streamBuffer = append(c.streamBuffer, msg.Data)
		c.messageMutex.Unlock()
		return
	}

	// handle request/response
	if msg.Data.RequestId != 0 {
		log.Printf("Handling request response for ID: %d", msg.Data.RequestId)
		c.handleRequestResponse(msg.Data.RequestId, msg.Data)
	}
}

func (c *FirebaseClient) reconnect() error {
	// close existing connection if it exists
	if c.conn != nil {
		c.conn.Close()
	}

	// attempt to establish new connection
	u := url.URL{
		Scheme:   "wss",
		Host:     c.host,
		Path:     "/.ws",
		RawQuery: fmt.Sprintf("v=5&ns=%s", c.namespace),
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to establish new connection: %v", err)
	}
	c.conn = conn

	// re-authenticate
	authRes, err := c.sendRequest("auth", QueryBody{
		Cred: &c.idToken,
	})
	if err != nil {
		return fmt.Errorf("failed to re-authenticate: %v", err)
	}

	// update user ID from auth response
	if data, ok := authRes.Data.(map[string]interface{}); ok {
		if auth, ok := data["auth"].(map[string]interface{}); ok {
			if userID, ok := auth["user_id"].(string); ok {
				c.userId = userID
			}
		}
	}

	return nil
}

func (c *FirebaseClient) handleRequestResponse(reqID int, data QueryData) {
	c.messageMutex.Lock()
	defer c.messageMutex.Unlock()

	// handle completion or standalone response
	if data.Body.Status != nil && *data.Body.Status == "ok" {
		if ch, found := c.pending[reqID]; found {
			// if we have buffered messages, this is a completion
			if len(c.streamBuffer) > 0 {
				// send all buffered messages
				for _, msg := range c.streamBuffer {
					ch <- msg.Body
				}
				// clear the buffer
				c.streamBuffer = make([]QueryData, 0)
			} else {
				// this is a standalone successful response
				ch <- data.Body
			}
			delete(c.pending, reqID)
			return
		}
	}

	// handle regular response (non-streaming, non-ok status)
	if ch, found := c.pending[reqID]; found {
		ch <- data.Body
		delete(c.pending, reqID)
	} else {
		log.Printf("No pending channel found for request %d", reqID)
	}
}

func (c *FirebaseClient) next() int {
	id := c.requestID
	c.requestID++
	return id
}

func (c *FirebaseClient) Close() {
	c.autoReconnect = false
	close(c.stopChan)
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *FirebaseClient) UserId() string {
	return c.userId
}

func (c *FirebaseClient) IdToken() string {
	return c.idToken
}
