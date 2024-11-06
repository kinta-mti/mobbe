package db

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type UserOrder struct {
	OrderID                string         `json:"orderID"`
	Usertoken              string         `json:"usertoken"`
	ReceivedTime           sql.NullTime   `json:"receivedTime"`
	PaymentValidateTime    sql.NullTime   `json:"paymentValidateTime"`
	PaymentReceivedTime    sql.NullTime   `json:"paymentReceivedTime"`
	PaymentValidatePayload sql.NullString `json:"paymentValidatePayload"`
	PaymentReceivedPayload sql.NullString `json:"paymentReceivedPayload"`
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

	if err != nil {
		log.Panic("[db.testConn] ", err.Error())
	} else {
		defer db.Close()
		log.Print("[db.testConn] Success!!")
		InsertNewUserOrder("test", "user token of test")
		log.Print("[db.testConn] " + GetOrder("test").Usertoken)
	}

}

func InsertNewUserOrder(orderId, usertoken string) int64 {
	db, err := sql.Open("mysql", db_user+":"+db_pass+"@/"+db_name+"?parseTime=true")
	if err != nil {
		log.Print("[db.InsertNewUserOrder]", err.Error())
	}
	defer db.Close()
	result, err := db.Exec("INSERT INTO UserOrder (OrderID, Usertoken, ReceivedTime) VALUES (?, ?, ?)", orderId, usertoken, time.Now())
	if err != nil {
		log.Print("[db.InsertNewUserOrder]" + err.Error())
		return 0
	} else {
		rowsAffected, _ := result.RowsAffected()
		return rowsAffected

	}

}

func GetOrder(orderId string) UserOrder {
	db, err := sql.Open("mysql", db_user+":"+db_pass+"@/"+db_name+"?parseTime=true")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var order UserOrder
	err = db.QueryRow("select OrderID,Usertoken,ReceivedTime,PaymentValidateTime,PaymentReceivedTime,PaymentValidatePayload,PaymentReceivedPayload from UserOrder where OrderID = ?", orderId).
		Scan(&order.OrderID, &order.Usertoken, &order.ReceivedTime, &order.PaymentValidateTime, &order.PaymentReceivedTime, &order.PaymentValidatePayload, &order.PaymentReceivedPayload)
	if err != nil {
		log.Panic("[db.GetOrder]" + err.Error())
	}
	return order
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
