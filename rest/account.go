package rest

import (
	"github.com/sebdehne/accountingserver/storage"
	"github.com/kataras/iris"
	"strconv"
)

type AccountApi struct {
	store storage.Storage
}

type AccountDto struct {
	Id              string
	StartingBalance int
}

func (aApi *AccountApi) GetAccounts(c *iris.Context) {
	root, err := aApi.store.Get()
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}

	c.SetHeader("ETag", strconv.Itoa(root.Version))
	result := make([]AccountDto, 0)

	for _, acc := range root.Accounts {
		result = append(result, AccountDto{acc.Id, acc.StartingBalance})
	}

	c.JSON(200, result)
}
