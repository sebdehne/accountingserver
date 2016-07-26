package rest

import (
	"github.com/sebdehne/accountingserver/storage"
	"github.com/kataras/iris"
	"github.com/sebdehne/accountingserver/domain"
	"strconv"
	"encoding/json"
)

type TransactionApi struct {
	store *storage.Storage
}

type TransactionSpecificationDto struct {
	Id          string `json:"id"`
	CategoryId  string `json:"category_id"`
	Amount      int `json:"amount"`
	Description string `json:"description"`
}

type TransactionDto struct {
	Id              string `json:"id"`
	Date            int64 `json:"date"`
	Amount          int `json:"amount"`
	AccountBalance  int `json:"account_balance"`
	RemoteAccountId string `json:"remote_account_id"`
	RemotePartyId   string `json:"remote_party_id"`
	Details         []TransactionSpecificationDto `json:"details"`
}

type TransactionsDto struct {
	BaseAmount   int `json:"base_amount"`
	Transactions []TransactionDto `json:"transactions"`
}

// inserts/moves the TX behind the last TX with the same date
func (tApi *TransactionApi) PutTransactionForAccount(c *iris.Context) {
	// try to unmarshall request body
	in := TransactionDto{}
	err := json.Unmarshal(c.Request.Body(), &in)
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}
	in.Id = c.Param("txId")
	inTx := MapInTransaction(in)

	root, err := tApi.store.Get()
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}

	acc, _, found := root.GetAccount(c.Param("id"))
	if !found {
		c.Error("Account does not exist", iris.StatusNotFound)
		return
	}
	acc.RemoveTransaction(inTx.Id)

	// TODO validate TX
	// date positive
	// account / party reference must exist
	// in detail:
	// -- category must exist
	// -- amount 0 or more

}

func (tApi *TransactionApi) ListTransactionsForAccount(c *iris.Context) {
	pageFilter := ExtractPageFilter(c)
	dateFilter, err := ExtractDateFilter(c)
	if err != nil {
		c.Error(err.Error(), iris.StatusBadRequest)
		return
	}
	root, err := tApi.store.Get()
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}

	acc, _, found := root.GetAccount(c.Param("id"))
	if !found {
		c.Error("Account does not exist", iris.StatusNotFound)
		return
	}

	r := acc.GetTransactions(dateFilter, pageFilter)

	c.SetHeader("ETag", strconv.Itoa(root.Version))
	c.JSON(iris.StatusOK, TransactionsDto{BaseAmount:r.BaseAmount, Transactions:MapOutTransactions(r.Transactions, r.BaseAmount)})
}

func MapOutTransactions(in []domain.Transaction, previousAmount int) []TransactionDto {
	result := make([]TransactionDto, 0)
	var tmp TransactionDto

	for _, tx := range in {
		tmp, previousAmount = MapOutTransaction(tx, previousAmount)
		result = append(result, tmp)
	}
	return result
}

func MapOutTransaction(in domain.Transaction, previousBalance int) (TransactionDto, int) {
	amount := Sum(in.Details)
	newBalance := previousBalance + amount

	return TransactionDto{
		Id:in.Id,
		Date:in.Date,
		Amount:amount,
		AccountBalance:newBalance,
		RemoteAccountId:in.RemoteAccountId,
		RemotePartyId:in.RemotePartyId,
		Details:MapOutTransactionSpecifications(in.Details)}, newBalance
}

func MapOutTransactionSpecifications(in []domain.TransactionSpecification) []TransactionSpecificationDto {
	result := make([]TransactionSpecificationDto, 0)
	for _, txS := range in {
		result = append(result, MapOutTransactionSpecification(txS))
	}
	return result
}

func Sum(in []domain.TransactionSpecification) int {
	result := 0
	for _, txS := range in {
		result += txS.Amount
	}
	return result
}

func MapOutTransactionSpecification(in domain.TransactionSpecification) TransactionSpecificationDto {
	return TransactionSpecificationDto{Id:in.Id, CategoryId:in.CategoryId, Amount:in.Amount, Description:in.Description}
}

func MapInTransaction(in TransactionDto) domain.Transaction {
	return domain.Transaction{
		Id:in.Id,
		Date:in.Date,
		RemoteAccountId:in.RemoteAccountId,
		RemotePartyId:in.RemotePartyId,
		Details:MapInTransactionSpecifications(in.Details)}
}

func MapInTransactionSpecifications(in []TransactionSpecificationDto) []domain.TransactionSpecification {
	result := make([]domain.TransactionSpecification, 0)
	for _, txS := range in {
		result = append(result, MapInTransactionSpecification(txS))
	}
	return result
}

func MapInTransactionSpecification(in TransactionSpecificationDto) domain.TransactionSpecification {
	return domain.TransactionSpecification{Id:in.Id, CategoryId:in.CategoryId, Amount:in.Amount, Description:in.Description}
}
