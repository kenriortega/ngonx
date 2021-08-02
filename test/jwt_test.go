package test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gbrlsnchs/jwt/v3"
	"github.com/kenriortega/goproxy/pkg/errors"
	"github.com/kenriortega/goproxy/pkg/logger"
)

var hs = jwt.NewHS256([]byte("secret"))

// JWTPayload struct for customs payload definition
type JWTPayload struct {
	jwt.Payload
}

// Test_jwt_sign test for sing & verification of JWT
func Test_jwt_sign(t *testing.T) {
	now := time.Now()
	// Using ParseDuration() function
	experiationTime, _ := time.ParseDuration("1h")
	notBefore, _ := time.ParseDuration("30m")
	pl := JWTPayload{
		Payload: jwt.Payload{
			Issuer:         "TestJWT",
			Subject:        "System JWT Test",
			Audience:       jwt.Audience{"http://example.com:3000"},
			ExpirationTime: jwt.NumericDate(now.Add(experiationTime)),
			NotBefore:      jwt.NumericDate(now.Add(notBefore)),
			IssuedAt:       jwt.NumericDate(now),
			JWTID:          "test-server-1",
		},
	}

	token, errToken := jwt.Sign(pl, hs)
	if errToken != nil {
		t.Error(errToken.Error())
	}
	strToken := string(token)
	fmt.Println(len(strToken))
	if len(strToken) != 276 {
		t.Error("Token len expected are 267 and result are ", len(strToken))
	}

	plToVerify := JWTPayload{}
	expValidator := jwt.ExpirationTimeValidator(now)
	validatePayload := jwt.ValidatePayload(&plToVerify.Payload, expValidator)
	_, err := jwt.Verify([]byte(strToken), hs, &plToVerify, validatePayload)

	if errors.ErrorIs(err, jwt.ErrExpValidation) {
		t.Error(errors.ErrTokenExpValidation)

	}
	if errors.ErrorIs(err, jwt.ErrHMACVerification) {
		logger.LogError(errors.ErrTokenHMACValidation.Error())
		t.Error(errors.ErrTokenHMACValidation)

	}
}
