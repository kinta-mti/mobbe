package endpoint

import (
	"encoding/json"
	"io"
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
	if port == "" {
		log.Println("[endpoint.init] configuration missing, please check server configuration")
	} else {
		router := gin.Default()
		router.POST("/checkout", PostCheckout)
		router.POST("/webhook", PostWebhook)
		router.GET("/hello", GetWorld)
		router.Run(":" + port)
		log.Println("[endpoint.init] server run on port: " + port)
	}

}

// post a checkout
func PostCheckout(c *gin.Context) {
	var checkout Checkout
	// Call BindJSON to bind the received JSON to
	// checkout.
	if err := c.BindJSON(&checkout); err != nil {
		log.Print("[endpoint.postcheckout] error BindJSON:" + err.Error())
		return
	}
	//insert new order to
	if db.InsertNewUserOrder(checkout.Id, checkout.FCMToken) {

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
		log.Println("[endpoint.postcheckout] Error on read body.\n[ERROR] -", err)
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
	var webhookRequest WHReq

	// Call BindJSON to bind the received JSON to
	// newAlbum.
	log.Print("[endpoint.postWebhook] function called")

	// check request structure
	if err := c.BindJSON(&webhookRequest); err != nil {
		log.Print("error BindJSON:" + err.Error())
		return
	}

	//GET raw request body
	requestRawBody, err := io.ReadAll(c.Request.Body)
	if err != nil {
		// Handle error
		log.Print("[endpoint.postWebhook] error get request body:" + err.Error())
		return
	}

	//extract Signature
	var signature = ""
	for name, values := range c.Request.Header {
		// Loop over all values for the name.
		for _, value := range values {
			log.Print("[endpoint.postWebhook]", name, "-", value)

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
			ErrorCode:    "0001",
			ErrorMessage: "Header Missing Signature"}
	} else {
		//payment validate response, implement always ok
		var signsplit = strings.Split(signature, ";")

		//validate Received Signature
		if ypg.IsValidSignature(requestRawBody, signsplit[0], signsplit[1]) {
			if webhookRequest.Type == "payment.validate" {
				//payment.validate response. implement always OK
				var whResponse WHPaymentValidateRes = WHPaymentValidateRes{
					Status:            "ok",
					ValidateSignature: ypg.SignatureResponse(signsplit[0], signsplit[1]),
					Inquiry:           webhookRequest.Inquiry}
				return http.StatusOK, whResponse
			} else if webhookRequest.Type == "payment.received" {
				//payment.received response. implement always OK
				return http.StatusOK, WHPaymentReceivedRes{
					Status:            "ok",
					ValidateSignature: ypg.SignatureResponse(signsplit[0], signsplit[1]),
				}
			} else {
				return http.StatusBadRequest, WHError{
					ErrorCode:    "9001",
					ErrorMessage: "unhandled webhook type"}
			}

		} else {
			//invalid signature
			log.Print("INVALID SIGNATURE")
			return http.StatusBadRequest, WHError{
				ErrorCode:    "1001",
				ErrorMessage: "Invalid Signature"}
		}
	}
}

func GetWorld(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, HelloWorld{Request: "hello", Response: "world"})
}
