package auth

import (
	"crypto/ed25519"
	"encoding/asn1"
	"encoding/pem"
	"fmt"
	"github.com/golang-jwt/jwt"
	"go.uber.org/zap"
	"time"
)

type JWTManager struct {
	log       *zap.Logger
	now       func() time.Time
	secretKey ed25519PrivateKey
	publicKey ed25519PubKey
}

func NewJWTManager(log *zap.Logger, secret, public string) (*JWTManager, error) {
	privBlock, _ := pem.Decode([]byte(secret))
	if privBlock == nil {
		return nil, fmt.Errorf("invalid secret key format")
	}
	var asn1PrivKey ed25519PrivateKey
	if _, err := asn1.Unmarshal(privBlock.Bytes, &asn1PrivKey); err != nil {
		return nil, err
	}

	pubBlock, _ := pem.Decode([]byte(public))
	if pubBlock == nil {
		return nil, fmt.Errorf("invalid public key format")
	}
	var asn1PubKey ed25519PubKey
	if _, err := asn1.Unmarshal(pubBlock.Bytes, &asn1PubKey); err != nil {
		return nil, err
	}

	return &JWTManager{
		log:       log,
		now:       time.Now,
		secretKey: asn1PrivKey,
		publicKey: asn1PubKey,
	}, nil
}

func (j *JWTManager) PrivateKey() ed25519.PrivateKey {
	return ed25519.NewKeyFromSeed(j.secretKey.PrivateKey[2:])
}

func (j *JWTManager) PubKey() ed25519.PublicKey {
	return j.publicKey.PublicKey.Bytes
}

func (j *JWTManager) GenerateToken(username string) (t string, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("SigningString panicked, recover: %v", r)
		}
	}()
	token := jwt.New(jwt.SigningMethodEdDSA)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = j.now().Add(10 * time.Minute).Unix()
	claims["authorized"] = true
	claims["user"] = username

	tokenString, err := token.SignedString(j.PrivateKey())
	if err != nil {
		j.log.Warn(fmt.Errorf("error signing token: %w", err).Error())
		return "", fmt.Errorf("error signing token")
	}

	return tokenString, nil
}

func (j *JWTManager) GetValidToken(t string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}
		return j.PubKey(), nil
	})
	if err != nil {
		return nil, fmt.Errorf("token parse: %w", err)
	}

	return token, nil
}

func (j *JWTManager) GetClaims(t string) (jwt.MapClaims, error) {
	token, err := j.GetValidToken(t)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
