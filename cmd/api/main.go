package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/sahil1si18ec083/Social-media-app-Golang/internal/env"
)

func main() {

	_ = godotenv.Load()
	x := env.GetString("ADDR", ":8080")
	fmt.Println(x)
	cfg := config{
		addr: env.GetString("ADDR", ":8080"),
	}

	app := &application{
		config: cfg,
	}

	mux := app.mount()
	err := app.run(mux)
	if err != nil {
		log.Fatal(err)
	}

}
