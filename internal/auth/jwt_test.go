package auth

import (
	"crypto/ed25519"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	testPubKey       = "-----BEGIN PUBLIC KEY-----\nMCowBQYDK2VwAyEAPBpEzo17xfJdvpT99CvhnU/VFgkR/KJtCEN7diBDVFY=\n-----END PUBLIC KEY-----\n"
	testSecretKey    = "-----BEGIN PRIVATE KEY-----\nMC4CAQAwBQYDK2VwBCIEIDj1z4gysEiKIHZ+SOI4guBidWpV6D8tLrqx9HE2W6Sw\n-----END PRIVATE KEY-----\n"
	testRSASecretKey = "-----BEGIN PRIVATE KEY-----\nMIIBVAIBADANBgkqhkiG9w0BAQEFAASCAT4wggE6AgEAAkEAp271VIkyePeprbsU\nI/TYchkOcCce4dnCt2zFWXKOqRRWKK5jzhZbr6hb92ltLjC2QKHaDlNbVzqaSOr/\naafv4wIDAQABAkEAiDIsAe3wTpI3RgjNo0oB3x4eroBEELeQOqCSD+atwT5mTJ8A\n8S/7JmPuNNdULmmRRUoB6KBlA2m4kii0V4QHcQIhANtgKl+NyoIXMD6CRUrqSKSN\nfvURvpizopMeDdslirPdAiEAw2LZVTQ9E7zVyqV9rbNHLlrFgc+ER/6KyKxopl+F\n1r8CIHOiaBOAGPujn3GDl2Taw7nBP+eMF+xD2/EySVl3m3odAiBv++XImeovt9lp\nDjTcK5aukMQGxKNyiAePQJGyWaliDQIgGj7UNc9+CE8vY2PtI9GXzKKawGgwDJBX\nqmamvkVMMCw=\n-----END PRIVATE KEY-----\n"
)

func TestJWTManagerGenerateToken(t *testing.T) {
	username := "test@email.com"

	t.Run("success creation", func(t *testing.T) {
		j, err := NewJWTManager(nil, testSecretKey, testPubKey)
		require.NoError(t, err)
		token, err := j.GenerateToken(username)
		require.NoError(t, err)

		got := mustParseJWTString(token, j.PubKey())
		require.NoError(t, err)
		require.Equal(t, token, got.Raw)
		require.True(t, got.Valid)

		claims, ok := got.Claims.(jwt.MapClaims)
		require.True(t, ok)
		require.Equal(t, claims["user"], username)
		require.True(t, claims["authorized"].(bool))
	})

	t.Run("invalid signing algorithm", func(t *testing.T) {
		j, err := NewJWTManager(nil, testRSASecretKey, testPubKey)
		require.NoError(t, err)
		_, err = j.GenerateToken(username)
		require.Error(t, err)
	})
}

func mustParseJWTString(token string, pubKey ed25519.PublicKey) *jwt.Token {
	t, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("invalid token signing method")
		}
		return pubKey, nil
	})
	if err != nil {
		panic(err)
	}
	return t
}

func TestJWTManagerGetClaims(t *testing.T) {
	username := "test@email.com"
	j, err := NewJWTManager(nil, testSecretKey, testPubKey)
	require.NoError(t, err)

	t.Run("successfully got claims", func(t *testing.T) {
		now := func() time.Time { return time.Date(2050, 11, 11, 22, 22, 22, 0, time.UTC) }
		j.now = now

		token, err := j.GenerateToken(username)
		require.NoError(t, err)

		got, err := j.GetClaims(token)
		require.NoError(t, err)

		require.Equal(t, got["user"], username)
		require.True(t, got["authorized"].(bool))
		require.Equal(t, float64(now().Add(10*time.Minute).Unix()), got["exp"].(float64))
		require.NoError(t, got.Valid())
	})

	t.Run("invalid token", func(t *testing.T) {
		j.now = func() time.Time { return time.Date(1970, 11, 11, 22, 22, 22, 0, time.UTC) }
		token, err := j.GenerateToken(username)
		require.NoError(t, err)

		_, err = j.GetClaims(token)
		require.Error(t, err)
	})
}

func TestJWTManagerVerifyToken(t *testing.T) {
	username := "test@email.com"
	j, err := NewJWTManager(nil, testSecretKey, testPubKey)
	require.NoError(t, err)

	token, err := j.GenerateToken(username)
	require.NoError(t, err)

	got, err := j.GetValidToken(token)
	require.NoError(t, err)

	require.Equal(t, got.Raw, token)
}
