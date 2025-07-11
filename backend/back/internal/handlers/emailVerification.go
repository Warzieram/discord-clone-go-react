package handlers

import (
	"back/internal/models/user"
	"log"
	"net/http"
)

func VerifyEmail(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	if token != "" {
		user, err := user.GetUserByVerificationToken(token)
		if err != nil {
			log.Println("Error retriving user from token: ", err)
			http.Error(w, "Something went wrong", http.StatusBadRequest)
			return
		}
		err = user.VerifyEmail()
		if err != nil {
			if err.Error() == "email already verified" {
				log.Println("ERROR: ", err)
				http.Error(w, "Email is already verified", http.StatusBadRequest)
				return
			}
			if err.Error() == "verification token expired" {
				log.Println("ERROR: ", err)
				http.Error(w, "Verification Token is expired", http.StatusBadRequest)
				return
			}
			log.Println("ERROR: ", err)
			http.Error(w, "Couldn't verify email", http.StatusInternalServerError)
			return
		}
	} else {
		log.Println("param token not present")
		http.Error(w, "Missing token param", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "http://192.168.1.151:5173/login", http.StatusMovedPermanently)

}

func ResendVerification(w http.ResponseWriter, r *http.Request) {

}
