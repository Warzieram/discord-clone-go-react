package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"back/internal/models/user"
	"back/internal/utils"
)

type AuthResponse struct {
	Token string    `json:"token"`
	User  user.User `json:"user"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	log.Println("Registering: ", r.RemoteAddr)
	var creds user.UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		log.Println(err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if creds.Email == "" || creds.Password == "" {
		http.Error(w, "Email or password required", http.StatusBadRequest)
		return
	}

	u, err := user.CreateUser(creds.Email, creds.Password)
	if err != nil {
		http.Error(w, "Something went wrong, Try again later", http.StatusInternalServerError)
		log.Fatal("Couldn't create the user: ", err)
		return
	}

	err = u.Save()
	if err != nil {
		if strings.Contains(err.Error(), "clé dupliquée") {
			http.Error(w, "Email Already Used", http.StatusConflict)
			return
		}
		log.Println(err)
		http.Error(w, "Error creating the user", http.StatusInternalServerError)
		return

	}

	_, err = user.GetUserByEmail(u.Email)
	if err != nil {
		log.Println("Couldn't retrieve created user: ", err)
		http.Error(w, "Error registering the user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(u)
	if err != nil {
		log.Println("ERROR: Coulnd't encode response: ", err)
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	log.Println("Logging: ", r.RemoteAddr)
	var creds user.UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	user, err := user.GetUserByEmail(creds.Email)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	log.Println("User found: ", user)

	if !user.CheckPassword(creds.Password) {
		log.Println(err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	log.Println("Password match")

	if !user.EmailVerified {
		log.Println("User not verified")
		http.Error(w, "Please Verify your email", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email, user.CreatedAt)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error generating JWT", http.StatusInternalServerError)
		return
	}
	log.Println("Token Generated: ", token)

	response := AuthResponse{Token: token, User: *user}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
	}
}

func Profile(w http.ResponseWriter, r *http.Request) {

	log.Println("Getting profile for:", r.RemoteAddr)
	userID := r.Context().Value("user_id").(int)
	email := r.Context().Value("email").(string)
	createdAt := r.Context().Value("created_at")

	response := map[string]any{
		"user_id":    userID,
		"email":      email,
		"created_at": createdAt,
		"message":    "User Profile",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
