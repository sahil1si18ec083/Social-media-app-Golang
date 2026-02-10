package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/db"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/env"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/store"

	_ "github.com/lib/pq"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file: ", err)
	}
	x := env.GetString("ADDR", ":8080")
	fmt.Println(x)
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5432/socialnetwork?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}
	fmt.Print(cfg.db.addr)
	db, err := db.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	defer db.Close()

	store := store.NewStorage(db)
	fmt.Print(store)

	app := &application{
		config: cfg,
		store:  store,
	}

	mux := app.mount()
	err = app.run(mux)
	fmt.Print("bye")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("bye.............")

}
