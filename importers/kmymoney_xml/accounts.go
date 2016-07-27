package kmymoney_xml

import (
	"github.com/sebdehne/accountingserver/domain"
)

func extractAccounts(rootNode Node) Accounts {
	result := Accounts{}

	accountInfos := extractAccountIds(rootNode.findNodeWithName("INSTITUTIONS"))
	accountInfos = enrichNames(accountInfos, rootNode)

	for _, accountInfo := range accountInfos {
		result = append(result, domain.NewAccount(accountInfo.AccountId, accountInfo.InstitutionName + " - " + accountInfo.AccountName, 0))
	}

	return result
}

func enrichNames(accountInfos []AccountInfo, rootNode Node) []AccountInfo {
	result := make([]AccountInfo, 0)

	for _, accountNode := range rootNode.findNodeWithName("ACCOUNTS").Nodes {
		for _, accountInfo := range accountInfos {
			if accountNode.IdAttr == accountInfo.AccountId {
				accountInfo.AccountName = accountNode.NameAttr
				result = append(result, accountInfo)
			}
		}
	}

	return result
}

func extractAccountIds(institutionsNode Node) []AccountInfo {
	result := make([]AccountInfo, 0)

	for _, institutionNode := range institutionsNode.Nodes {
		accountIds := institutionNode.findNodeWithName("ACCOUNTIDS")

		for _, accountId := range accountIds.Nodes {
			result = append(result, AccountInfo{AccountId:accountId.IdAttr, InstitutionName:institutionNode.NameAttr})
		}
	}

	return result
}

