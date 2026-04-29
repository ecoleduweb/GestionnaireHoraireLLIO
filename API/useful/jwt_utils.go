package useful

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// VerifyTokenValidity vérifie si le jeton JWT est valide selon les dates
// fournies dans ses *claims*. Utile pour valider un token d'accès avant
// de faire un appel vers un API externe.
//
// P.S. 1 : Ce code NE VÉRIFIE PAS L'AUTHENTICITÉ DU JETON.
//
// P.S. 2 : Un jeton sans date d'expiration sera automatiquement rejeté.
func VerifyTokenValidity(tokenString string) bool {
	parser := jwt.NewParser()
	token, _, err := parser.ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return false
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false
	}

	// Obtention de la date d'expiration / nil si non définie.
	exp, err := claims.GetExpirationTime()
	if err != nil || exp == nil {
		return false
	}

	return exp.After(time.Now())
}
