package kmymoney_xml

import (
	"io/ioutil"
	"log"
	"bytes"
	"encoding/xml"
	"fmt"
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
	fmt.Println(accounts)

	extractTransactions(n, accounts)

}

