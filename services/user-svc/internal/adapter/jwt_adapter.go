package adapter

import (
	"be-realtime-chat-app/services/commoner/utils"
	"be-realtime-chat-app/services/user-svc/internal/entity"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type JWTAdapter interface {
	GenerateAdminAccessToken(userId string) (*entity.AdminAccessToken, error)
	GenerateUserAccessToken(userId string) (*entity.UserAccessToken, error)
	VerifyAdminAccessToken(token string) (*entity.AdminAccessToken, error)
	VerifyUserAccessToken(token string) (*entity.UserAccessToken, error)
}

type jwtAdapter struct {
	adminAccessSecretByte    []byte
	employeeAccessSecretByte []byte
	adminAccessExpireTime    time.Duration
	employeeAccessExpireTime time.Duration
}

func NewJWTAdapter() JWTAdapter {
	adminAccessSecret := utils.GetEnv("ADMIN_ACCESS_TOKEN_SECRET")
	employeeAccessSecret := utils.GetEnv("EMPLOYEE_ACCESS_TOKEN_SECRET")

	adminAccessExpireStr := utils.GetEnv("ADMIN_ACCESS_TOKEN_EXP_MINUTE")
	employeeAccessExpireStr := utils.GetEnv("EMPLOYEE_ACCESS_TOKEN_EXP_MINUTE")

	adminAccessExpireInt, _ := strconv.Atoi(adminAccessExpireStr)
	employeeAccessExpireInt, _ := strconv.Atoi(employeeAccessExpireStr)

	return &jwtAdapter{
		adminAccessSecretByte:    []byte(adminAccessSecret),
		employeeAccessSecretByte: []byte(employeeAccessSecret),
		adminAccessExpireTime:    time.Duration(adminAccessExpireInt),
		employeeAccessExpireTime: time.Duration(employeeAccessExpireInt),
	}
}

func (c *jwtAdapter) GenerateAdminAccessToken(userID string) (*entity.AdminAccessToken, error) {
	expirationTime := time.Now().Add(time.Minute * c.adminAccessExpireTime)

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	stringToken, err := token.SignedString(c.adminAccessSecretByte)
	if err != nil {
		return nil, err
	}

	return &entity.AdminAccessToken{
		UserID:    userID,
		Token:     stringToken,
		ExpiresAt: expirationTime,
	}, nil
}

func (c *jwtAdapter) GenerateUserAccessToken(userID string) (*entity.UserAccessToken, error) {
	expirationTime := time.Now().Add(time.Minute * c.employeeAccessExpireTime)

	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_id"] = userID
	claims["exp"] = expirationTime.Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	stringToken, err := token.SignedString(c.employeeAccessSecretByte)
	if err != nil {
		return nil, err
	}

	return &entity.UserAccessToken{
		UserID:    userID,
		Token:     stringToken,
		ExpiresAt: expirationTime,
	}, nil
}

func (c *jwtAdapter) VerifyAdminAccessToken(token string) (*entity.AdminAccessToken, error) {
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return c.adminAccessSecretByte, nil
	})
	if err != nil {
		return nil, err
	}

	accessTokenDetail := &entity.AdminAccessToken{}
	claims, ok := tokenClaims.Claims.(jwt.MapClaims)
	if ok && tokenClaims.Valid {
		userIdStr, ok := claims["user_id"].(string)
		if !ok {
			log.Println("user_id not a string")
			return nil, fmt.Errorf("Invalid token claims")
		}

		authorized, ok := claims["authorized"].(bool)
		if !ok {
			log.Println("authorized is not a bool")
			return nil, fmt.Errorf("Invalid token claims")
		}

		if !authorized {
			log.Println("unathorize")
			return nil, fmt.Errorf("Invalid token claims")
		}

		_, err := uuid.Parse(userIdStr)
		if err != nil {
			log.Println("failed to parse ulid:", err)
			return nil, fmt.Errorf("Invalid token claims")
		}

		accessTokenDetail.UserID = userIdStr
		expFloat, ok := claims["exp"].(float64)
		if !ok {
			log.Println("exp is not a float")
			return nil, fmt.Errorf("Invalid exp in token claims")
		}

		expiresAt := time.Unix(int64(expFloat), 0)
		accessTokenDetail.ExpiresAt = expiresAt
		accessTokenDetail.Token = token
	}

	return accessTokenDetail, nil

}

func (c *jwtAdapter) VerifyUserAccessToken(token string) (*entity.UserAccessToken, error) {
	tokenClaims, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return c.employeeAccessSecretByte, nil
	})
	if err != nil {
		return nil, err
	}

	accessTokenDetail := &entity.UserAccessToken{}
	claims, ok := tokenClaims.Claims.(jwt.MapClaims)
	if ok && tokenClaims.Valid {
		userIdStr, ok := claims["user_id"].(string)
		if !ok {
			log.Println("user_id not a string")
			return nil, fmt.Errorf("Invalid token claims")
		}

		authorized, ok := claims["authorized"].(bool)
		if !ok {
			log.Println("authorized is not a bool")
			return nil, fmt.Errorf("Invalid token claims")
		}

		if !authorized {
			log.Println("unathorize")
			return nil, fmt.Errorf("Invalid token claims")
		}

		_, err := uuid.Parse(userIdStr)
		if err != nil {
			log.Println("failed to parse ulid:", err)
			return nil, fmt.Errorf("Invalid token claims")
		}

		accessTokenDetail.UserID = userIdStr
		expFloat, ok := claims["exp"].(float64)
		if !ok {
			log.Println("exp is not a float")
			return nil, fmt.Errorf("Invalid exp in token claims")
		}

		expiresAt := time.Unix(int64(expFloat), 0)
		accessTokenDetail.ExpiresAt = expiresAt
		accessTokenDetail.Token = token
	}

	return accessTokenDetail, nil

}
