package utils

import (
	"crypto/rand"
	"regexp"
	"unicode"

	"github.com/lindell/go-burner-email-providers/burner"
	"golang.org/x/crypto/bcrypt"
)

type RandType string

const (
	RandTypeAlphaNum             RandType = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	RandTypeAlphaNumNoSimilarity RandType = "2346789abcdefghijkmnpqrtwxyzABCDEFGHJKLMNPQRTUVWXYZ"
	RandTypeAlpha                RandType = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	RandTypeNumber               RandType = "0123456789"
)

func (rt RandType) String() string {
	return string(rt)
}

func (rt RandType) IsValid() bool {
	switch rt {
	case
		RandTypeAlpha,
		RandTypeAlphaNum,
		RandTypeNumber,
		RandTypeAlphaNumNoSimilarity:
		return true
	}
	return false
}

func RandomID(strSize int, randType ...RandType) string {
	var dictionary string

	if len(randType) == 0 {
		dictionary = RandTypeAlphaNum.String()
	}

	if len(randType) > 0 {
		dictionary = randType[0].String()
		if !randType[0].IsValid() {
			dictionary = RandTypeAlphaNum.String()
		}
	}

	bytes := make([]byte, strSize)
	rand.Read(bytes)
	for k, v := range bytes {
		bytes[k] = dictionary[v%byte(len(dictionary))]
	}
	return string(bytes)
}

// PasswordValidate: validates plain password against the rules defined below.
//
// upp: at least one upper case letter.
// low: at least one lower case letter.
// num: at least one digit.
// sym: at least one special character.
// tot: at least theRequired passwordLength as passed in arguement.
// No empty string or whitespace.
func PasswordValidate(pass string, pwdLen int) bool {
	var (
		upp, low, num, sym bool
		tot                uint8
	)

	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			upp = true
			tot++
		case unicode.IsLower(char):
			low = true
			tot++
		case unicode.IsNumber(char):
			num = true
			tot++
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			sym = true
			tot++
		default:
			return false
		}
	}

	if !upp || !low || !num || !sym || tot < uint8(pwdLen) {
		return false
	}
	return true
}

// HashPassword hashes a password to bycrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// ComparePasswordAndHash compares a given password to its bcrypt hash
func ComparePasswordAndHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// IsBurnerEmail   checks if an email is a burner email or not
func IsBurnerEmail(val string) bool {
	return burner.IsBurnerEmail(val)
}

// IsE164PhoneNumber validates a number in the E164 Format
func IsE164PhoneNumber(val string) bool {
	return regexp.MustCompile(`^\+?[1-9]\d{1,14}$`).Match([]byte(val))
}

// StringPtr returns a pointer to the string value s
func StringPtr(s string) *string {
	return &s
}

// BoolPtr returns a pointer to the bool value b
func BoolPtr(b bool) *bool {
	return &b
}

// IntPtr returns a pointer to the int value i
func IntPtr(i int) *int {
	return &i
}
