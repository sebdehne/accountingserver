package kmymoney_xml

import (
	"strings"
	"strconv"
	"time"
	"github.com/sebdehne/accountingserver/domain"
)

func addTransactions(rootNode Node, accounts *Accounts) {

	for _, txNode := range rootNode.findNodeWithName("TRANSACTIONS").Nodes {
		splits := make([]domain.TransactionSplit, 0)

		splitsNode := txNode.findNodeWithName("SPLITS").Nodes

		firstSplit := splitsNode[0]
		accountId := firstSplit.AccountAttr
		payeeId := firstSplit.PayeeAttr
		localAccount, found := accounts.GetAccount(accountId)
		if !found {
			panic("Could not find local account")
		}

		// special handling of transfer between accounts
		if len(splitsNode) == 2 {
			if remoteAcc, found := accounts.GetAccount(splitsNode[1].AccountAttr); found {
				// this is a transfer
				// TX on local account
				localAccount.AddTransaction(domain.Transaction{
					Id:txNode.IdAttr,
					Date:dateStrToUnixTimestamp(txNode.PostDateAttr),
					RemoteAccountId:splitsNode[1].AccountAttr,
					Splits:[]domain.TransactionSplit{{
						CategoryId:"TRANSFER",
						Amount:valueToAmount(splitsNode[1].ValueAttr),
						Description:splitsNode[0].MemoAttr}}})

				// TX on remote account
				remoteAcc.AddTransaction(domain.Transaction{
					Id:txNode.IdAttr,
					Date:dateStrToUnixTimestamp(txNode.PostDateAttr),
					RemoteAccountId:accountId,
					Splits:[]domain.TransactionSplit{{
						CategoryId:"TRANSFER",
						Amount:valueToAmount(splitsNode[0].ValueAttr),
						Description:splitsNode[0].MemoAttr}}})

				continue
			}
		}

		// else a payment

		for i, splitNode := range splitsNode {
			if i == 0 {
				continue
			}

			splits = append(splits, domain.TransactionSplit{
				Description:splitNode.MemoAttr,
				CategoryId:splitNode.AccountAttr,
				Amount:valueToAmount(splitNode.ValueAttr)})
		}

		localAccount.AddTransaction(domain.Transaction{
			Id:txNode.IdAttr,
			Date:dateStrToUnixTimestamp(txNode.PostDateAttr),
			RemotePartyId:payeeId,
			Splits:splits})
	}
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