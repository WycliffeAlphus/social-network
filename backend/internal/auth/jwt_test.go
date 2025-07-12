package auth

import (
	"os"
	"testing"
	"time"
)

func TestGenerateAndValidateJWT_Valid(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	userID := "user123"
	token, err := GenerateJWT(userID, 1*time.Hour)
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}
	claims, err := ValidateJWT(token)
	if err != nil {
		t.Fatalf("ValidateJWT failed: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("expected user_id %s, got %s", userID, claims.UserID)
	}
}

func TestValidateJWT_Expired(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	userID := "user123"
	token, err := GenerateJWT(userID, -1*time.Hour) // already expired
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}
	_, err = ValidateJWT(token)
	if err == nil || err.Error() != "token expired" {
		t.Errorf("expected token expired error, got %v", err)
	}
}

func TestValidateJWT_InvalidSignature(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	token, err := GenerateJWT("user123", 1*time.Hour)
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}
	// Tamper with the token
	if len(token) < 10 {
		t.Fatalf("token too short to tamper")
	}
	tampered := token[:len(token)-1] + "x"
	_, err = ValidateJWT(tampered)
	if err == nil || err.Error() != "invalid signature" {
		t.Errorf("expected invalid signature error, got %v", err)
	}
}

func TestValidateJWT_InvalidHeader(t *testing.T) {
	os.Setenv("JWT_SECRET", "testsecret")
	userID := "user123"
	token, err := GenerateJWT(userID, 1*time.Hour)
	if err != nil {
		t.Fatalf("GenerateJWT failed: %v", err)
	}
	// Replace header with invalid base64
	parts := []byte(token)
	parts[0] = '!' // not base64
	_, err = ValidateJWT(string(parts))
	if err == nil {
		t.Errorf("expected invalid header encoding error, got nil")
	}
}
