package authorizer

// docs for jwkfetch: https://godoc.org/github.com/Soluto/fetch-jwk
// docs for jwt: https://godoc.org/github.com/dgrijalva/jwt-go
// docs for lambda authorizer: https://github.com/awslabs/aws-apigateway-lambda-authorizer-blueprints

import (
	"errors"
	"fmt"

	jwkfetch "github.com/Soluto/fetch-jwk"
	jwt "github.com/dgrijalva/jwt-go"
)

const (
	clientID = "2084ukslsc831pt202t2dudt7c"
)

// Validate function
// Other than the following: algorithm, expiry, and token structure, the only other thing
// we're checking is the Cognito clientID
func Validate(cognitoClientID, tokenString string) (principalID string, err error) {

	var errStr string

	jwk := jwkfetch.FromIssuerClaim
	token, err := jwt.Parse(tokenString, jwk())
	claims := token.Claims.(jwt.MapClaims)

	if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			return principalID, errors.New("Invalid token")
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			// Token is either expired or not active yet
			return principalID, errors.New("Expired token")
		} else {
			return principalID, errors.New("Invalid token with unknown type")
		}
	}

	// Validate the expected alg
	if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
		errStr = fmt.Sprintf("unexpected signing method: %v", token.Header["alg"])
		return principalID, errors.New(errStr)
	}

	// Check clientID
	if cognitoClientID != claims["client_id"] {
		return principalID, errors.New("Invalid client id")
	}

	if claims["username"] == "" {
		return principalID, errors.New("Missing username")
	}
	principalID = fmt.Sprintf("%s|%s", claims["username"], claims["client_id"])

	return principalID, err
}
