package kmymoney_xml

import (
	"io/ioutil"
	"log"
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/sebdehne/accountingserver/storage"
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
	addTransactions(n, &accounts)
	cats := extractCategories(accounts, n)
	parties := extractParties(accounts, n)

	store := storage.New("data", "accounting.json")
	err = store.InitStorage()
	if err != nil {
		log.Fatal(err)
	}

	root, err := store.Get()
	if err != nil {
		panic(err)
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

	// import all accounts
	for _, acc := range accounts {
		if _, _, found := root.GetAccount(acc.Id); !found {
			root.Accounts = append(root.Accounts, acc)
		} else {
			fmt.Println("Account " + acc.Id + " already exists")
		}
	}

	// Save
	root.Version++
	err = store.SaveAndCommit(root)
	if err != nil {
		panic(err)
	}

}

