package utils

import (
	"gotransact/config"
	logger "gotransact/log"
	// "os"
	// "strconv"

	"github.com/sirupsen/logrus"
	gomail "gopkg.in/mail.v2"
)

// func getEmailPort() int {
// 	// Get the environment variable as a string
// 	emailPortStr := os.Getenv("EMAIL_PORT")

// 	// Convert the string to an integer
// 	emailPort, err := strconv.Atoi(emailPortStr)
// 	if err != nil {
// 		logger.ErrorLogger.Fatalf("Error converting EMAIL_PORT to integer: %v", err)
// 	}

// 	return emailPort
// }

// var EmailPort = getEmailPort()

// "user has been created successfully ,this is a confirmation mail"
func SendMail(email string, subject string, body string) {
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("attempted Sendmail() email to", email)
	mail := gomail.NewMessage()

	mail.SetHeader("From", config.FromEmail)
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/html", body)

	dialer := gomail.NewDialer(config.EmailHost, 587, config.EmailUser, config.EmailPass)
	if err := dialer.DialAndSend(mail); err != nil {
		logger.ErrorLogger.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Error("error sending mail")
		logger.ErrorLogger.Fatal(err.Error())
	}
	logger.InfoLogger.WithFields(logrus.Fields{}).Info("mail sent successfully to ", email)
}

func SendMailWithAttachment(email, subject, body, filePath string) {
	m := gomail.NewMessage()
	m.SetHeader("From", config.FromEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)
	m.Attach(filePath)


	d := gomail.NewDialer(config.EmailHost, 587, config.EmailUser, config.EmailPass)

	if err := d.DialAndSend(m); err != nil {
		logger.ErrorLogger.Printf("could not send email: %v", err)
	}
	logger.InfoLogger.Printf("Email sent successfully")
}


