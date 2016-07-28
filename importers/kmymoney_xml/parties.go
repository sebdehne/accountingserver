package kmymoney_xml

import (
	"github.com/sebdehne/accountingserver/domain"
)

func extractParties(txs []Transaction, root Node) map[string]domain.Party {
	result := make(map[string]domain.Party)

	for _, tx := range txs {
		if tx.RemotePartyId == "" {
			// skip transfers
			continue
		}

		if _, found := result[tx.RemotePartyId]; !found {
			result[tx.RemotePartyId] = extractParty(tx.RemotePartyId, root)
		}
	}

	return result
}

func extractParty(payeeId string, root Node) domain.Party {
	for _, payeeNode := range root.findNodeWithName("PAYEES").Nodes {
		if (payeeNode.IdAttr == payeeId) {
			return domain.Party{Id:payeeNode.IdAttr, Name:payeeNode.NameAttr}
		}
	}
	panic("Could not find payee with id " + payeeId)
}