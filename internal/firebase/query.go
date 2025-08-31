package firebase

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type QueryMessage struct {
	Type string    `json:"t"`
	Data QueryData `json:"d"`
}

type QueryData struct {
	Action    string    `json:"a"`
	Body      QueryBody `json:"b"`
	RequestId int       `json:"r"`
}

type QueryBody struct {
	Path   *string      `json:"p,omitempty"`
	Query  *QueryParams `json:"q,omitempty"`
	Type   *int         `json:"t,omitempty"`
	Cred   *string      `json:"cred,omitempty"`
	Status *string      `json:"s,omitempty"`
	Data   interface{}  `json:"d,omitempty"`
}

type QueryParams struct {
	Index          *string `json:"i,omitempty"`
	StartTime      *string `json:"sp,omitempty"`
	EndTime        *string `json:"ep,omitempty"`
	StartInclusive *string `json:"sin,omitempty"`
	StartName      *string `json:"sn,omitempty"`
	EndName        *string `json:"en,omitempty"`
	EndInclusive   *string `json:"ein,omitempty"`
	Limit          *int    `json:"l,omitempty"`
	ViewFrom       *string `json:"vf,omitempty"`
}

func (c *FirebaseClient) Query(path string, queryParams QueryParams) (QueryBody, error) {
	msg := QueryMessage{
		Type: "d",
		Data: QueryData{
			Action: "q",
			Body: QueryBody{
				Path:  &path,
				Query: &queryParams,
				Type:  &[]int{1}[0],
			},
			RequestId: c.next(),
		},
	}

	jsonData, _ := json.MarshalIndent(msg, "", "  ")
	log.Printf("Sending query message:\n%s", string(jsonData))

	// create a channel to receive the response
	responseChan := make(chan QueryBody, 1)

	// register the response channel
	c.messageMutex.Lock()
	c.pending[msg.Data.RequestId] = responseChan
	c.messageMutex.Unlock()

	// send the query
	c.writeMutex.Lock()
	err := c.conn.WriteJSON(msg)
	c.writeMutex.Unlock()

	if err != nil {
		c.messageMutex.Lock()
		delete(c.pending, msg.Data.RequestId)
		c.messageMutex.Unlock()
		return QueryBody{}, err
	}

	// wait for the response
	select {
	case response := <-responseChan:
		return response, nil
	case <-time.After(30 * time.Second):
		return QueryBody{}, fmt.Errorf("timeout waiting for response")
	}
}

func NewQueryParams(params map[string]interface{}) QueryParams {
	queryParams := QueryParams{}

	if viewFrom, ok := params["viewFrom"].(*string); ok {
		queryParams.ViewFrom = viewFrom
	}
	if limit, ok := params["limit"].(*int); ok {
		queryParams.Limit = limit
	}
	if startTime, ok := params["startTime"].(*string); ok {
		queryParams.StartTime = startTime
	}
	if endTime, ok := params["endTime"].(*string); ok {
		queryParams.EndTime = endTime
	}
	if startName, ok := params["startName"].(*string); ok {
		queryParams.StartName = startName
	}
	if index, ok := params["index"].(*string); ok {
		queryParams.Index = index
	}
	if endName, ok := params["endName"].(*string); ok {
		queryParams.EndName = endName
	}

	return queryParams
}
