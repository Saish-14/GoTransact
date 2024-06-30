package test

import (
	accountmodels "gotransact/apps/accounts/models"
	models "gotransact/apps/base/models"
	transactionmodels "gotransact/apps/transaction/models"
	transaction_handlers "gotransact/apps/transaction/handlers"
	"gotransact/pkg/db"
	log "gotransact/log"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestConfirmPayment_Success(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	log.Init()
	user := accountmodels.User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password",
	}
	db.DB.Create(&user)

	gateway := transactionmodels.PaymentGateway{
		Slug:  "card",
		Label: "Card",
	}
	db.DB.Create(&gateway)

	transactionRequest := transactionmodels.TransactionRequest{
		UserID:             user.ID,
		Status:             transactionmodels.StatusProcessing,
		Description:        "Test Transaction",
		Amount:             100.0,
		PaymentGatewayMethodID: gateway.ID,
		Base: models.Base{
			InternalID: uuid.New(),
		},
	}
	db.DB.Create(&transactionRequest)

	status, message, data := transaction_handlers.ConfirmPayment(transactionRequest.InternalID.String(), "true")

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "Transaction successfull", message)
	assert.Equal(t, map[string]interface{}{}, data)

	var updatedTransaction transactionmodels.TransactionRequest
	db.DB.Where("internal_id = ?", transactionRequest.InternalID).First(&updatedTransaction)
	assert.Equal(t, transactionmodels.StatusSuccess, updatedTransaction.Status)

	var transactionHistory transactionmodels.TransactionHistory
	db.DB.Where("transaction_id = ?", updatedTransaction.ID).First(&transactionHistory)
	assert.Equal(t, transactionmodels.StatusSuccess, transactionHistory.Status)

	CloseTestDb()
}

func TestConfirmPayment_Failed(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	log.Init()
	user := accountmodels.User{
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
		Password:  "password",
	}
	db.DB.Create(&user)

	gateway := transactionmodels.PaymentGateway{
		Slug:  "card",
		Label: "Card",
	}
	db.DB.Create(&gateway)

	transactionRequest := transactionmodels.TransactionRequest{
		UserID:             user.ID,
		Status:             transactionmodels.StatusProcessing,
		Description:        "Test Transaction",
		Amount:             100.0,
		PaymentGatewayMethodID: gateway.ID,
		Base: models.Base{
			InternalID: uuid.New(),
		},
	}
	db.DB.Create(&transactionRequest)

	status, message, data := transaction_handlers.ConfirmPayment(transactionRequest.InternalID.String(), "false")

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "Transaction Canceled", message)
	assert.Equal(t, map[string]interface{}{}, data)

	var updatedTransaction transactionmodels.TransactionRequest
	db.DB.Where("internal_id = ?", transactionRequest.InternalID).First(&updatedTransaction)
	assert.Equal(t, transactionmodels.StatusFailed, updatedTransaction.Status)

	var transactionHistory transactionmodels.TransactionHistory
	db.DB.Where("transaction_id = ?", updatedTransaction.ID).First(&transactionHistory)
	assert.Equal(t, transactionmodels.StatusFailed, transactionHistory.Status)

	CloseTestDb()
}

func TestConfirmPayment_InvalidTransactionID(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	log.Init()
	invalidTransactionID := "invalid-uuid"
	gateway := transactionmodels.PaymentGateway{
		Slug:  "card",
		Label: "Card",
	}
	db.DB.Create(&gateway)

	status, message, data := transaction_handlers.ConfirmPayment(invalidTransactionID, "true")

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "Invalid transaction ID", message)
	assert.Equal(t, map[string]interface{}{}, data)

	CloseTestDb()
}

func TestConfirmPayment_TransactionNotFound(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	log.Init()
	validUUID := uuid.New().String()

	gateway := transactionmodels.PaymentGateway{
		Slug:  "card",
		Label: "Card",
	}
	db.DB.Create(&gateway)

	status, message, data := transaction_handlers.ConfirmPayment(validUUID, "true")

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "transaction request not found", message)
	assert.Equal(t, map[string]interface{}{}, data)

	CloseTestDb()
}
