package models

//import (
//	orm "BillSystem/database"
//	"gorm.io/gorm"
//)

type Fest struct {
	ID		int64	`json:"id"`
	Test 	string	`json:"test"`
}

type AlyBill struct {
	ID				int64 	`json:"id"`
	ProductName		string 	`json:"ProductName"`
	PaymentAmount	string	`json:"PaymentAmount"`
	BillAccountName	string	`json:"BillAccountName"`
	ProductDetail	string 	`json:"ProductDetail"`
	BillingCycle	string	`json:"BillingCycle"`
	ApportionDepart	string	`json:"ApportionDepart"`
}


type CenterDepart struct {
	ID		int64	`json:"id"`
	Center	string	`json:"Center"`
}

