package main

import (
	"fmt"
	"log"
	"net/http"
	"./gql"
	"./postgres"
	"./server"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/graphql-go/graphql"
)

func main()  {
	router, db := initializeAPI()
	defer db.Close()

	// Listen on port 4000 and if there's an error log it and exit
	log.Fatal(http.ListenAndServe(":4000", router))
}

func initializeAPI(*chi.Mux, *postgres.Db)  {
	router := chi.NewRouter()

	// Create a new connection to our pg database
	db, err := postgres.New(
		postgres.ConnString("localhost", 5432, "postgres", "graphql"),
	) if err != nil {
		log.Fatal(err)
	}

	rootQuery := gql.NewRoot(db)

	// Create a new graphql schema, passing in the root query
	sc, err := graphql.NewSchema(
		graphql.SchemaConfig{Query: rootQuery.Query}
	) if err != nil {
		fmt.Println("Error creating schema: ", err)
	}

	// Create a server struct that holds a pointer to our database as well
	// as the address of our graphql schema

	s: = server.Server{
		GqlSchema: &sc
	}

	// Add some middleware to our router
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.StripSlashes,
		middleware.Recover,
	)

	// Create the graphql route with a Server method to handle it
	router.Post("/graphql", s.GraphQL())

	return router, db
}