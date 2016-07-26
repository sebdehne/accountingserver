package rest

import "github.com/sebdehne/accountingserver/storage"

func New(s storage.Storage) RestApi {
	return RestApi{s, CategoryApi{s}, AccountApi{s}}
}

type RestApi struct {
	store       storage.Storage
	CategoryApi CategoryApi
	AccountApi  AccountApi
}

