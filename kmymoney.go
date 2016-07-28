package main

import "github.com/sebdehne/accountingserver/importers/kmymoney_xml"

// Imports data from KMyMoney (XML format)
func main() {
	kmymoney_xml.ImportKMyMoneyXml("/Users/sebas/Desktop/kmymoney_backup.xml")
}
