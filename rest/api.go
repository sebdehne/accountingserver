package rest

import "github.com/sebdehne/accountingserver/storage"

func New(s storage.Storage) RestApi {
	return RestApi{store:s, CategoryApi:CategoryApi{store:s}}
}

type RestApi struct {
	store       storage.Storage
	CategoryApi CategoryApi
}

