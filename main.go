package main

import (
	"github.com/sebdehne/accountingserver/server"
	"github.com/sebdehne/accountingserver/storage"
	"github.com/sebdehne/accountingserver/rest"
	"log"
	"github.com/kataras/iris"
)

func main() {

	store := storage.New("data", "accounting.json")
	err := store.InitStorage()
	if err != nil {
		log.Fatal(err)
	}

	restApi := rest.New(store)

	server.RunServer(":8081", server.Api{Prefix:"accounting", Version:1, Routes:[]server.Route{
		{"GET", "/categories", restApi.CategoryApi.ListCategories},
		{"PUT", "/category/:id", restApi.CategoryApi.PutCategory},
		{"DELETE", "/category/:id", restApi.CategoryApi.DeleteCategory},

		{"GET", "/accounts", restApi.AccountApi.GetAccounts},
		{"PUT", "/account/:id", restApi.AccountApi.PutAccount},
		{"DELETE", "/account/:id", restApi.AccountApi.DeleteAccount},

		{"GET", "/parties", restApi.PartApi.ListParties},
		{"PUT", "/party/:id", restApi.PartApi.PutParty},
		{"DELETE", "/party/:id", restApi.PartApi.DeleteParty},

		{"GET", "/account/:id/transactions", restApi.TransactionApi.ListTransactionsForAccount},
		{"PUT", "/account/:id/transaction/:txId", restApi.TransactionApi.PutTransactionForAccount},
		{"DELETE", "/account/:id/transaction/:txId", restApi.TransactionApi.DeleteTransactionFromAccount},
	}}, server.Api{Prefix:"webapp", Version:1, Routes:[]server.Route{
		{"GET", "/*filepath", iris.StaticHandler("./webapp", 2, true, false, []string{"index.html"})},
	}})

}
