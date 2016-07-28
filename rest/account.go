package rest

import (
	"github.com/sebdehne/accountingserver/storage"
	"github.com/kataras/iris"
	"strconv"
	"encoding/json"
	"github.com/sebdehne/accountingserver/domain"
)

type AccountApi struct {
	store *storage.Storage
}

type AccountDto struct {
	Id              string `json:"id"`
	Name            string `json:"name"`
	StartingBalance int `json:"starting_balance"`
}

func (aApi *AccountApi) DeleteAccount(c *iris.Context) {
	// get the existing data
	root, err := aApi.store.Get()
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}
	id := c.Param("id")

	// validate the ETag header
	expectedVersion, err := strconv.Atoi(c.RequestHeader("ETag"))
	if err != nil {
		c.Error("Invalid ETag header", iris.StatusBadRequest)
		return
	}
	if expectedVersion != root.Version {
		c.Error("Invalid ETag header", iris.StatusConflict)
		return
	}

	if root.IsAccountInUse(id) {
		c.Error("Account is in use", iris.StatusConflict)
		return
	}

	if !root.RemoveAccount(id) {
		c.SetStatusCode(iris.StatusNotFound)
		return
	}

	root.Version++
	err = aApi.store.SaveAndCommit(root)
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}

	c.SetStatusCode(iris.StatusNoContent)
	c.SetHeader("ETag", strconv.Itoa(root.Version))
}

func (aApi *AccountApi) PutAccount(c *iris.Context) {
	// try to unmarshall request body
	in := AccountDto{}
	err := json.Unmarshal(c.Request.Body(), &in)
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}
	in.Id = c.Param("id")
	inAcc := domain.NewAccount(in.Id, in.Name, in.StartingBalance)

	// get the existing data
	root, err := aApi.store.Get()
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}

	// validate the ETag header
	expectedVersion, err := strconv.Atoi(c.RequestHeader("ETag"))
	if err != nil {
		c.Error("Invalid ETag header", iris.StatusBadRequest)
		return
	}
	if expectedVersion != root.Version {
		c.Error("Invalid ETag header", iris.StatusConflict)
		return
	}
	// validate account name
	if len(inAcc.Name) == 0 {
		c.Error("Account name cannot be empty", iris.StatusBadRequest)
		return
	}
	for _, acc := range root.Accounts {
		if acc.Id != inAcc.Id && acc.Name == inAcc.Name {
			c.Error("Account name already exists on account with id " + acc.Id, iris.StatusBadRequest)
			return
		}
	}

	// all good, update the category now
	existingAccount, _, found := root.GetAccount(inAcc.Id)
	if !found {
		root.Accounts = append(root.Accounts, inAcc)
	} else {
		// update existing account
		existingAccount.Name = inAcc.Name
		if inAcc.StartingBalance != existingAccount.StartingBalance && len(existingAccount.Transactions) > 0 {
			c.Error("Cannot change startingBalance for account because there are transactions", iris.StatusBadRequest)
			return
		}
		existingAccount.StartingBalance = inAcc.StartingBalance
	}
	root.Version++
	err = aApi.store.SaveAndCommit(root)
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}

	// prepare a response
	c.SetHeader("ETag", strconv.Itoa(root.Version))
	if found {
		c.JSON(iris.StatusOK, MapAccount(inAcc))
	} else {
		c.JSON(iris.StatusCreated, MapAccount(inAcc))
	}
}

func (aApi *AccountApi) GetAccounts(c *iris.Context) {
	root, err := aApi.store.Get()
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}

	c.SetHeader("ETag", strconv.Itoa(root.Version))
	c.JSON(200, MapAccounts(root.Accounts))
}

func MapAccounts(in []domain.Account) []AccountDto {
	result := make([]AccountDto, 0)
	for _, acc := range in {
		result = append(result, MapAccount(acc))
	}
	return result
}

func MapAccount(in domain.Account) AccountDto {
	return AccountDto{Id:in.Id, Name:in.Name, StartingBalance:in.StartingBalance}
}