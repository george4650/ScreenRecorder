package usecase

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-ldap/ldap/v3"
	"github.com/golang-jwt/jwt/v5"
)

type AuthUseCases struct {
	key  string
	conn *ldap.Conn
}

func NewAuthUseCases(conn *ldap.Conn, key string) *AuthUseCases {
	return &AuthUseCases{
		key:  key,
		conn: conn,
	}
}

const tokenLifetime = time.Hour * 24 * 30 //время жизни токена = 1 месяц

func (a *AuthUseCases) Authentication(login, password string) error {
	err := a.conn.Bind(login, password)
	if err != nil {
		if err, ok := err.(*ldap.Error); ok {
			if err.ResultCode == ldap.LDAPResultInvalidCredentials {
				//log.Println("invalid credentials")
				return errors.New("invalid credentials")
			}
		}
		log.Printf(fmt.Sprintf("l.Bind: %s", err.Error()))
		return fmt.Errorf("AuthUseCases - Authentication - a.conn.Bind: %w", err)
	}
	return nil
}

func (s *AuthUseCases) GenerateToken(login string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenLifetime)),
		Subject:   login,
	})

	tokenString, err := token.SignedString([]byte(s.key))
	if err != nil {
		return "", fmt.Errorf("AuthUseCases - GenerateToken - token.SignedString: %w", err)
	}
	return tokenString, nil
}

func (s *AuthUseCases) ParseToken(tokenString string) (string, error) {

	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.key), nil
	}, jwt.WithLeeway(5*time.Second))

	if err != nil {
		return "", fmt.Errorf("AuthUseCases - ParseToken -  jwt.Parse: %v", err)
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	} else {
		return "", fmt.Errorf("AuthUseCases - ParseToken -  jwt.Parse: %s", "token claims are not of type *tokenClaims")
	}

}
