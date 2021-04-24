package usecase

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type JWTCase interface {
	GenerateToken(userID string, role string) string
	ValidateToken(token string) (*jwt.Token, error)
}

type authCustomClaim struct {
	jwt.StandardClaims
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

type jwtCase struct {
	secretKey string
	issuer    string
}

func NewJWTCase() JWTCase {
	return &jwtCase{
		secretKey: getSecretKey(),
		issuer:    getIssuer(),
	}
}

func (service *jwtCase) GenerateToken(userID string, role string) string {
	claims := &authCustomClaim{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			Issuer:    service.issuer,
			IssuedAt:  time.Now().Unix(),
		},
		role,
		userID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (service *jwtCase) ValidateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("Invalid token %v", token.Header["alg"])

		}
		return []byte(service.secretKey), nil
	})

}

func getSecretKey() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "raiz"
	}
	return secret
}

func getIssuer() string {
	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "raiz"
	}
	return issuer
}
