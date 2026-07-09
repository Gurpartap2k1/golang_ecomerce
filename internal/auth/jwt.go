package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtManager struct {
	secret []byte
}

type Claims struct {
	UserId int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJwtManager(secret string) *JwtManager {
	return &JwtManager{
		secret: []byte(secret),
	}
}

func (m *JwtManager) Generate(userID int64) (string, error) {

	claims := Claims{
		UserId: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(
				time.Now().Add(24 * time.Hour),
			),
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(m.secret)

}
func (m *JwtManager) Verify(tokenString string) (*Claims, error) {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {

			return m.secret, nil

		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}

	return claims, nil
}
