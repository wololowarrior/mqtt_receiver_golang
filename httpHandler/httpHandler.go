package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
)

var redisClient *redis.Client

const secret = "secret-shh"

func init() {
	redisHost := os.Getenv("REDIS_HOST")
	redisClient = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:6379", redisHost),
		DB:   0,
	})
}

type TokenClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

func generateToken(email string) (string, error) {
	claims := TokenClaims{
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret)) // Replace with your secret key
}

func postHandler(w http.ResponseWriter, r *http.Request) {
	var requestBody map[string]string
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	email, ok := requestBody["email"]
	if !ok || email == "" {
		http.Error(w, "Invalid or empty email", http.StatusBadRequest)
		return
	}

	token, err := generateToken(email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func getHandler(w http.ResponseWriter, r *http.Request) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	fmt.Printf("%v %v", claims.Email, claims.StandardClaims.ExpiresAt)

	// Retrieve the latest speed from Redis
	val, err := redisClient.Get(r.Context(), "latest_speed").Result()
	if err == redis.Nil {
		http.Error(w, "No data in Redis", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"speed": val})
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", postHandler).Methods("POST")
	r.HandleFunc("/", getHandler).Methods("GET")

	fmt.Println("Server is running on :4000")
	log.Fatal(http.ListenAndServe(":4000", r))
}
