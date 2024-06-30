package tests

import (
	"gotransact/apps/accounts/handlers"
	"gotransact/apps/accounts/models"
	accountValidator "gotransact/apps/accounts/validators"
	log "gotransact/log"
	"gotransact/pkg/db"
	"net/http"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSignup_success(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	accountValidator.Init()
	log.Init()
	input := handlers.SignupUser{
		FirstName:   "testfirstname",
		LastName:    "testlastname",
		Email:       "test@gmail.com",
		CompanyName: "trellissoft",
		Password:    "Password@123",
	}

	status, message, data := handlers.Signup(input)

	assert.Equal(t, http.StatusOK, status)
	assert.Equal(t, "User created successfully", message)
	assert.Equal(t, map[string]interface{}{}, data)

	var user models.User
	err := db.DB.Where("email = ?", input.Email).First(&user).Error
	assert.NoError(t, err)
	assert.Equal(t, input.FirstName, user.FirstName)
	assert.Equal(t, input.LastName, user.LastName)
	assert.Equal(t, input.Email, user.Email)

	var company models.Company
	err = db.DB.Where("name = ?", input.CompanyName).First(&company).Error
	assert.NoError(t, err)
	assert.Equal(t, input.CompanyName, company.Name)
	ClearDatabase()
	CloseTestDb()
}

func TestSignup_EmailAreadyExist(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	accountValidator.Init()

	existingUser := models.User{
		FirstName: "testfirstname",
		LastName:  "testlastname",
		Email:     "test@gmail.com",
		Company: models.Company{
			Name: "trellissoft",
		},
	}

	db.DB.Create(&existingUser)

	input := handlers.SignupUser{
		FirstName:   "otherfirstname",
		LastName:    "otherlastname",
		Email:       "test@gmail.com",
		CompanyName: "Google",
		Password:    "Password@123",
	}

	status, message, data := handlers.Signup(input)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "email already exists", message)
	assert.Equal(t, map[string]interface{}{}, data)
	ClearDatabase()
	CloseTestDb()
}

func TestSignup_CompanyAlreadyExist(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	accountValidator.Init()
	log.Init()
	existingUser := models.User{
		FirstName: "testfirstname",
		LastName:  "testlastname",
		Email:     "test@gmail.com",
		Company: models.Company{
			Name: "trellissoft",
		},
	}

	db.DB.Create(&existingUser)

	input := handlers.SignupUser{
		FirstName:   "othername",
		LastName:    "otherlastname",
		Email:       "testother@gmail.com",
		CompanyName: "trellissoft",
		Password:    "Password@123",
	}

	status, message, data := handlers.Signup(input)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "company already exists", message)
	assert.Equal(t, map[string]interface{}{}, data)
	ClearDatabase()
	CloseTestDb()
}

func TestSignup_InvaldPassword(t *testing.T) {
	SetupTestDb()
	ClearDatabase()
	log.Init()
	input := handlers.SignupUser{
		FirstName:   "otherfirstname",
		LastName:    "otherlastname",
		Email:       "test@gmail.com",
		CompanyName: "trellissoft",
		Password:    "password@123",
	}

	status, message, data := handlers.Signup(input)

	assert.Equal(t, http.StatusBadRequest, status)
	assert.Equal(t, "Password should contain atleast one upper case character,one lower case character,one number and one special character", message)
	assert.Equal(t, map[string]interface{}{}, data)
	ClearDatabase()
	CloseTestDb()
}
