package handlers

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"encoding/base64"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/golang-jwt/jwt"
)

type Auth struct {
	db  *sql.DB
	key []byte
}

func NewAuth() (*Auth, error) {
	dsn := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("AUTH_DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	encodedKey := os.Getenv("JWT_KEY")
	key, err := base64.StdEncoding.DecodeString(encodedKey)
	if err != nil {
		return nil, err
	}

	return &Auth{db, key}, nil
}

type NonceReq struct {
	Email string
}

type NonceResp struct {
	Nonce []byte
}

func (a *Auth) NonceRoute(w http.ResponseWriter, r *http.Request) {
	req, ok := recieveJSON[NonceReq](w, r)
	if !ok {
		fmt.Println("failed to recieve email")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	nonce, err := a.Nonce(req.Email)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sendJSON(NonceResp{Nonce: nonce}, w)
}

type TokenReq struct {
	Email string
	Proof []byte
}

type TokenResp struct {
	Email    string
	ID       int
	Role     string
	External bool
	Access   string
	Refresh  string
}

func (a *Auth) TokensRoute(w http.ResponseWriter, r *http.Request) {
	req, ok := recieveJSON[TokenReq](w, r)
	if !ok {
		return
	}
	id, access, refresh, role, external, err := a.Tokens(req.Email, req.Proof)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	sendJSON(TokenResp{ID: id, Access: access, Refresh: refresh, Role: role, External: external}, w)
}

type Claims struct {
	ID   int    `json:"id"`
	Kind string `json:"kind"`
	jwt.StandardClaims
}

type User struct {
	ID       int
	Email    string
	Role     string
	External bool
	Merchant bool
	Company  int
}

const (
	accessKind  = "access"
	refreshKind = "refresh"
)

type UserKey string

const userKey UserKey = "user"

func AddUserToContext(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func UserFromContext(ctx context.Context) (*User, bool) {
	user, ok := ctx.Value(userKey).(*User)
	return user, ok
}

func (a *Auth) Nonce(email string) ([]byte, error) {
	nonce := make([]byte, 64)
	rand.Read(nonce)

	var (
		id     int
		secret []byte
	)
	tx, err := a.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = tx.QueryRow(`
		SELECT
			id,
			hash
		FROM users
		WHERE email = $1
		`, email).Scan(&id, &secret)
	if err != nil {
		return nil, fmt.Errorf("user may not exist: %w", err)
	}
	if len(secret) == 0 {
		return nil, fmt.Errorf("secret has no length")
	}

	salted := slices.Concat(nonce, secret)
	salted = append(salted, nonce...)
	hashed := sha512.Sum512(salted)

	_, err = tx.Exec(`
		DELETE FROM challenge
		WHERE user_id = $1`, id)
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(`
		INSERT INTO challenge
			(user_id, hash)
		VALUES
			($1, $2)
		`, id, hashed[:])
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return nonce, nil
}

func (a *Auth) Tokens(email string, hash []byte) (int, string, string, string, bool, error) {
	tx, err := a.db.Begin()
	if err != nil {
		return 0, "", "", "", true, err
	}
	defer tx.Rollback()

	var (
		id        int
		role      string
		external  bool
		challenge []byte
	)
	err = tx.QueryRow(`
		SELECT
			users.id,
			users.role,
			companies.external,
			challenge.hash
		FROM users
		JOIN challenge
		ON users.id = challenge.user_id
		JOIN companies
		ON companies.id = users.company_id
		WHERE users.email = $1
		`, email).Scan(&id, &role, &external, &challenge)
	if err != nil {
		return 0, "", "", "", true, err
	}

	if !bytes.Equal(challenge, hash) {
		return 0, "", "", "", true, fmt.Errorf("failed to provide proof")
	}

	accessToken, refreshToken, err := a.makeTokens(id, tx)
	if err != nil {
		return 0, "", "", "", true, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, "", "", "", true, err
	}
	return id, accessToken, refreshToken, role, external, nil
}

func (a *Auth) makeTokens(id int, tx *sql.Tx) (string, string, error) {
	now := time.Now()
	accessExpires := now.Add(time.Hour * 72)
	refreshExpires := now.Add(time.Hour * 24 * 14)

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, &Claims{
		ID:   id,
		Kind: accessKind,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: accessExpires.Unix(),
		},
	}).SignedString(a.key)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS512, &Claims{
		ID:   id,
		Kind: refreshKind,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: refreshExpires.Unix(),
		},
	}).SignedString(a.key)
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}

	_, err = tx.Exec(`
		INSERT INTO tokens (
			user_id,
			fingerprint,
			access_expiry,
			refresh_expiry,
			access_token,
			refresh_token
		) VALUES ($1, '', $2, $3, $4, $5)
		`, id, accessExpires, refreshExpires, accessToken, refreshToken)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (a *Auth) Check(accessToken string) (*User, bool, error) {
	token := strings.TrimPrefix(accessToken, "Bearer ")

	t, err := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (any, error) {
		_, ok := t.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return a.key, nil
	})
	if err != nil {
		slog.Warn("Errored claims", "Error", err, "Token", accessToken)
		return nil, false, err
	}

	claims, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		slog.Warn("Bad claims", "Token", accessToken)
		return nil, false, nil
	}

	if claims.ExpiresAt < time.Now().Unix() {
		slog.Warn("Claims expired", "Token", accessToken)
		return nil, false, nil
	}

	var u User
	err = a.db.QueryRow(`
		SELECT
			users.id,
			users.email,
			users.role,
			companies.external,
			companies.merchant,
			companies.id
		FROM users
		JOIN companies ON users.company_id = companies.id
		JOIN tokens ON tokens.user_id = users.id
		WHERE
			tokens.access_token = $1 AND
			users.active = true AND
			tokens.access_expiry > current_timestamp
		`, token).Scan(&u.ID, &u.Email, &u.Role, &u.External, &u.Merchant, &u.Company)
	if err != nil {
		return nil, false, err
	}

	if u.ID != claims.ID {
		slog.Warn("Claims don't match", "claims", claims, "user", u)
		return nil, false, nil
	}

	return &u, true, nil
}
