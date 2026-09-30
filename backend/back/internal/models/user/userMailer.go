package user

import (
	"log"
	"os"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

var (
	publicKey  = os.Getenv("MJ_APIKEY_PUBLIC")
	privateKey = os.Getenv("MJ_APIKEY_PRIVATE")
	domainName = os.Getenv("DOMAIN_NAME")
	env        = os.Getenv("APP_ENV")
)

// verificationURL builds the link that activates an account. In dev the API is
// reached over the LAN; everywhere else it is served from DOMAIN_NAME, which
// must be the host the backend answers on (that is where /api/verify lives).
func verificationURL(token string) string {
	if env == "dev" {
		return "http://192.168.1.151:8080/api/verify?token=" + token
	}
	return "https://" + domainName + "/api/verify?token=" + token
}

// sendVerificationMail delivers a verification link to the user. With no
// Mailjet keys configured it logs the URL and returns, so a local signup can
// still be completed by hand.
func sendVerificationMail(u *User, subject string, heading string, body string) {
	url := verificationURL(u.VerificationToken.String)

	if publicKey == "" || privateKey == "" {
		log.Println("WARN: Mailjet API keys not set, skipping verification email. Verification URL:", url)
		return
	}

	mj := mailjet.NewMailjetClient(publicKey, privateKey)
	messageInfo := []mailjet.InfoMessagesV31{
		{
			From: &mailjet.RecipientV31{
				Email: "no-reply@lucramassamy.fr",
				Name:  "Luc RAMASSAMY",
			},
			To: &mailjet.RecipientsV31{
				mailjet.RecipientV31{
					Email: u.Email,
					Name:  u.Email,
				},
			},
			Subject:  subject,
			TextPart: body + "\n" + url,
			HTMLPart: "<h1>" + heading + "</h1><p>" + body + "\n" + url + "</p>",
		},
	}

	messages := mailjet.MessagesV31{Info: messageInfo}
	res, err := mj.SendMailV31(&messages)
	if err != nil {
		log.Println("ERROR: Couldn't send the mail:", err)
		return
	}

	log.Printf("Data: %+v\n", res)
}

func SendCreationEmail(u *User) {
	sendVerificationMail(
		u,
		"Welcome !",
		"Congratulations !",
		"Congratulation ! You created your account !",
	)
}

func ReSendVerificationEmail(u *User) {
	sendVerificationMail(
		u,
		"Your verification link",
		"Here is your new link",
		"Here is a fresh link to verify your account:",
	)
}
