package auth

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "tokensecret"
	expiresIn := time.Second * 3

	// Create the token
	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("error making JWT: %v", err)
	}

	if len(token) == 0 {
		t.Fatal("the string of the token is empty")
	}

	result, err := ValidateJWT(token, tokenSecret)
	if err != nil {
		t.Fatalf("error validating JWT: %v", err)
	}
	if result == uuid.Nil {
		t.Fatal("the result UUID is empty")
	}

	fmt.Println(userID)
	fmt.Println(result)

	if result != userID {
		t.Fatal("the input uuid is not equal to the result uuid")
	}

}

func TestExpiredJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "bootsismyfavoritecodingwizardbear"
	// Set a very short expiration
	expiresIn := time.Second * 1

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("error making JWT: %v", err)
	}

	// Wait for token to expire
	time.Sleep(time.Second * 2)

	// Now validate the token - error expected
	_, err = ValidateJWT(token, tokenSecret)
	if err == nil {
		t.Fatalf("expected expired token error")
	}

	if !strings.Contains(err.Error(), "token is expired") {
		t.Fatalf("expected error to contain 'token is expired', got '%v' instead", err)
	}
}

func TestInvalidSecretJWT(t *testing.T) {
	userID := uuid.New()
	tokenSecret := "tokensecret"
	wrongSecret := "thisIsTheWrongSecret"
	expiresIn := time.Second * 10

	token, err := MakeJWT(userID, tokenSecret, expiresIn)
	if err != nil {
		t.Fatalf("error making JWT: %v", err)
	}

	// Try to validate with wrong secret
	_, err = ValidateJWT(token, wrongSecret)

	if !strings.Contains(err.Error(), "signature is invalid") {
		t.Fatalf("expected signature to be invalid, 'got %v' instead", err.Error())
	}
}

func TestMalformedJWT(t *testing.T) {
	tokenSecret := "tokensecret"
	malformedToken := "this.is.not.a.valid.jwt"

	// Try to validate the malformed token
	_, err := ValidateJWT(malformedToken, tokenSecret)
	if err == nil {
		t.Fatal("expected error with malformed token")
	}

	if !strings.Contains(err.Error(), "token is malformed") {
		t.Fatalf("expected error to contain 'token is malformed', got '%v' instead", err)
	}
}
