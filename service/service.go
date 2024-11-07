package service

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kinta-mti/mobbe/db"
	"github.com/kinta-mti/mobbe/ypg"
)

type HelloWorld struct {
	Request  string `json:"request"`
	Response string `json:"response"`
}

type Item struct {
	Sku   string `json:"sku"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
	Qty   int64  `json:"qty"`
	Url   string `json:"url"`
	Type  string `json:"type"`
}

type Checkout struct {
	Id                  string `json:"id"` //order ID
	Items               []Item `json:"items"`
	CustomerName        string `json:"customerName"`
	CustomerEmail       string `json:"customerEmail"`
	CustomerPhoneNumber string `json:"customerPhoneNumber"`
	CustomerCountry     string `json:"customerCountry"`
	CustomerPostalCode  string `json:"customerPostalCode"`
	CustomerLocality    string `json:"customerLocality"`
	CustomerAddress     string `json:"customerAddress"`
	FCMToken            string `json:"fcmToken"`
}

type CheckoutRes struct {
	Url string `json:"url"`
}

type WHReq struct {
	Type           string          `json:"type"` // value: payment.validate / payment.received
	Transaction    ypg.Transaction `json:"transaction"`
	Inquiry        ypg.Inquiry     `json:"inquiry"`
	Token          string          `json:"token"`
	TokenExpiredAt string          `json:"token_expired_at"`
}

type WHPaymentValidateRes struct {
	Status            string      `json:"status"`
	ValidateSignature string      `json:"validateSignature"`
	Inquiry           ypg.Inquiry `json:"inquiry"`
}

type WHPaymentReceivedRes struct {
	Status            string `json:"status"`
	ValidateSignature string `json:"validateSignature"`
}

type WHError struct {
	ErrorCode    string `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}

func Init(port string) {
	log.Println("[service.init] called!!")
	if port == "" {
		log.Println("[service.init] configuration missing, please check server configuration")
	} else {
		router := gin.Default()
		router.POST("/checkout", PostCheckout)
		router.POST("/webhook", PostWebhook)
		router.GET("/hello", GetWorld)
		router.Run(":" + port)
		log.Println("[service.init] server run on port: " + port)
	}

}

// post a checkout
func PostCheckout(c *gin.Context) {
	var checkout Checkout
	// Call BindJSON to bind the received JSON to
	// checkout.
	if err := c.BindJSON(&checkout); err != nil {
		log.Print("[service.postcheckout] error BindJSON:" + err.Error())
		return
	}
	//insert new order to
	if db.InsertNewUserOrder(checkout.Id, checkout.FCMToken) > 0 {

	} else {
		c.JSON(http.StatusInternalServerError, WHError{ErrorCode: "3000", ErrorMessage: "Order ID already created"})
	}

	var inquiryRequest ypg.InquiryReq

	inquiryRequest.Currency = "IDR"
	inquiryRequest.Customer = ypg.Customer{
		Name:        checkout.CustomerName,
		Email:       checkout.CustomerEmail,
		PhoneNumber: checkout.CustomerPhoneNumber,
		Address:     checkout.CustomerAddress,
		PostalCode:  checkout.CustomerPostalCode,
		Locality:    checkout.CustomerLocality,
		Country:     checkout.CustomerCountry,
		Language:    "id",
	}
	inquiryRequest.Order = ypg.Order{}
	inquiryRequest.Order.Id = checkout.Id
	inquiryRequest.Order.DisablePromo = true
	inquiryRequest.Order.TokenOption = "true"
	inquiryRequest.Order.TokenType = "static"
	inquiryRequest.Order.UseLastToken = false
	inquiryRequest.Order.Items = []ypg.OrderItem{}
	var totalAmount int64 = 0
	for i := 0; i < len(checkout.Items); i++ {
		totalAmount += (checkout.Items[i].Qty * checkout.Items[i].Price)
		inquiryRequest.Order.Items = append(
			inquiryRequest.Order.Items,
			ypg.OrderItem{
				Name:     checkout.Items[i].Name,
				Quantity: checkout.Items[i].Qty,
				Amount:   checkout.Items[i].Price,
				Url:      checkout.Items[i].Url,
				Type:     checkout.Items[i].Type,
			})
	}
	inquiryRequest.Amount = totalAmount
	//inquiryRequest.PaymentSource
	//inquiryRequest.PaymentSourceMethod
	inquiryRequest.ReferenceUrl = "https://ypgmerchant.test.negriku.id/afterPayment"
	inquiryRequest.Token = ""
	//ypg.RefreshAccessToken()
	ypg.RefreshAccessToken()
	var payload, err = json.Marshal(inquiryRequest)
	if err != nil {
		log.Println("[service.postcheckout] Error on read body.\n[ERROR] -", err)
	} else {
		var checkouturl = ypg.Inquiries(payload)
		//c.IndentedJSON(http.StatusCreated, CheckoutRes{Url: checkouturl})
		if checkouturl == "error" {
			c.JSON(http.StatusInternalServerError, "")
		} else {
			c.JSON(http.StatusOK, CheckoutRes{Url: checkouturl})
		}

	}

}

func PostWebhook(c *gin.Context) {
	log.Print("[service.postWebhook] function called")

	//GET raw request body
	requestRawBody, err := c.GetRawData()
	if err != nil {
		// Handle error
		log.Print("[service.postWebhook] error get request body:" + err.Error())
		return
	} else {
		log.Print("[service.postWebhook] requestrawbody length " + string(requestRawBody))
	}

	// convert request raw data to a struct
	var webhookRequest WHReq
	if err := json.Unmarshal(requestRawBody, &webhookRequest); err != nil {
		log.Print("[service.postWebhook] error BindJSON:" + err.Error())
		return
	} else {
		log.Print("[service.postWebhook] after bindjson:" + webhookRequest.Type)
	}

	//extract Signature
	var signature = ""
	for name, values := range c.Request.Header {
		// Loop over all values for the name.
		for _, value := range values {
			log.Print("[service.postWebhook]", name, "-", value)

			if name == "Signature" {
				signature = value
			}
		}
	}

	c.IndentedJSON(webHookResponse(requestRawBody, signature, webhookRequest))

}

func webHookResponse(requestRawBody []byte, signature string, webhookRequest WHReq) (int, any) {
	if signature == "" {
		return http.StatusBadRequest, WHError{
			ErrorCode:    e_headerMissingSignature_c,
			ErrorMessage: e_headerMissingSignature_m}
	} else {
		//payment validate response, implement always ok
		var signsplit = strings.Split(signature, ";")

		//validate Received Signature
		if ypg.IsValidSignature(requestRawBody, signsplit[0], signsplit[1]) {
			if webhookRequest.Type == "payment.validate" {
				//payment.validate response. implement always OK
				if db.UpdatePaymentValidate(webhookRequest.Inquiry.Order.Id, string(requestRawBody)) {
					var whResponse WHPaymentValidateRes = WHPaymentValidateRes{
						Status:            "ok",
						ValidateSignature: ypg.SignatureResponse(signsplit[0], signsplit[1]),
						Inquiry:           webhookRequest.Inquiry}
					return http.StatusOK, whResponse
				} else {
					var whResponse WHError = WHError{
						ErrorCode:    e_orderIDNotFound_c,
						ErrorMessage: e_orderIDNotFound_m,
					}
					return http.StatusBadRequest, whResponse
				}

			} else if webhookRequest.Type == "payment.received" {
				//payment.received response. implement always OK
				if db.UpdatePaymentReceived(webhookRequest.Inquiry.Order.Id, string(requestRawBody)) {
					return http.StatusOK, WHPaymentReceivedRes{
						Status:            "ok",
						ValidateSignature: ypg.SignatureResponse(signsplit[0], signsplit[1]),
					}
				} else {
					var whResponse WHError = WHError{
						ErrorCode:    e_orderIDNotFound_c,
						ErrorMessage: e_orderIDNotFound_m,
					}
					return http.StatusBadRequest, whResponse
				}

			} else {
				return http.StatusBadRequest, WHError{
					ErrorCode:    e_unhandledException_c,
					ErrorMessage: e_unhandledException_m + ": webhookRequest.Type==" + webhookRequest.Type}
			}

		} else {
			//invalid signature
			log.Print("[service.webHookResponse] INVALID SIGNATURE")
			return http.StatusBadRequest, WHError{
				ErrorCode:    e_invalidSignature_c,
				ErrorMessage: e_invalidSignature_m}
		}
	}
}

func GetWorld(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, HelloWorld{Request: "hello", Response: "world"})
}

// constant list
const (
	e_headerMissingSignature_c = "0001"
	e_headerMissingSignature_m = "Header Missing Signature"
	e_invalidSignature_c       = "1001"
	e_invalidSignature_m       = "Invalid Signature"
	e_orderIDNotFound_c        = "2001"
	e_orderIDNotFound_m        = "Order ID not found"
	e_unhandledException_c     = "9999"
	e_unhandledException_m     = "unhandled exception"
)
