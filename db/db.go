package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/kinta-mti/mobbe/config"
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

func InsertNewUserOrder(orderId, usertoken string, cfgDB config.DBConnInfo) bool {
	db, err := sql.Open("mysql", cfgDB.User+":"+cfgDB.Pass+"@/"+cfgDB.Name)
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

func GetOrder(orderId string, cfgDB config.DBConnInfo) {
	db, err := sql.Open("mysql", cfgDB.User+":"+cfgDB.Pass+"@/"+cfgDB.Name)
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
