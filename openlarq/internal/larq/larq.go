package larq

import (
	"log"
	"net/http"
	"os"

	"barking.dev/openlarq/internal/auth"
	"barking.dev/openlarq/internal/firebase"
	"barking.dev/openlarq/internal/goals"
	"barking.dev/openlarq/internal/health"
	"barking.dev/openlarq/internal/intake"
	"barking.dev/openlarq/internal/user"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type App struct {
	fb *firebase.FirebaseClient
}

func NewApp() *App {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("err loading: %v", err)
	}

	app := &App{}

	return app
}

func (a *App) StartApp() error {
	r := mux.NewRouter()

	a.authenticate()
	a.registerApiRoutes(r)
	a.startServer(r)

	return nil
}

func (a *App) registerApiRoutes(r *mux.Router) {
	r.HandleFunc("/health", health.Health).Methods("GET")
	r.HandleFunc("/user-info", user.GetUserInfo(a.fb)).Methods("GET")
	r.HandleFunc("/liquid-intake", intake.GetLiquidIntake(a.fb)).Methods("GET")
	r.HandleFunc("/hydration-goal", goals.GetHydrationGoals(a.fb)).Methods("GET")
}

func (a *App) startServer(r *mux.Router) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := ":" + port

	log.Printf("API listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func (a *App) authenticate() error {
	idToken, err := auth.Authenticate()
	if err != nil {
		log.Fatal(err)
	}

	config := firebase.LoadConfig()

	// create a new Firebase client instance
	fb, err := firebase.NewFirebaseClient(config.ProjectID, config.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to create Firebase client:", err)
	}

	a.fb = fb

	// authenticate the user
	if err := a.fb.AuthenticateUser(idToken); err != nil {
		log.Fatal("Failed to authenticate user:", err)
	}

	return nil
}
