package services

import (
	"fmt"
	"os"
	"time"

	"github.com/clutchtechnology/hisense-vmi-dataserver/src/models"
	"github.com/golang-jwt/jwt/v4"
)

func GetSecretKey() string {
	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "secret"
	}
	return secret
}

type AuthClaims struct {
	Identifier string                `json:"identifier"`
	Role       models.JwtServiceRole `json:"role"`
	ID         int64                 `json:"id"`
	jwt.RegisteredClaims
}

type JwtService struct {
	secretKey string
	issuer    string
}

func NewJWTService() (IJwtService, error) {
	return &JwtService{
		secretKey: GetSecretKey(),
		issuer:    "Clutch",
	}, nil
}

func (service *JwtService) GenerateToken(identifier string, id int64, role models.JwtServiceRole) string {
	claims := &AuthClaims{
		identifier,
		role,
		id,
		jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * 30)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 365)),
			Issuer:    service.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

func (service *JwtService) ValidateToken(encodedToken string, role models.JwtServiceRole) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(encodedToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("invalid token %s", token.Header["alg"])
		}
		return []byte(service.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	if claims.Role != role {
		return nil, fmt.Errorf("invalid role")
	}

	return token, nil
}
