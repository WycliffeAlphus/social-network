package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"time"
)

type Claims struct {
	UserID string `json:"user_id"`
	Exp    string `json:"exp"`
	Iat    string `json:"iat"`
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // Load secret from environment variable
var eat = time.FixedZone("EAT", 3*60*60)        // East Africa Time (UTC+3)

// GenerateJWT creates a JWT token for a given user ID and expiry duration.
func GenerateJWT(userID string, expiry time.Duration) (string, error) {
	now := time.Now().In(eat).Format("2006-01-02 15:04:05 MST")
	exp := time.Now().Add(expiry).In(eat).Format("2006-01-02 15:04:05 MST")
	claims := Claims{UserID: userID, Exp: exp, Iat: now}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"HS256","typ":"JWT"}`))
	payload := base64.RawURLEncoding.EncodeToString(claimsJSON)
	unsigned := header + "." + payload
	sig := sign(unsigned)
	token := unsigned + "." + sig
	return token, nil
}

// ValidateJWT parses and validates a JWT token, returning the claims if valid.
func ValidateJWT(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
	}
	// Validate header
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return nil, errors.New("invalid header encoding")
	}
	var header struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil {
		return nil, errors.New("invalid header JSON")
	}
	if header.Alg != "HS256" {
		return nil, errors.New("unsupported JWT alg")
	}
	unsigned := parts[0] + "." + parts[1]
	if !verify(unsigned, parts[2]) {
		return nil, errors.New("invalid signature")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, err
	}
	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, err
	}
	expTime, err := time.ParseInLocation("2006-01-02 15:04:05 MST", claims.Exp, eat)
	if err != nil {
		return nil, errors.New("invalid exp time format")
	}
	if time.Now().In(eat).After(expTime) {
		return nil, errors.New("token expired")
	}
	return &claims, nil
}

// sign creates an HMAC SHA256 signature for the given data using the secret.
func sign(data string) string {
	h := hmac.New(sha256.New, jwtSecret)
	h.Write([]byte(data))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))
}

// verify checks if the provided signature matches the expected signature for the data.
func verify(data, sig string) bool {
	expected := sign(data)
	return hmac.Equal([]byte(expected), []byte(sig))
}
