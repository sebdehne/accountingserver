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

func (c *Root) Get(id string) (Category, int, bool) {
	for i, cat := range c.Categories {
		if id == cat.Id {
			return cat, i, true
		}
	}
	return Category{}, 0, false
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
	Date int64
	Amount int
	NewBalance int
	RemoteAccountId string
	RemotePartId string
	Details []TransactionSpecification
}

type TransactionSpecification struct {
	CategoryId string
	Amount int
	Description string
}
