package email

import (
	"github.com/pkg/errors"
	"log"
	"net/smtp"

	"github.com/spf13/viper"
)

// SendMail is a handler for sending out an email to an entity, reading required params from the config file
func SendMail(body string, to string) error {
	var err error
	// read from config.yaml in the working directory
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "error while reading email values from config file")
	}
	log.Println("SENDING EMAIL: ", viper.Get("email"), viper.Get("password"))
	from := viper.Get("email").(string)    // interface to string
	pass := viper.Get("password").(string) // interface to string
	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")
	// to can also be an array of addresses if needed
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: OpenSolar Notification\n\n" + body

	err = smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg))
	if err != nil {
		return errors.Wrap(err, "smtp error")
	}
	return nil
}
