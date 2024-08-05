package model

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/Ararat25/go_final_project/customError"
	"github.com/Ararat25/go_final_project/dbManager"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// Service структура для хранения ссылки на бд
type Service struct {
	DB        *dbManager.SchedulerStore
	TokenTTL  time.Duration
	TokenSalt []byte
}

// NewService создание нового объекта Service
func NewService(db *dbManager.SchedulerStore, tokenTTL time.Duration, tokenSalt []byte) *Service {
	return &Service{
		DB:        db,
		TokenTTL:  tokenTTL,
		TokenSalt: tokenSalt,
	}
}

func (s *Service) CheckToken(token, envPassword string) error {
	var claims TokenClaims

	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("incorrect method")
		}

		return s.TokenSalt, nil
	})
	if err != nil {
		return err
	}

	if !jwtToken.Valid {
		return customError.ErrNotValidJwtToken
	}

	if !s.doPasswordsMatch(claims.PasswordChecksum, envPassword) {
		return customError.ErrNotValidJwtToken
	}

	return nil
}

func (s *Service) VerifyUser(envPassword, password string) (string, error) {
	if envPassword != password {
		return "", customError.ErrPasswordNotValid
	}

	return s.generateJwtToken(password)
}

func (s *Service) generateJwtToken(password string) (string, error) {
	now := time.Now()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, TokenClaims{
		PasswordChecksum: s.checksum(password),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.TokenTTL)),
		},
	})

	signedToken, err := jwtToken.SignedString(s.TokenSalt)
	if err != nil {
		return "", err
	}

	fmt.Printf("Токен для входа и тестов: %s", signedToken)

	return signedToken, nil
}

// checksum создает контрольную сумму для проверки токена
func (s *Service) checksum(password string) string {
	result := sha256.Sum256([]byte(password))

	hashString := hex.EncodeToString(result[:])

	return hashString
}

// doPasswordsMatch сравнивает контрольную сумму паролей
func (s *Service) doPasswordsMatch(checksum, password string) bool {
	return checksum == s.checksum(password)
}
