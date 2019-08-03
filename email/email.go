package email

import (
	"github.com/pkg/errors"
	"net/smtp"
)

// SendMail is a handler for sending out an email to an entity, reading required params from the config file
func SendMail(body string, to string) error {
	var err error
	auth := smtp.PlainAuth("", From, Pass, "smtp.gmail.com")
	// to can also be an array of addresses if needed
	msg := "From: " + From + "\n" +
		"To: " + to + "\n" +
		"Subject: OpenSolar Notification\n\n" + body

	err = smtp.SendMail("smtp.gmail.com:587", auth, From, []string{to}, []byte(msg))
	if err != nil {
		return errors.Wrap(err, "smtp error")
	}
	return nil
}
