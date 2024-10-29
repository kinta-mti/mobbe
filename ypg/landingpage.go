package ypg

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// based on PG Documentation version 1.9.8 draft (internal)

type Customer struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"` //format: 8xxxxx
	Country     string `json:"country"`     //eg: IDN
	PostalCode  string `json:"postalCode"`
	Locality    string `json:"locality"` // eg: Bandung, Jakarta
	Address     string `json:"address"`
	Language    string `json:"language"` // values: en, id
}

type OrderItem struct {
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
	Amount   int64  `json:"amount"`
	Url      string `json:"url"`
	Type     string `json:"type"`
	Id       string `json:"id"` // sku of item

}

type Order struct {
	Id           string      `json:"id"` //order ID
	DisablePromo bool        `json:"disablePromo"`
	Items        []OrderItem `json:"items"`
	TokenOption  string      `json:"tokenOption"`
	TokenType    string      `json:"tokenType"`
	UseLastToken bool        `json:"useLastToken"`
}

type Urls struct {
	Selections string `json:"selections"`
	Checkout   string `json:"checkout"`
}

// data type for API gateway/IPGAPI/v1/inquiries
type InquiryReq struct {
	Amount              int64    `json:"amount"`   //total amount
	Currency            string   `json:"currency"` // ISO4217 currency code (eg. IDR)
	Customer            Customer `json:"customer"`
	Order               Order    `json:"order"`
	ReferenceUrl        string   `json:"referenceUrl"`        //Merchant redirect url after end-user completed the payment page
	PaymentSource       string   `json:"paymentSource"`       //leave it blank to use all possible payment channel
	PaymentSourceMethod string   `json:"paymentSourceMethod"` //leave it blank for normal sales flow. used for enabling authcapture process.
	Token               string   `json:"token"`               //The token given by PG from payment.received webhook, if the end-user chooses to save the credit/debit card information for future transactions.
}

type InquiryRes struct {
	Id             string   `json:"id"` // inquiry ID, NOT order ID
	CreatedTime    string   `json:"createdTime"`
	ReferenceId    string   `json:"referenceId"` //order ID yu have passed on request
	Status         string   `json:"status"`      //Status of the payment request (Unpaid, Paid)
	Amount         int64    `json:"amount"`
	Currency       string   `json:"currency"`
	PaymentSources []string `json:"paymentSources"`
	Urls           Urls     `json:"urls"`
}

type InquiryOrderItem struct {
	Name     string `json:"name"`
	Quantity int64  `json:"quantity"`
	Amount   int64  `json:"amount"`
}

type InquiryOrder struct {
	Id           string             `json:"id"` //order ID
	DisablePromo bool               `json:"disablePromo"`
	Items        []InquiryOrderItem `json:"items"`
}

type InquiryCustomer struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"` //format: 8xxxxx
	Country     string `json:"country"`     //eg: IDN
	PostalCode  string `json:"postalCode"`
}

type Inquiry struct {
	Id          string          `json:"id"` // Inquiry ID
	CreatedTime string          `json:"createdTime"`
	UpdatedTime string          `json:"updatedTime"`
	Status      string          `json:"status"` // inquiry status: failed, unpaid, partial, pending, paid
	Amount      int64           `json:"amount"`
	Currency    string          `json:"currency"`
	Customer    InquiryCustomer `json:"customer"`
	Order       InquiryOrder    `json:"order"`
}

type TransactionStatusData struct {
	Message                     string `json:"message"`
	QrCode                      string `json:"qrCode"`
	ExpireTime                  string `json:"expireTime"`
	VaNumber                    string `json:"vaNumber"`
	AuthenticationModule        string `json:"authenticationModule"`
	ChallengeAuthenticationCode string `json:"challengeAuthenticationCode"`
	ProcessingCode              string `json:"processingCode"`
	AuthenticationCode          string `json:"authenticationCode"`
	CardType                    string `json:"cardType"`
	CardNetwork                 string `json:"cardNetwork"`
}

type Transaction struct {
	Id                 string                `json:"id"` // Transaction ID
	CreatedTime        string                `json:"createdTime"`
	UpdatedTime        string                `json:"updatedTime"`
	Status             string                `json:"status"`     // transaction status: submitted, declined, pending, validated, failed, processing, authorized, captured
	StatusCode         string                `json:"statusCode"` // (payment received) Raw status code from PG or acquiring transaction response.
	StatusData         TransactionStatusData `json:"statusData"`
	Amount             int64                 `json:"amount"`
	Currency           string                `json:"currency"`
	PaymentSource      string                `json:"paymentSource"`
	NetworkReferenceId string                `json:"networkReferenceId"`
	AuthorizationCode  string                `json:"authorizationCode"`
	PaymentSourceData  interface{}           `json:"paymentSourceData"`
	VoidStatus         string                `json:"voidStatus"`
}

type AccessTokenRequest struct {
	GrantType string `json:"grant_type"` //value always client_credentials
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"` //Access token expiration period in seconds
}

var accessToken *AccessToken
var lastAccessTokenTime = time.Now()
var merchant_apiKey = ""
var merchant_secretkey = ""
var apim_key = ""
var apim_secret = ""
var path_uri = ""
var path_accessToken = ""
var path_inquiries = ""

func Init(merchantApiKey, merchantSecretKey, apimKey, apimSecret, pathUri, pathAccessToken, pathInquiries string) {
	if merchantApiKey == "" || merchantSecretKey == "" || apimKey == "" || apimSecret == "" || pathUri == "" || pathAccessToken == "" || pathInquiries == "" {
		log.Println("[ypg.init] configuration missing, please check ypg configuration")
	} else {
		merchant_apiKey = merchantApiKey
		merchant_secretkey = merchantSecretKey
		apim_key = apimKey
		apim_secret = apimSecret
		path_uri = pathUri
		path_accessToken = pathAccessToken
		path_inquiries = pathInquiries
	}
}

func RefreshAccessToken() *AccessToken {
	if accessToken == nil || (accessToken != nil && lastAccessTokenTime.Add(time.Second*time.Duration(accessToken.ExpiresIn)).Before(time.Now())) {

		data := url.Values{}
		data.Set("grant_type", "client_credential")

		u, _ := url.ParseRequestURI(path_uri)
		u.Path = path_accessToken
		urlStr := u.String()

		client := &http.Client{}
		r, _ := http.NewRequest(http.MethodPost, urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
		r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		r.SetBasicAuth(apim_key, apim_secret)

		resp, err := client.Do(r)
		if err != nil {
			log.Print("[ypg.RefreshAccessToken]", err)
			return nil
		}
		defer resp.Body.Close()
		bodyText, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Print("[ypg.RefreshAccessToken]", err)
			return nil
		}

		json.Unmarshal(bodyText, &accessToken)
		lastAccessTokenTime = time.Now()
		return accessToken
	} else {
		return accessToken
	}

}

func Inquiries(payload []byte) string {
	log.Print("[ypg.Inquiries] function called")
	u, _ := url.ParseRequestURI(path_uri)
	u.Path = path_inquiries
	urlStr := u.String()
	r, err := http.NewRequest(http.MethodPost, urlStr, bytes.NewReader(payload))
	if err != nil {
		log.Println("[ypg.Inquiries] Error creating request.\n[ERROR] -", err)
		return "error"
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("X-Api-Key", merchant_apiKey)
	r.Header.Add("Authorization", "Bearer "+accessToken.AccessToken)
	//log.Println("accessToken.\n[data] -", "Bearer "+accessToken.AccessToken)
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		log.Println("[ypg.Inquiries]\n[ERROR] -", err)
		return "error"
	}

	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	//log.Println("inquiry raw body.\n[data] -", string(bodyText))
	if err != nil {
		log.Println("[ypg.Inquiries] \n[ERROR] -", err)
		return "error"
	}
	var inquiryResponse *InquiryRes
	json.Unmarshal(bodyText, &inquiryResponse)
	return inquiryResponse.Urls.Checkout
}

// HMAC SHA256(Requested raw body + `.` + timestamp, merchant.SECRETKEY).digest('hex')
func SignatureValidation(rawBody string) string {
	return ""
}

// MD5(merchantSECRETKEY + signature + timestamp).digest('hex')
func SignatureResponse(signature string, timestamp string) string {
	hash := md5.Sum([]byte(merchant_secretkey + signature + timestamp))
	return hex.EncodeToString(hash[:])
}

// HMAC SHA256(Requested raw body + `.` + timestamp, merchant.SECRETKEY).digest('hex')
func IsValidSignature(requestRawBody []byte, signature, timestamp string) bool {
	mac := hmac.New(sha256.New, []byte(merchant_secretkey))

	mac.Write(append(requestRawBody, "."+timestamp...))
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(expectedMAC, []byte(signature))
}
