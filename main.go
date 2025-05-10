package main

import (
	"flag"
	"fmt"
	"github.com/jdmukiibs/femProject/internal/app"
	"github.com/jdmukiibs/femProject/internal/routes"
	"net/http"
	"time"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 8080, "port to run the server on")
	flag.Parse()
	myApp, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	myApp.Logger.Printf("we are running on port: %d\n", port)

	r := routes.SetupRoutes(myApp)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	err = server.ListenAndServe()
	if err != nil {
		myApp.Logger.Fatal(err)
	}
}
