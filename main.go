package main

import (
	"github.com/sebdehne/accountingserver/server"
	"github.com/sebdehne/accountingserver/storage"
	"github.com/sebdehne/accountingserver/rest"
	"log"
)

func main() {

	store := storage.New("data", "accounting.json")
	err := store.InitStorage()
	if err != nil {
		log.Fatal(err)
	}

	restApi := rest.New(store)

	server.RunServer("accounting", server.Api{Version:1, Routes:[]server.Route{

		{"GET", "/categories", restApi.CategoryApi.ListCategories},
		{"PUT", "/category/:id", restApi.CategoryApi.PutCategory},
		{"DELETE", "/category/:id", restApi.CategoryApi.DeleteCategory},

		{"GET", "/accounts", restApi.AccountApi.GetAccounts},
	}})
}
