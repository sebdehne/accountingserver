package main

import "github.com/sebdehne/accountingserver/importers/kmymoney_xml"

func main() {
	kmymoney_xml.ImportKMyMoneyXml("/Users/sebas/Desktop/kmymoney_backup.xml")
}
