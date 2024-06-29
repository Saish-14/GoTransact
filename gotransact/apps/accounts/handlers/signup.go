package handlers

import (
	"gotransact/apps/accounts/models"
	// "gotransact/apps/accounts/utils"
	"gotransact/apps/accounts/validators"
	base_utils "gotransact/apps/base/utils"
	logger "gotransact/log"
	"gotransact/pkg/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Struct for the signup request
type SignupUser struct {
	FirstName   string `json:"first_name" binding:"required"`
	LastName    string `json:"last_name" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
	CompanyName string `json:"company_name" binding:"required"`
}

// SignupHandler godoc
// @Summary Registers a new user
// @Description Create a new user account
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param   user     body    SignupUser     true        "User data"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /register [post]
func SignupHandler(c *gin.Context) {
	// Log the incoming request
	logger.InfoLogger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"url":    c.Request.URL.String(),
	}).Info("Attempted signup")

	// Bind the incoming JSON to SignupUser struct
	var req SignupUser
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, base_utils.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid input",
			Data:    map[string]interface{}{"error": err.Error()},
		})
		return
	}

	// Call the signup function
	statusCode, message, data := Signup(req)

	// Return the appropriate JSON response
	c.JSON(statusCode, base_utils.Response{
		Status:  statusCode,
		Message: message,
		Data:    data,
	})
}

func Signup(requestInputs SignupUser) (int, string, map[string]interface{}) {
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("attempted signup method with email", requestInputs.Email, "and company", requestInputs.CompanyName)

	// Validate the input
	if err := validators.GetValidator().Struct(requestInputs); err != nil {
		return http.StatusBadRequest, "Invalid input", map[string]interface{}{"error": err.Error()}
	}

	// Check if the email already exists
	var count int64
	if err := db.DB.Model(&models.User{}).Where("email = ?", requestInputs.Email).Count(&count).Error; err != nil {
		return http.StatusInternalServerError, "Database error", map[string]interface{}{}
	}
	if count > 0 {
		return http.StatusBadRequest, "Email already exists", map[string]interface{}{}
	}

	// Check if the company already exists
	if err := db.DB.Model(&models.Company{}).Where("name = ?", requestInputs.CompanyName).Count(&count).Error; err != nil {
		return http.StatusInternalServerError, "Database error", map[string]interface{}{}
	}
	if count > 0 {
		return http.StatusBadRequest, "Company already exists", map[string]interface{}{}
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(requestInputs.Password), bcrypt.DefaultCost)
	if err != nil {
		return http.StatusInternalServerError, "Failed to hash the password", map[string]interface{}{}
	}

	// Create a new User instance
	user := models.User{
		FirstName: requestInputs.FirstName,
		LastName:  requestInputs.LastName,
		Email:     requestInputs.Email,
		Password:  string(hashedPassword),
		Company: models.Company{
			Name: requestInputs.CompanyName,
		},
	}

	// Save the user to the database
	if err := db.DB.Create(&user).Error; err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error while creating user in database")
		return http.StatusInternalServerError, "Failed to create user in database", map[string]interface{}{}
	}

	// Send a confirmation email (placeholder function)
	go base_utils.SendMail(user.Email, "subject", "body")

	logger.InfoLogger.WithFields(logrus.Fields{}).Info("User created successfully with email", requestInputs.Email, "and company", requestInputs.CompanyName)
	return http.StatusOK, "User Created Successfully", map[string]interface{}{}
}
