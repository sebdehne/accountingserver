package domain

func New() Root {
	return Root{Version:1, Parts:make([]Part, 0), Categories:make([]Category, 0), Accounts:make([]Account, 0)}
}

type Root struct {
	Version int
	Parts []Part
	Categories []Category
	Accounts []Account
}

type Part struct {
	Id string
	Name string
}

type Category struct {
	Id string
	Name string
}

type Account struct {
	Id string
	StartingBalance int
	Transactions []Transaction
}

type Transaction struct {
	Id string
	Date int64
	Amount int
	NewAccountBalance int
	RemoteAccountId string
	RemotePartId string
	Details []TransactionSpecification
}

type TransactionSpecification struct {
	Id string
	Parent string
	CategoryId string
	Amount int
	Description string
}
