package task

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Ararat25/go_final_project/cmd/config"
	"github.com/Ararat25/go_final_project/errors"
	"github.com/Ararat25/go_final_project/model/entity"
	"github.com/golang-jwt/jwt/v5"
)

var serviceInstance *Service

// Service структура для хранения данных для работы с сервисом
type Service struct {
	db       Storage
	TokenTTL time.Duration
	Config   *config.Config
}

// NewService создание нового объекта Service
func NewService(db *Storage, tokenTTL time.Duration, config *config.Config) *Service {
	if serviceInstance == nil {
		serviceInstance = &Service{
			db:       *db,
			TokenTTL: tokenTTL,
			Config:   config,
		}
	}

	return serviceInstance
}

func (s *Service) CheckToken(token, envPassword string) error {
	var claims entity.TokenClaims

	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("incorrect method")
		}

		return s.Config.TokenSalt, nil
	})
	if err != nil {
		return err
	}

	if !jwtToken.Valid {
		return errors.ErrNotValidJwtToken
	}

	if !s.doPasswordsMatch(claims.PasswordChecksum, envPassword) {
		return errors.ErrNotValidJwtToken
	}

	return nil
}

func (s *Service) VerifyUser(envPassword, password string) (string, error) {
	if envPassword != password {
		return "", errors.ErrPasswordNotValid
	}

	return s.generateJwtToken(password)
}

func (s *Service) generateJwtToken(password string) (string, error) {
	now := time.Now()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, entity.TokenClaims{
		PasswordChecksum: s.checksum(password),
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.TokenTTL)),
		},
	})

	signedToken, err := jwtToken.SignedString(s.Config.TokenSalt)
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
