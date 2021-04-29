package helper

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gopkg.in/gomail.v2"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587

var (
	CONFIG_AUTH_EMAIL    string = os.Getenv("EMAIL_ADDRESS")
	CONFIG_AUTH_PASSWORD string = os.Getenv("EMAIL_PASS")
	CONFIG_SENDER_NAME   string = "PT. MY COMPANY <" + CONFIG_AUTH_EMAIL + ">"
)

func SendMail(to string, subject string, body string) {
	// log.Println("GOROUTINES")
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Println(err.Error())
	} else {

		log.Println("Mail sent!")
	}
	// log.Println(CONFIG_AUTH_EMAIL)
	// log.Println(CONFIG_AUTH_PASSWORD)
	// log.Println(CONFIG_SENDER_NAME)
}
