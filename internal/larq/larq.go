package larq

import (
	"log"
	"net/http"

	"barking.dev/larq-api/internal/auth"
	"barking.dev/larq-api/internal/firebase"
	"barking.dev/larq-api/internal/goals"
	"barking.dev/larq-api/internal/health"
	"barking.dev/larq-api/internal/intake"
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

	r.HandleFunc("/liquid-intake", intake.GetLiquidIntake(a.fb)).Methods("GET")
	r.HandleFunc("/hydration-goal", goals.GetHydrationGoals(a.fb)).Methods("GET")
}

func (a *App) startServer(r *mux.Router) {
	addr := ":8080"

	log.Printf("API listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}

func (a *App) authenticate() error {
	idToken, err := auth.Authenticate()
	if err != nil {
		log.Fatal(err)
	}

	config := firebase.LoadConfig()

	// Create a new Firebase client instance
	fb, err := firebase.NewFirebaseClient(config.ProjectID, config.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to create Firebase client:", err)
	}

	a.fb = fb

	// Authenticate the user
	if err := a.fb.AuthenticateUser(idToken); err != nil {
		log.Fatal("Failed to authenticate user:", err)
	}

	return nil
}
