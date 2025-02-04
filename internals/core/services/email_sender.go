package services

import "github.com/go-gomail/gomail"

type EmailService struct {
	Host     string
	Port     int
	Email    string
	Password string
}

// sample input = {"smtp.gmail.com", 587, "inecrft4747@gmail.com", "nyqf jefi yuga rcgc"}
func NewEmailService(host string, port int, email string, password string) *EmailService {
	return &EmailService{Host: host, Port: port, Email: email, Password: password}
}

// call this method in user registration function
// token: some string that can be use to Get which user it is from later
func (e *EmailService) SendVerificationEmail(email, token string) error {
	message := gomail.NewMessage()
	message.SetHeader("From", e.Email)
	message.SetHeader("To", email)
	message.SetHeader("Subject", "ConDormHub Email Verification")

	// have not yet implement API for /verify
	message.SetBody("text/plain", "Click the link to verify your account: http://localhost:3000/verify/"+token)

	dailer := gomail.NewDialer(e.Host, e.Port, e.Email, e.Password)

	return dailer.DialAndSend(message)
}
