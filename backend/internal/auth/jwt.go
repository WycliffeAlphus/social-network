package auth

import (
	"errors"
	"time"

	// "os"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
)

type Claims struct {
	UserID string `json:"user_id"`
	Exp    int64  `json:"exp"`
}

var jwtSecret = []byte("changeme") // TODO: load from env

// GenerateJWT creates a JWT token for a given user ID and expiry duration.
func GenerateJWT(userID string, expiry time.Duration) (string, error) {
	exp := time.Now().Add(expiry).Unix()
	claims := Claims{UserID: userID, Exp: exp}
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
	parts := split(token, '.')
	if len(parts) != 3 {
		return nil, errors.New("invalid token format")
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
	if time.Now().Unix() > claims.Exp {
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

// split splits a string s by the separator sep and returns a slice of substrings.
func split(s string, sep rune) []string {
	var out []string
	start := 0
	for i, c := range s {
		if c == sep {
			out = append(out, s[start:i])
			start = i + 1
		}
	}
	out = append(out, s[start:])
	return out
}
