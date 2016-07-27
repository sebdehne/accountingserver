package domain

func New() Root {
	return Root{Version:1, Parties:make([]Party, 0), Categories:make([]Category, 0), Accounts:make([]Account, 0)}
}

func NewAccount(id, name string, startingBalance int) Account {
	return Account{Id:id, Name:name, StartingBalance:startingBalance, Transactions:make([]Transaction, 0)}
}

type Root struct {
	Version    int
	Parties    []Party
	Categories []Category
	Accounts   []Account
}

type Party struct {
	Id   string
	Name string
}

type Category struct {
	Id   string
	Name string
}

type Account struct {
	Id              string
	Name            string
	StartingBalance int
	Transactions    []Transaction
}

type Transaction struct {
	Id              string
	Date            int64
	RemoteAccountId string
	RemotePartyId   string
	Splits          []TransactionSplit
}

type TransactionSplit struct {
	CategoryId  string
	Amount      int
	Description string
}

type DateFilter struct {
	FromDate *int64
	ToDate   *int64
}

type PageFilter struct {
	Offset int
	Limit  int
}
