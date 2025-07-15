package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/utils"
	"errors"
	"net/mail"
	"strings"
	"time"
)

// UserService provides methods for user-related operations
type UserService struct {
	Repo *repository.UserRepository // Handles database operations for users
}

// Constants defining age limits for user registration
const (
	MinAge = 13  // Minimum allowed age for registration
	MaxAge = 120 // Maximum allowed age for registration
)

// RegisterUser handles new user registration with validation
func (s *UserService) RegisterUser(user *model.User) error {
	// Validate all required fields are present and valid
	if err := s.validateRequiredFields(user); err != nil {
		return err
	}

	// Validate email format and structure
	if err := s.validateEmail(user.Email); err != nil {
		return err
	}

	// Check password meets complexity requirements
	if err := s.validatePassword(user.Password); err != nil {
		return err
	}

	// Verify user's age is within allowed range
	if err := s.validateAge(user.DOB); err != nil {
		return err
	}

	// Validate optional fields if provided
	if err := s.validateOptionalFields(user); err != nil {
		return err
	}

	// Clean and standardize input data
	s.sanitizeInput(user)

	// Check for existing users with same email or nickname
	if err := s.checkDuplicates(user); err != nil {
		return err
	}

	// Set default values for new user
	user.ID = utils.GenerateUUID() // Generate unique identifier
	if user.ProfileVisibility == "" {
		user.ProfileVisibility = "public" // Default visibility setting
	}

	// Securely hash password before storage
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = hashed // Replace plaintext with hash

	// Save validated user to database
	return s.Repo.CreateUser(user)
}

// validateRequiredFields checks all mandatory registration fields
func (s *UserService) validateRequiredFields(user *model.User) error {
	// Check email is provided and not empty
	if strings.TrimSpace(user.Email) == "" {
		return errors.New("email is required")
	}
	// Check password is provided and not empty
	if strings.TrimSpace(user.Password) == "" {
		return errors.New("password is required")
	}
	// Check first name is provided and not empty
	if strings.TrimSpace(user.FirstName) == "" {
		return errors.New("first name is required")
	}
	// Check last name is provided and not empty
	if strings.TrimSpace(user.LastName) == "" {
		return errors.New("last name is required")
	}
	// Check date of birth is provided
	if user.DOB.IsZero() {
		return errors.New("date of birth is required")
	}
	return nil
}

// validateEmail performs comprehensive email validation
func (s *UserService) validateEmail(email string) error {
	// Normalize email by trimming and lowercasing
	email = strings.TrimSpace(strings.ToLower(email))

	// Use standard library email parser
	_, err := mail.ParseAddress(email)
	if err != nil {
		return errors.New("invalid email format")
	}

	// Check maximum email length
	if len(email) > 254 {
		return errors.New("email too long (max 254 characters)")
	}

	// Basic email structure validation
	if !strings.Contains(email, "@") {
		return errors.New("invalid email format")
	}

	// Split into local and domain parts
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return errors.New("invalid email format")
	}

	localPart := parts[0]
	domain := parts[1]

	// Validate local part length
	if len(localPart) == 0 || len(localPart) > 64 {
		return errors.New("invalid email format")
	}

	// Validate domain part length
	if len(domain) == 0 || len(domain) > 255 {
		return errors.New("invalid email format")
	}

	// Domain must contain a dot
	if !strings.Contains(domain, ".") {
		return errors.New("invalid email format")
	}

	// Domain edge character validation
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") ||
		strings.HasPrefix(domain, "-") || strings.HasSuffix(domain, "-") {
		return errors.New("invalid email format")
	}

	return nil
}

// validatePassword enforces password strength rules
func (s *UserService) validatePassword(password string) error {
	// Check password length requirements
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(password) > 128 {
		return errors.New("password too long (max 128 characters)")
	}

	// Track character type requirements
	var (
		hasUpper   = false // Uppercase letter flag
		hasLower   = false // Lowercase letter flag
		hasNumber  = false // Number flag
		hasSpecial = false // Special character flag
	)

	// Define allowed character sets
	upperChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	lowerChars := "abcdefghijklmnopqrstuvwxyz"
	numberChars := "0123456789"
	specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"

	// Check each character against requirements
	for _, char := range password {
		charStr := string(char)
		switch {
		case strings.Contains(upperChars, charStr):
			hasUpper = true
		case strings.Contains(lowerChars, charStr):
			hasLower = true
		case strings.Contains(numberChars, charStr):
			hasNumber = true
		case strings.Contains(specialChars, charStr):
			hasSpecial = true
		}
	}

	// Return specific error for missing character types
	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

// validateAge checks if user's age is within allowed range
func (s *UserService) validateAge(dob time.Time) error {
	now := time.Now()
	// Calculate age in years
	age := int(now.Sub(dob).Hours() / 24 / 365.25)

	// Minimum age check
	if age < MinAge {
		return errors.New("user must be at least 13 years old")
	}

	// Maximum age check
	if age > MaxAge {
		return errors.New("invalid date of birth")
	}

	// Future date check
	if dob.After(now) {
		return errors.New("date of birth cannot be in the future")
	}

	return nil
}

// validateOptionalFields validates non-required user fields
func (s *UserService) validateOptionalFields(user *model.User) error {
	// Validate nickname if provided
	if user.Nickname != "" {
		nickname := strings.TrimSpace(user.Nickname)
		// Length validation
		if len(nickname) > 30 {
			return errors.New("nickname too long (max 30 characters)")
		}
		if len(nickname) < 3 {
			return errors.New("nickname must be at least 3 characters long")
		}

		// Character validation
		if !s.isValidNickname(nickname) {
			return errors.New("nickname can only contain letters, numbers, and underscores")
		}
	}

	// Profile visibility validation
	if user.ProfileVisibility != "" && user.ProfileVisibility != "public" && user.ProfileVisibility != "private" {
		return errors.New("profile visibility must be 'public' or 'private'")
	}

	// About section length check
	if len(user.About) > 1000 {
		return errors.New("about section too long (max 1000 characters)")
	}

	// Image URL length check
	if user.ImgURL != "" && len(user.ImgURL) > 255 {
		return errors.New("image URL too long (max 255 characters)")
	}

	return nil
}

// isValidNickname checks nickname contains only allowed characters
func (s *UserService) isValidNickname(nickname string) bool {
	// Allowed characters: letters, numbers, underscore
	allowedChars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"

	// Verify each character is allowed
	for _, char := range nickname {
		if !strings.ContainsRune(allowedChars, char) {
			return false
		}
	}
	return true
}

// sanitizeInput cleans and standardizes user input
func (s *UserService) sanitizeInput(user *model.User) {
	// Clean email: trim and lowercase
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))

	// Trim whitespace from names
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)

	// Trim whitespace from optional fields
	user.Nickname = strings.TrimSpace(user.Nickname)
	user.About = strings.TrimSpace(user.About)
	user.ImgURL = strings.TrimSpace(user.ImgURL)
}

// checkDuplicates verifies email and nickname uniqueness
func (s *UserService) checkDuplicates(user *model.User) error {
	// Check for existing email
	existingUser, err := s.Repo.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		return errors.New("email already exists")
	}

	// Check for existing nickname if provided
	if user.Nickname != "" {
		existingUser, err := s.Repo.GetUserByNickname(user.Nickname)
		if err == nil && existingUser != nil {
			return errors.New("nickname already exists")
		}
	}

	return nil
}