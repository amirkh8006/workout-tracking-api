package main

import (
	"femProject/internal/app"
	"femProject/internal/routes"
	"flag"
	"fmt"
	"net/http"
	"time"
)

func main() {
	var port int
	flag.IntVar(&port , "port" , 8080, "Go Backend Server Port")
	flag.Parse()


	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	defer app.DB.Close()
	

	app.Logger.Println("We Are Running Our App")



	r := routes.SetUpRoutes(app)

	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: r,
		IdleTimeout: time.Minute,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	app.Logger.Printf("We Are Running On Port %d", port)
	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatal(err)
	}
}
