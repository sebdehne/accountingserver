package kmymoney_xml

import (
	"strings"
	"strconv"
	"time"
)

func extractTransactions(rootNode Node, accounts Accounts) []Transaction {
	result := make([]Transaction, 0)

	txNodes := rootNode.findNodeWithName("TRANSACTIONS").Nodes

	for _, txNode := range txNodes {
		splits := make([]Split, 0)

		splitsNode := txNode.findNodeWithName("SPLITS").Nodes

		// special handling of transfer between accounts
		if len(splitsNode) == 2 {
			if _, found := accounts.GetAccount(splitsNode[1].AccountAttr); found {
				// this is a transfer
				tx := Transaction{
					Id:txNode.IdAttr,
					Date:dateStrToUnixTimestamp(txNode.PostDateAttr),
					AccountId:splitsNode[0].AccountAttr,
					RemoteAccountId:splitsNode[1].AccountAttr,
					Splits:[]Split{{
						Memo:splitsNode[0].MemoAttr,
						Amount:valueToAmount(splitsNode[1].ValueAttr),
						CategoryAccountId:"TRANSFER"}}}
				result = append(result, tx)
				continue
			}
		}

		// else a payment

		accountId := ""
		payeeId := ""
		for i, splitNode := range splitsNode {
			if i == 0 {
				accountId = splitNode.AccountAttr
				payeeId = splitNode.PayeeAttr
				continue
			}

			splits = append(splits, Split{
				Memo:splitNode.MemoAttr,
				CategoryAccountId:splitNode.AccountAttr,
				Amount:valueToAmount(splitNode.ValueAttr)})
		}

		tx := Transaction{
			Id:txNode.IdAttr,
			Date:dateStrToUnixTimestamp(txNode.PostDateAttr),
			AccountId:accountId,
			RemotePartyId:payeeId,
			Splits:splits,
		}

		result = append(result, tx)
	}

	return result
}

func valueToAmount(value string) int {
	parts := strings.Split(value, "/")
	ore, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	ore *= 100
	devider, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	ore /= devider
	return ore * -1
}

const dateLayout = "2006-01-02"

func dateStrToUnixTimestamp(dateStr string) int64 {
	t, err := time.Parse(dateLayout, dateStr)
	if err != nil {
		panic(err)
	}
	return t.Unix()
}