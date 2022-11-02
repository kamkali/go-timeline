package auth

import (
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/kamkali/go-timeline/internal/logger"
	"time"
)

type JWTManager struct {
	log       logger.Logger
	secretKey ed25519PrivateKey
	publicKey ed25519PubKey
}

func NewJWTManager(log logger.Logger, secret, public string) (*JWTManager, error) {
	privBlock, _ := pem.Decode([]byte(secret))
	var asn1PrivKey ed25519PrivateKey
	if _, err := asn1.Unmarshal(privBlock.Bytes, &asn1PrivKey); err != nil {
		return nil, err
	}

	pubBlock, _ := pem.Decode([]byte(public))
	var asn1PubKey ed25519PubKey
	if _, err := asn1.Unmarshal(pubBlock.Bytes, &asn1PubKey); err != nil {
		return nil, err
	}

	return &JWTManager{
		log:       log,
		secretKey: asn1PrivKey,
		publicKey: asn1PubKey,
	}, nil
}

func (j *JWTManager) GenerateToken(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(10 * time.Minute).Unix()
	claims["authorized"] = true
	claims["user"] = username

	tokenString, err := token.SignedString(j.secretKey.getPrivateKey())
	if err != nil {
		j.log.Warn(fmt.Errorf("error signing token: %w", err).Error())
		return "", fmt.Errorf("error signing token")
	}

	return tokenString, nil
}

func (j *JWTManager) VerifyToken(t string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}
		return j.publicKey.getPubKey(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("token parse: %w", err)
	}

	return token, nil
}

func (j *JWTManager) GetClaims(t string) (jwt.MapClaims, error) {
	token, err := j.VerifyToken(t)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
