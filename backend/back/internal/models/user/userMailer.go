package user

import (
	"log"
	"os"

	"github.com/mailjet/mailjet-apiv3-go/v4"
)

var (
	publicKey  = os.Getenv("MJ_APIKEY_PUBLIC")
	privateKey = os.Getenv("MJ_APIKEY_PRIVATE")
)

func SendCreationEmail(u *User) {
	tokenString := u.VerificationToken.String

	verificationURL := "http://192.168.1.151:8080/api/verify?token=" + tokenString

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
			Subject:  "Welcome !",
			TextPart: "Congratulation ! You created your account !\n" + verificationURL,
			HTMLPart: "<h1>Congratulations !</h1><p>You created your account\n " + verificationURL + "</p>",
		},
	}
	messages := mailjet.MessagesV31{Info: messageInfo}
	res, err := mj.SendMailV31(&messages)
	if err != nil {
		log.Fatal("Couldn't send the mail :", err)
	}

	log.Printf("Data: %+v\n", res)

}

func ReSendVerificationEmail(u *User) {
	tokenString := u.VerificationToken.String
	verificationURL := "http://192.168.1.151:8080/api/verify?token=" + tokenString

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
			Subject:  "Welcome !",
			TextPart: "Congratulation ! You created your account !\n" + verificationURL,
			HTMLPart: "<h1>Congratulations !</h1><p>You created your account\n" + verificationURL + "</p>",
		},
	}
	log.Println(mj)
	log.Println(messageInfo)

	return

}
