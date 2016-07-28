package kmymoney_xml

import (
	"io/ioutil"
	"log"
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/sebdehne/accountingserver/storage"
	"github.com/sebdehne/accountingserver/domain"
)

func ImportKMyMoneyXml(filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer(data)
	dec := xml.NewDecoder(buf)

	var n Node
	err = dec.Decode(&n)
	if err != nil {
		log.Fatal(err)
	}

	accounts := extractAccounts(n)
	txs := extractTransactions(n, accounts)
	cats := extractCategories(txs, n)
	parties := extractParties(txs, n)

	store := storage.New("data", "accounting.json")
	err = store.InitStorage()
	if err != nil {
		log.Fatal(err)
	}

	root, err := store.Get()
	if err != nil {
		panic(err)
	}

	// import all accounts
	for _, acc := range accounts {
		if _, _, found := root.GetAccount(acc.Id); !found {
			root.Accounts = append(root.Accounts, acc)
		} else {
			fmt.Println("Account " + acc.Id + " already exists")
		}
	}

	// import all parties
	for partyId, party := range parties {
		if _, _, found := root.GetParty(partyId); !found {
			root.Parties = append(root.Parties, party)
		} else {
			fmt.Println("Party " + partyId + " already exists")
		}
	}

	// import all categories
	for catId, cat := range cats {
		if _, _, found := root.GetCategory(catId); !found {
			root.Categories = append(root.Categories, cat)
		} else {
			fmt.Println("Category " + catId + " already exists")
		}
	}

	// import all transactions
	for _, tx := range txs {
		acc, _, found := root.GetAccount(tx.AccountId)
		if !found {
			panic("Could not import tx " + tx.Id + " because the account " + tx.AccountId + " could not be found")
		}

		if _, _, found := acc.GetTransaction(tx.Id); found {
			fmt.Println("TX " + tx.Id + " already exists")
			continue
		}

		// map to domain.*
		newSplits := make([]domain.TransactionSplit, 0)
		for _, split := range tx.Splits {
			newSplits = append(newSplits, domain.TransactionSplit{
				CategoryId:split.CategoryAccountId,
				Amount:split.Amount,
				Description:split.Memo})
		}
		newTx := domain.Transaction{
			Id:tx.Id,
			Date:tx.Date,
			RemoteAccountId:tx.RemoteAccountId,
			RemotePartyId:tx.RemotePartyId,
			Splits:newSplits}

		acc.AddTransaction(newTx)
	}

	// Save
	root.Version++
	err = store.SaveAndCommit(root)
	if err != nil {
		panic(err)
	}

}

