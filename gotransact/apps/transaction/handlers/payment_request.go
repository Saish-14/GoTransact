package handlers

import (
	"bytes"
	"fmt"
	account "gotransact/apps/accounts/models"
	base_utils "gotransact/apps/base/utils"
	transaction "gotransact/apps/transaction/models"
	transaction_validators "gotransact/apps/transaction/validators"
	"gotransact/pkg/db"
	"html/template"
	"strconv"
	"time"

	// "gotransact/apps/transaction/functions"
	// "gotransact/apps/transaction/utils"
	logger "gotransact/log"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type PostPaymentInput struct {
	CardNumber  string `json:"cardnumber" binding:"required" validate:"card_number" `
	ExpiryDate  string `json:"expirydate" binding:"required" validate:"expiry_date" `
	Cvv         string `json:"cvv" validate:"cvv" binding:"required"`
	Amount      string `json:"amount" binding:"required" validate:"amount"`
	Description string `json:"description" `
}


func PaymentRequest(c *gin.Context) {

	logger.InfoLogger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"url":    c.Request.URL.String(),
	}).Info("Post Payment Request received")

	var Postpaymentinput PostPaymentInput
	if err := c.ShouldBindJSON(&Postpaymentinput); err != nil {
		c.JSON(http.StatusBadRequest, base_utils.Response{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data:    map[string]interface{}{"data": err.Error()},
		})
		return
	}

	UserFromRequest, exist := c.Get("user")
	if !exist {
		c.JSON(http.StatusBadRequest, base_utils.Response{
			Status:  http.StatusBadRequest,
			Message: "User not found in token",
			Data:    map[string]interface{}{},
		})
		return
	}

	user, ok := UserFromRequest.(account.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assert user type"})
		return
	}

	status, message, data := PostPayment(Postpaymentinput, user)

	c.JSON(status, base_utils.Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}


func PostPayment(Postpaymentinput PostPaymentInput, user account.User) (int, string, map[string]interface{}) {

	logger.InfoLogger.WithFields(logrus.Fields{}).Info("Attempted to create transaction request with email ", user.Email, " id ", user.InternalID)

	if err := transaction_validators.GetValidator().Struct(Postpaymentinput); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		errors := make(map[string]string)
		for _, fieldErr := range validationErrors {
			fieldName := fieldErr.Field()
			tag := fieldErr.Tag()
			errors[fieldName] = transaction_validators.CustomErrorMessages[tag]
		}
		return http.StatusBadRequest, "error while validating", map[string]interface{}{}
	}

	floatAmount, _ := strconv.ParseFloat(Postpaymentinput.Amount, 64)

	var gateway transaction.PaymentGateway
	if err := db.DB.Where("slug = ?", "card").First(&gateway).Error; err != nil {
		return http.StatusBadRequest, "invalid payment type", map[string]interface{}{}
	}

	TransactionRequest := transaction.TransactionRequest{
		UserID:             user.ID,
		Status:             transaction.StatusProcessing,
		Description:        Postpaymentinput.Description,
		Amount:             floatAmount,
		PaymentGatewayMethodID: gateway.ID,
	}

	if err := db.DB.Create(&TransactionRequest).Error; err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error creating record in transaction request transaction-id=", TransactionRequest.InternalID)
		return http.StatusInternalServerError, "internal server error", map[string]interface{}{}
	}
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("created record in transaction request with email ", user.Email, " id ", user.InternalID)

	TransactionHistory := transaction.TransactionHistory{
		TransactionID: TransactionRequest.ID,
		Status:        TransactionRequest.Status,
		Description:   TransactionRequest.Description,
		Amount:        TransactionRequest.Amount,
	}

	if err := db.DB.Create(&TransactionHistory).Error; err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error creating record in transaction history")
		return http.StatusInternalServerError, "internal server error", map[string]interface{}{}
	}

	logger.InfoLogger.WithFields(logrus.Fields{}).Info("Created record in transaction history with email ", user.Email, " id ", user.InternalID)

	go SendTransactionMail(user, TransactionRequest)

	return http.StatusOK, "success", map[string]interface{}{"transaction ID": TransactionRequest.InternalID}
}


type TemplateData struct {
	Username     string
	TrasactionID uuid.UUID
	Amount       float64
	ConfirmURL   string
	CancelURL    string
	DateTime     time.Time
}

func SendTransactionMail(user account.User, request transaction.TransactionRequest) {

	logger.InfoLogger.WithFields(logrus.Fields{
		"email": user.Email,
		"id":    user.InternalID,
	}).Info("Attempted to send confirm payment mail")

	// Parse the HTML template
	tmpl, err := template.ParseFiles("gotransact/apps/transaction/templates/transaction_email.html")
	if err != nil {
		fmt.Printf("Error parsing email template: %s", err)
	}

	// Create a buffer to hold the executed template
	var body bytes.Buffer

	baseURL := "http://localhost:8080/api/confirm-payment"
	params := url.Values{}
	params.Add("transaction_id", request.InternalID.String())
	params.Add("status", "true")
	ConfirmActionURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	baseURL = "http://localhost:8080/api/confirm-payment"
	params = url.Values{}
	params.Add("transaction_id", request.InternalID.String())
	params.Add("status", "false")
	CancelActionURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Execute the template with the data
	TemplateData := TemplateData{
		Username:     user.FirstName,
		TrasactionID: request.InternalID,
		Amount:       request.Amount,
		ConfirmURL:   ConfirmActionURL,
		CancelURL:    CancelActionURL,
	}
	fmt.Println(TemplateData)
	if err := tmpl.Execute(&body, TemplateData); err != nil {
		fmt.Printf("Error executing email template: %s", err)
	}

	fmt.Println()
	// Set E-Mail body as HTML
	base_utils.SendMail(user.Email, "Payment Confirmation Required", body.String())

}