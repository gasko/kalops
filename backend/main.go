package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB
var supabaseProject = "cjksralldmtlpbubrwjs"
var supabaseJWKSURL = "https://" + supabaseProject + ".supabase.co/auth/v1/.well-known/jwks.json"

// JWKS cache
var jwks *keyfunc.JWKS

func initJWKS() {
	var err error
	fmt.Println(supabaseJWKSURL)
	jwks, err = keyfunc.Get(supabaseJWKSURL, keyfunc.Options{
		RefreshInterval: 5 * time.Minute,
		RefreshErrorHandler: func(err error) {
			log.Printf("error refreshing JWKS: %v", err)
		},
	})
	if err != nil {
		log.Fatalf("failed to get JWKS: %v", err)
	}
}

func verifyJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, jwks.Keyfunc)
}

/*func handler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := verifyJWT(tokenStr)
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Exempel på query mot lokala Postgres 17
	var now string
	err = db.QueryRow("SELECT NOW()").Scan(&now)
	if err != nil {
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Hello! Database time (Postgres 17): %s", now)
}*/

func handler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing token", http.StatusUnauthorized)
		return
	}
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	token, err := verifyJWT(tokenStr) // kollar mot Supabase JWKS
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Hej %v! 🎉", token.Claims.(jwt.MapClaims)["email"])
}
func main() {
	// init DB
	var err error
	db, err = sql.Open("pgx", "postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// testa DB-connection
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("could not connect to postgres: %v", err)
	}

	// init JWKS (Supabase keys)
	initJWKS()

	http.HandleFunc("/api/hello", handler)
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
