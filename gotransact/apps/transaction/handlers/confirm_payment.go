package handlers

import (
	base_utils "gotransact/apps/base/utils"
	log "gotransact/log"
	transaction "gotransact/apps/transaction/models"
	"fmt"
	"html/template"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gotransact/pkg/db"
	"strings"
	"github.com/google/uuid"
)

func ConfirmPaymentHandler(c *gin.Context) {

	log.InfoLogger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"url":    c.Request.URL.String(),
	}).Info("confirm payment Request received")

	transactionIdStr := c.Query("transaction_id")
	statusStr := c.Query("status")

	_, message, data := ConfirmPayment(transactionIdStr, statusStr)

	// Create a map for template data
	tmplData := map[string]interface{}{
		"TransactionID": transactionIdStr,
		"Amount":        data["Amount"],
		"Message":       message,
	}
	// Select the template based on the message
	var tmpl *template.Template
	var err error

	if message == "Transaction successful" {
		tmpl, err = template.ParseFiles("gotransact/apps/transaction/templates/payment_success.html")
	} else if message == "Transaction Canceled" {
		tmpl, err = template.ParseFiles("gotransact/apps/transaction/templates/payment_fail.html")
	} else {
		c.JSON(http.StatusInternalServerError, base_utils.Response{
			Status:  http.StatusInternalServerError,
			Message: "Unknown transaction status",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, base_utils.Response{
			Status:  http.StatusInternalServerError,
			Message: "Template parsing error",
		})
		return
	}

	// Render the template
	c.Writer.Header().Set("Content-Type", "text/html")
	tmpl.Execute(c.Writer, tmplData)
}



func ConfirmPayment(transactionIdStr, statusStr string) (int, string, map[string]interface{}) {
	log.InfoLogger.WithFields(logrus.Fields{}).Info("Attempted to confirm/cancel payment transaction-id=", transactionIdStr)
	transactionId, err := uuid.Parse(transactionIdStr)
	fmt.Println("parsed", transactionId)
	if err != nil {
		return http.StatusBadRequest, "Invalid transaction ID", map[string]interface{}{}
	}

	var transactionRequest transaction.TransactionRequest
	if err := db.DB.Where("internal_id = ?", transactionId).First(&transactionRequest).Error; err != nil {
		return http.StatusBadRequest, "transaction request not found", map[string]interface{}{}
	}

	var trasactionHistory transaction.TransactionHistory
	trasactionHistory.TransactionID = transactionRequest.ID
	trasactionHistory.Description = transactionRequest.Description
	trasactionHistory.Amount = transactionRequest.Amount

	if strings.EqualFold(statusStr, "true") {

		if err := db.DB.Model(&transactionRequest).Where("id = ?", transactionRequest.ID).Update("status", transaction.StatusSuccess).Error; err != nil {
			log.ErrorLogger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Error completing payment transaction-id=", transactionRequest.InternalID)
			return http.StatusInternalServerError, "Failed to confirm the payment", map[string]interface{}{}
		}
		log.InfoLogger.WithFields(logrus.Fields{}).Info("Payment completed transaction-id=", transactionRequest.InternalID)
		{
			trasactionHistory.Status = transaction.StatusSuccess
		}
	} else {

		if err := db.DB.Model(&transactionRequest).Where("id = ?", transactionRequest.ID).Update("status", transaction.StatusFailed).Error; err != nil {
			log.ErrorLogger.WithFields(logrus.Fields{
				"error": err.Error(),
			}).Error("Error canceling the payment transaction-id=", transactionRequest.InternalID)
			return http.StatusInternalServerError, "Failed to confirm the payment", map[string]interface{}{}
		}
		log.InfoLogger.WithFields(logrus.Fields{}).Info("Payment canceled transaction-id=", transactionRequest.InternalID)
		{
			trasactionHistory.Status = transaction.StatusFailed
		}

	}

	if err := db.DB.Create(&trasactionHistory).Error; err != nil {
		log.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("failed to update transaction history to confirm/cancel transaction-id=", transactionRequest.InternalID)
		return http.StatusInternalServerError, "Failed to update transaction history", map[string]interface{}{}
	}
	log.InfoLogger.WithFields(logrus.Fields{}).Info("transaction history updated to confirm/cancel transaction-id=", transactionRequest.InternalID)
	if strings.EqualFold(statusStr, "true") {
		return http.StatusOK, "Transaction successful", map[string]interface{}{"Amount":transactionRequest.Amount}
	} else {
		return http.StatusOK, "Transaction Canceled", map[string]interface{}{"Amount":transactionRequest.Amount}
	}
}