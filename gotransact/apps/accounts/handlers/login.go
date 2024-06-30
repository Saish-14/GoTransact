package handlers

import (
	"gotransact/apps/accounts/models"
	"gotransact/apps/accounts/utils"
	base_utils "gotransact/apps/base/utils"
	validators "gotransact/apps/accounts/validators"
	"gotransact/pkg/db"
	logger "gotransact/log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

)

// Placeholder for your actual imports
// import "your_project/astuctutils"
// import "your_project/models"
// import "your_project/validators"
// import "your_project/db"

// Assuming these packages and structs are defined elsewhere in your codebase.
// var logger = logrus.New()

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	Email    string
	Password string
}


func Login(input LoginInput) (int, string, map[string]interface{}) {
	// Log the login attempt
	logger.InfoLogger.WithFields(logrus.Fields{
		"email": input.Email,
	}).Info("Attempted login")

	// Validate the input struct
	if err := validators.GetValidator().Struct(input); err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error while validating fields")
		return http.StatusBadRequest, "Error while validating fields", map[string]interface{}{}
	}

	// Fetch the user from the database
	var user models.User
	if err := db.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error finding email in database")
		return http.StatusInternalServerError, "No user found in database", map[string]interface{}{}
	}

	// Compare the provided password with the stored hashed password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Invalid password")
		return http.StatusUnauthorized, "Invalid password", map[string]interface{}{}
	}

	// Generate a PASETO token for the user
	token, err := utils.GeneratePasetoToken(user)
	if err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("Error generating token")
		return http.StatusInternalServerError, "Error generating token", map[string]interface{}{}
	}

	// Log the successful login
	logger.InfoLogger.WithFields(logrus.Fields{
		"email": input.Email,
	}).Info("Logged in successfully")

	return http.StatusOK, "Logged in successfully", map[string]interface{}{"token": token}
}


// @BasePath /api
// @Summary 			Login a user
// @Description 		User login
// @Tags 				Auth
// @Accept 				json
// @Produce 			json
// @Param 				Login body   LoginInput true "Login input"
// @in 					header
// @Success 			200 {object} base_utils.Response
// @Failure 			400 {object} base_utils.Response
// @Failure 			401 {object} base_utils.Response
// @Router 				/login [post]
func LoginHandler(c *gin.Context) {
	logger.InfoLogger.WithFields(logrus.Fields{
		"method": c.Request.Method,
		"url":    c.Request.URL.String(),
	}).Info("Attempted login")

	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, base_utils.Response{
			Status:  http.StatusBadRequest,
			Message: "Invalid input",
			Data:    map[string]interface{}{"error": err.Error()},
		})
		return
	}

	status, message, data := Login(input)
	c.JSON(status, base_utils.Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

