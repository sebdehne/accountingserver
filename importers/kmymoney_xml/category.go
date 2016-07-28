package kmymoney_xml

import (
	"github.com/sebdehne/accountingserver/domain"
)

func extractCategories(txs []Transaction, root Node) map[string]domain.Category {
	result := make(map[string]domain.Category)

	for _, tx := range txs {
		for _, split := range tx.Splits {
			if split.CategoryAccountId == "TRANSFER" {
				break
			}
			if _, found := result[split.CategoryAccountId]; !found {
				result[split.CategoryAccountId] = extractCategory(split.CategoryAccountId, root)
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
