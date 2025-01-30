package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/andresilvase/gobid/internal/api"
	"github.com/andresilvase/gobid/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Failed to load .env file")
		panic(err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s",
		os.Getenv("GOBID_DATABASE_USER"),
		os.Getenv("GOBID_DATABASE_PASSWORD"),
		os.Getenv("GOBID_DATABASE_HOST"),
		os.Getenv("GOBID_DATABASE_PORT"),
		os.Getenv("GOBID_DATABASE_NAME"),
	))

	if err != nil {
		fmt.Println("Failed to connect to database")
		panic(err)
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		fmt.Println("Failed to ping database")
		panic(err)
	}

	api := api.Api{
		Router:      chi.NewMux(),
		UserService: services.NewUserService(pool),
	}

	api.BindRoutes()

	fmt.Println("Server running on port :3080")

	if err := http.ListenAndServe("localhost:3080", api.Router); err != nil {
		fmt.Println("Failed to start server")
		panic(err)
	}
}
