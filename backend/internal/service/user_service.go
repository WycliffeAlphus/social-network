package service

import (
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/utils"
	"fmt"
	"log"
	"net/mail"
	"strings"
	"time"
)

type RegistrationErrors struct {
	Email       string
	Nickname    string
	Password    string
	FirstName   string
	LastName    string
	DateOfBirth string
	Avatar      string
	AboutMe     string
	Visibility  string
}

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
func (s *UserService) RegisterUser(user *model.User) (*RegistrationErrors, error) {
	errors := &RegistrationErrors{}

	// Validate all required fields are present and valid
	s.validateRequiredFields(user, errors)

	// Validate email format and structure
	s.validateEmail(user.Email, errors)

	// Check password meets complexity requirements
	s.validatePassword(user.Password, errors)

	// Verify user's age is within allowed range
	s.validateAge(user.DOB, errors)

	// Validate optional fields if provided
	s.validateOptionalFields(user, errors)

	// Clean and standardize input data
	s.sanitizeInput(user)

	// Check for existing users with same email or nickname
	s.checkDuplicates(user, errors)

	if errors.HasErrors() {
		return errors, nil
	}

	// Set default values for new user
	user.ID = utils.GenerateUUID() // Generate unique identifier
	if user.ProfileVisibility == "" {
		user.ProfileVisibility = "public" // Default visibility setting
	}

	// Securely hash password before storage
	hashed, err := utils.HashPassword(user.Password)
	if err != nil {
		log.Println("error while hashing password during registration", err)
	}
	user.Password = hashed // Replace plaintext with hash

	// Save validated user to database
	return nil, s.Repo.CreateUser(user)
}

// validateRequiredFields checks all mandatory registration fields
func (s *UserService) validateRequiredFields(user *model.User, errors *RegistrationErrors) {
	// Check email is provided and not empty
	if strings.TrimSpace(user.Email) == "" {
		errors.Email = "Email is required"
		return
	}
	// Check password is provided and not empty
	if strings.TrimSpace(user.Password) == "" {
		errors.Password = "Password is required"
		return
	}
	// Check first name is provided and not empty
	if strings.TrimSpace(user.FirstName) == "" {
		errors.FirstName = "First name is required"
		return
	}
	// Check last name is provided and not empty
	if strings.TrimSpace(user.LastName) == "" {
		errors.LastName = "Last name is required"
		return
	}
	// Check date of birth is provided
	if user.DOB.IsZero() {
		errors.DateOfBirth = "Date of birth is required"
		return
	}
}

// validateEmail performs comprehensive email validation
func (s *UserService) validateEmail(email string, errors *RegistrationErrors) {
	// Normalize email by trimming and lowercasing
	email = strings.TrimSpace(strings.ToLower(email))

	// Use standard library email parser
	_, err := mail.ParseAddress(email)
	if err != nil {
		errors.Email = "Invalid email format"
		return
	}

	// Check maximum email length
	if len(email) > 254 {
		errors.Email = "Email too long (max 254 characters)"
		return
	}

	// Basic email structure validation
	if !strings.Contains(email, "@") {
		errors.Email = "Invalid email format"
		return
	}

	// Split into local and domain parts
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		errors.Email = "Invalid email format"
		return
	}

	localPart := parts[0]
	domain := parts[1]

	// Validate local part length
	if len(localPart) == 0 || len(localPart) > 64 {
		errors.Email = "Invalid email format"
		return
	}

	// Validate domain part length
	if len(domain) == 0 || len(domain) > 255 {
		errors.Email = "Invalid email format"
		return
	}

	// Domain must contain a dot
	if !strings.Contains(domain, ".") {
		errors.Email = "Invalid email format"
		return
	}

	// Domain edge character validation
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") ||
		strings.HasPrefix(domain, "-") || strings.HasSuffix(domain, "-") {
		errors.Email = "Invalid email format"
		return
	}
}

// validatePassword enforces password strength rules
func (s *UserService) validatePassword(password string, errors *RegistrationErrors) {
	// Check password length requirements
	if len(password) < 8 {
		errors.Password = "Password must be at least 8 characters long"
		return
	}
	if len(password) > 16 {
		errors.Password = "Password too long (max 16 characters)"
		return
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
		errors.Password = "Password must contain at least one uppercase letter"
		return
	}
	if !hasLower {
		errors.Password = "Password must contain at least one lowercase letter"
		return
	}
	if !hasNumber {
		errors.Password = "Password must contain at least one number"
		return
	}
	if !hasSpecial {
		errors.Password = "Password must contain at least one special character"
		return
	}
}

// validateAge checks if user's age is within allowed range
func (s *UserService) validateAge(dob time.Time, errors *RegistrationErrors) {
	now := time.Now()

	// Future date check
	if dob.After(now) {
		errors.DateOfBirth = "Date of birth cannot be in the future"
		return
	}

	// Calculate age in years
	age := int(now.Sub(dob).Hours() / 24 / 365.25)

	// Minimum age check
	if age < MinAge {
		errors.DateOfBirth = fmt.Sprintf("User must be at least %d years old", MinAge)
		return
	}

	// Maximum age check
	if age > MaxAge {
		errors.DateOfBirth = fmt.Sprintf("User must be at most %d years old", MaxAge)
		return
	}
}

// validateOptionalFields validates non-required user fields
func (s *UserService) validateOptionalFields(user *model.User, errors *RegistrationErrors) {
	// Validate nickname if provided
	if user.Nickname != "" {
		nickname := strings.TrimSpace(user.Nickname)
		// Length validation
		if len(nickname) > 30 {
			errors.Nickname = "Nickname too long (max 30 characters)"
			return
		}
		if len(nickname) < 3 {
			errors.Nickname = "Nickname must be at least 3 characters long"
			return
		}

		// Character validation
		if !s.isValidNickname(nickname) {
			errors.Nickname = "Nickname can only contain letters, numbers, and underscores"
			return
		}
	}

	// Profile visibility validation
	if user.ProfileVisibility != "" && user.ProfileVisibility != "public" && user.ProfileVisibility != "private" {
		errors.Visibility = "Profile visibility must be 'public' or 'private'"
		return
	}

	// About section length check
	if len(user.About) > 1000 {
		errors.AboutMe = "About section too long (max 1000 characters)"
		return
	}

	// Image URL length check
	if user.ImgURL != "" && len(user.ImgURL) > 255 {
		errors.Avatar = "Image URL too long (max 255 characters)"
		return
	}
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
func (s *UserService) checkDuplicates(user *model.User, errors *RegistrationErrors) {
	// Check for existing email
	emailExists := repository.GetUserByEmail(s.Repo.DB, user.Email)
	if emailExists {
		errors.Email = "Email already exists"
		return
	}

	// Check for existing nickname if provided
	if user.Nickname != "" {
		nicknameExists := repository.GetUserByNickname(s.Repo.DB, user.Nickname)
		if nicknameExists {
			errors.Nickname = "Nickname already exists"
			return
		}
	}
}

func (re *RegistrationErrors) HasErrors() bool {
	return re.Email != "" ||
		re.Nickname != "" ||
		re.Password != "" ||
		re.FirstName != "" ||
		re.LastName != "" ||
		re.DateOfBirth != "" ||
		re.Avatar != "" ||
		re.AboutMe != "" ||
		re.Visibility != ""
}
