package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type UserOrder struct {
	OrderID                string `json:"orderID"`
	Usertoken              string `json:"usertoken"`
	ReceivedTime           string `json:"receivedTime"`
	PaymentValidateTime    string `json:"paymentValidateTime"`
	PaymentReceivedTime    string `json:"paymentReceivedTime"`
	PaymentValidatePayload string `json:"paymentValidatePayload"`
	PaymentReceivedPayload string `json:"paymentReceivedPayload"`
}

var db_name = ""
var db_user = ""
var db_pass = ""

func Init(dbName, dbUser, dbPass string) {
	log.Println("[db.Init] called!!")
	if dbName == "" || dbUser == "" || dbPass == "" {
		log.Println("[db.Init] configuration missing, please check database configuration")
	} else {
		db_name = dbName
		db_user = dbUser
		db_pass = dbPass
		testConn()
	}
}

func testConn() {
	db, err := sql.Open("mysql", db_user+":"+db_pass+"@/"+db_name)
	defer db.Close()
	if err != nil {
		log.Panic("[db.testConn] ", err.Error())
	} else {
		log.Print("[db.testConn] Success!!")
	}

}

func InsertNewUserOrder(orderId, usertoken string) bool {
	db, err := sql.Open("mysql", db_user+":"+db_pass+"@/"+db_name)
	if err != nil {
		log.Print("[db.InsertNewUserOrder]", err.Error())
	}
	defer db.Close()
	result, err := db.Exec("INSERT INTO UserOrder (OrderID, Usertoken, ReceivedTime) VALUES ($1, $2, $3)", orderId, usertoken, time.Now())
	if err != nil {
		log.Print(err.Error())
	}
	row, _ := result.RowsAffected()
	if row > 0 {
		return true
	} else {
		return false
	}

}

func GetOrder(orderId string) {
	db, err := sql.Open("mysql", db_user+":"+db_pass+"@/"+db_name)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var order UserOrder
	err = db.QueryRow("select orderID,usertoken,receivedTime,PaymentValidateTime,PaymentReceivedTime,PaymentValidatePayload,PaymentReceivedPayload where orderID = ?", orderId).
		Scan(&order.OrderID, &order.Usertoken, &order.ReceivedTime, &order.PaymentValidateTime, &order.PaymentReceivedTime, &order.PaymentValidatePayload, &order.PaymentReceivedPayload)
	if err != nil {
		panic(err.Error())
	}

}

/*
func UpdatePaymentValidate(orderId, paymentValidatePayload string) bool {
	db, err := sql.Open("mysql", "root:tree@/negrikui_ypgmerchant")
	if err != nil {
		log.Print(err.Error())
	}
	defer db.Close()
	result, err := db.Exec("INSERT INTO UserOrder (OrderID, Usertoken, ReceivedTime) VALUES ($1, $2, $3)", orderId, usertoken, time.Now())
	if err != nil {
		log.Print(err.Error())
	}
	row, _ := result.RowsAffected()
	if row > 0 {
		return true
	} else {
		return false
	}
}
*/
