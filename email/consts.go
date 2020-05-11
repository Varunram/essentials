package email

// From is the source of the email
var From string

// Pass is the password of the sender's email account
var Pass string

// SetConsts sets the consts for email processing
func SetConsts(from string, pass string) {
	From = from
	Pass = pass
}
