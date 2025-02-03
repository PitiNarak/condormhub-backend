package email

import "github.com/go-gomail/gomail"

func SendEmail(email, token string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", "inecrft4747@gmail.com")
	message.SetHeader("To", email)
	message.SetHeader("Subject", "Hello Email!")

	message.SetBody("text/plain", "Test Body")

	dailer := gomail.NewDialer("smtp.gmail.com", 587, "inecrft4747@gmail.com", "nyqf jefi yuga rcgc")

	return dailer.DialAndSend(message)
}
