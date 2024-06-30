package tests

import (
	"gotransact/apps/accounts/handlers"
	"gotransact/apps/accounts/models"
	"gotransact/apps/accounts/utils"
	accountValidator "gotransact/apps/accounts/validators"
	log "gotransact/log"
	"gotransact/pkg/db"
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestLogin_success(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	accountValidator.Init()
	log.Init()
	//create a user
	password, _ := utils.HashPassword("Password@123")
	existinguser := models.User{
		FirstName: "testfirstname",
		LastName:  "testlastname",
		Email:     "test@gmail.com",
		Password:  password,
	}

	db.DB.Create(&existinguser)

	input := handlers.LoginInput{
		Email:    "test@gmail.com",
		Password: "Password@123",
	}

	status, message, data := handlers.Login(input)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "Logged in successfull", message)
	assert.NotEmpty(t, data["token"])
	ClearDatabase()
	CloseTestDb()
}

func TestLogin_InvalidEmail(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	accountValidator.Init()
	log.Init()
	password, _ := utils.HashPassword("Trellis@123")
	existinguser := models.User{
		FirstName: "Saish",
		LastName:  "Naik",
		Email:     "nsaish@trellissoft.ai",
		Password:  password,
	}

	db.DB.Create(&existinguser)

	// Attempt to log in with invalid email
	input := handlers.LoginInput{
		Email:    "nsaish@trellissoft.ai",
		Password: "Trellis@123",
	}

	status, message, data := handlers.Login(input)

	assert.Equal(t, http.StatusUnauthorized, status)
	assert.Equal(t, "invalid username or password", message)
	assert.Empty(t, data)
	CloseTestDb()
}

func TestLogin_InvalidPassword(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	accountValidator.Init()
	log.Init()
	// Create a user
	password, _ := utils.HashPassword("Password@123")
	existinguser := models.User{
		FirstName: "testfirstname",
		LastName:  "testlastname",
		Email:     "test@gmail.com",
		Password:  password,
	}

	db.DB.Create(&existinguser)

	input := handlers.LoginInput{
		Email:    "test@gmail.com",
		Password: "WrongPassword@123",
	}

	status, message, data := handlers.Login(input)

	assert.Equal(t, http.StatusUnauthorized, status)
	assert.Equal(t, "invalid username or password", message)
	assert.Empty(t, data)
	ClearDatabase()
	CloseTestDb()
}