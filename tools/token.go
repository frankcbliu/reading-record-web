package tools

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Keep this two config private, it should not expose to open source
const NBSecretPassword = "CBNB this#is@a very=strong password.szu"
const NBRandomPassword = "A String Very Very Very Niubilty!!@##$!@#4"

// A Util function to generate jwt_token which can be used in the request header
func GenToken(user_name string) string {
	jwt_token := jwt.New(jwt.GetSigningMethod("HS256"))
	// Set some claims
	jwt_token.Claims = jwt.MapClaims{
		"user_name": user_name,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}
	// Sign and get the complete encoded token as a string
	token, _ := jwt_token.SignedString([]byte(NBSecretPassword))
	return token
}
