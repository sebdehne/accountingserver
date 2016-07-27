package kmymoney_xml

import (
	"encoding/xml"
	"github.com/sebdehne/accountingserver/domain"
)

type Node struct {
	IdAttr       string `xml:"id,attr"`
	NameAttr     string `xml:"name,attr"`
	ValueAttr    string `xml:"value,attr"`
	MemoAttr     string `xml:"memo,attr"`
	AccountAttr  string `xml:"account,attr"`
	PostDateAttr string `xml:"postdate,attr"`
	PayeeAttr    string `xml:"payee,attr"`

	XMLName      xml.Name
	Content      []byte `xml:",innerxml"`
	Nodes        []Node `xml:",any"`
}

func (n *Node) findNodeWithName(name string) Node {
	for _, childNode := range n.Nodes {
		if childNode.XMLName.Local == name {
			return childNode
		}
	}
	panic("Could not find node " + name)
}

type Accounts []domain.Account

func (accs Accounts) GetAccount(accountId string) (domain.Account, bool) {
	for _, accInfo := range accs {
		if accInfo.Id == accountId {
			return accInfo, true
		}
	}
	return domain.Account{}, false
}

type AccountInfo struct {
	AccountId       string
	InstitutionName string
	AccountName     string
}

type Transaction struct {
	Id              string
	Date            int64
	AccountId       string
	RemoteAccountId string
	RemotePartyId   string
	Splits          []Split
}

type Split struct {
	Memo              string
	CategoryAccountId string
	Amount            int
}

