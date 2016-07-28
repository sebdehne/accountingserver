package kmymoney_xml

import (
	"github.com/sebdehne/accountingserver/domain"
)

func extractCategories(accounts Accounts, root Node) map[string]domain.Category {
	result := make(map[string]domain.Category)

	for _,acc := range accounts {
		for _, tx := range acc.Transactions {
			for _, split := range tx.Splits {
				if split.CategoryId == "TRANSFER" {
					break
				}
				if _, found := result[split.CategoryId]; !found {
					result[split.CategoryId] = extractCategory(split.CategoryId, root)
				}
			}
		}
	}

	return result
}

func extractCategory(categoryId string, root Node) domain.Category {
	for _, catNode := range root.findNodeWithName("ACCOUNTS").Nodes {
		if (catNode.IdAttr == categoryId) {
			return domain.Category{Id:catNode.IdAttr, Name:catNode.NameAttr}
		}
	}
	panic("Could not find category with id " + categoryId)
}
