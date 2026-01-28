package main

import (
	"fmt"
	"log"
)

func main() {

	fmt.Print("hello")
	cfg := config{
		addr: ":8080",
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
