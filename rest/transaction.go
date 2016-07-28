package rest

import (
	"github.com/sebdehne/accountingserver/storage"
	"github.com/kataras/iris"
	"github.com/sebdehne/accountingserver/domain"
	"strconv"
	"encoding/json"
	"errors"
)

type TransactionApi struct {
	store *storage.Storage
}

type TransactionSplitDto struct {
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
	Splits          []TransactionSplitDto `json:"details"`
}

type TransactionsDto struct {
	BaseAmount   int `json:"base_amount"`
	Transactions []TransactionDto `json:"transactions"`
}

func (tApi *TransactionApi) DeleteTransactionFromAccount(c *iris.Context) {
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

	if !acc.RemoveTransaction(c.Param("txId")) {
		c.SetStatusCode(iris.StatusNotFound)
		return
	}

	root.Version++
	err = tApi.store.SaveAndCommit(root)
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}

	c.SetHeader("ETag", strconv.Itoa(root.Version))
}

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

	if err := validateTransactionDto(root, in); err != nil {
		c.Error(err.Error(), iris.StatusBadRequest)
		return
	}

	// insert the tx now
	existingTx, i, found := acc.GetTransaction(inTx.Id)
	if found && existingTx.Date == inTx.Date {
		// update in place
		acc.Transactions[i] = inTx
	} else {
		if found {
			acc.RemoveTransaction(inTx.Id)
		}
		acc.AddTransaction(inTx)
	}

	root.Version++
	err = tApi.store.SaveAndCommit(root)
	if err != nil {
		c.Error(err.Error(), iris.StatusInternalServerError)
		return
	}

	c.SetHeader("ETag", strconv.Itoa(root.Version))
	if found {
		c.SetStatusCode(iris.StatusOK)
	} else {
		c.SetStatusCode(iris.StatusCreated)
	}
}

func validateTransactionDto(root domain.Root, in TransactionDto) error {
	if (in.Date < 0) {
		return errors.New("date cannot be negative")
	}

	if len(in.RemoteAccountId) > 0 {
		if len(in.RemotePartyId) > 0 {
			return errors.New("remote_party_id cannot be set when remote_account_id is set")
		}
		if _, _, found := root.GetAccount(in.RemoteAccountId); !found {
			return errors.New("remote_account_id does not exist")
		}
	} else if len(in.RemotePartyId) > 0 {
		if len(in.RemoteAccountId) > 0 {
			return errors.New("remote_account_id cannot be set when remote_party_id is set")
		}
		if _, _, found := root.GetParty(in.RemotePartyId); !found {
			return errors.New("remote_party_id does not exist")
		}
	} else {
		return errors.New("either remote_party_id or remote_account_id must be set")
	}

	for i, txS := range in.Splits {
		if _, _, found := root.GetCategory(txS.CategoryId); !found {
			return errors.New("CategoryId " + txS.CategoryId + " on " + strconv.Itoa(i) + "th transaction-detail does not exist")
		}
	}

	return nil
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
	amount := Sum(in.Splits)
	newBalance := previousBalance + amount

	return TransactionDto{
		Id:in.Id,
		Date:in.Date,
		Amount:amount,
		AccountBalance:newBalance,
		RemoteAccountId:in.RemoteAccountId,
		RemotePartyId:in.RemotePartyId,
		Splits:MapOutTransactionSpecifications(in.Splits)}, newBalance
}

func MapOutTransactionSpecifications(in []domain.TransactionSplit) []TransactionSplitDto {
	result := make([]TransactionSplitDto, 0)
	for _, txS := range in {
		result = append(result, MapOutTransactionSpecification(txS))
	}
	return result
}

func Sum(in []domain.TransactionSplit) int {
	result := 0
	for _, txS := range in {
		result += txS.Amount
	}
	return result
}

func MapOutTransactionSpecification(in domain.TransactionSplit) TransactionSplitDto {
	return TransactionSplitDto{CategoryId:in.CategoryId, Amount:in.Amount, Description:in.Description}
}

func MapInTransaction(in TransactionDto) domain.Transaction {
	return domain.Transaction{
		Id:in.Id,
		Date:in.Date,
		RemoteAccountId:in.RemoteAccountId,
		RemotePartyId:in.RemotePartyId,
		Splits:MapInTransactionSpecifications(in.Splits)}
}

func MapInTransactionSpecifications(in []TransactionSplitDto) []domain.TransactionSplit {
	result := make([]domain.TransactionSplit, 0)
	for _, txS := range in {
		result = append(result, MapInTransactionSpecification(txS))
	}
	return result
}

func MapInTransactionSpecification(in TransactionSplitDto) domain.TransactionSplit {
	return domain.TransactionSplit{CategoryId:in.CategoryId, Amount:in.Amount, Description:in.Description}
}
