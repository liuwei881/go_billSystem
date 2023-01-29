package main

import (
	"BillSystem/controllers"
	_ "BillSystem/database"
)

func main() {
	con := controllers.InitController()
	con.Run()
}
