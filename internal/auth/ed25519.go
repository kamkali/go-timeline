package auth

import (
	"crypto/ed25519"
	"encoding/asn1"
	"github.com/golang-jwt/jwt"
)

type ed25519PrivateKey struct {
	Version          int
	ObjectIdentifier struct {
		ObjectIdentifier asn1.ObjectIdentifier
	}
	PrivateKey []byte
}

type SigningMethodEd25519 struct{}

func (m *SigningMethodEd25519) Alg() string {
	return "EdDSA"
}

func (m *SigningMethodEd25519) Verify(signingString string, signature string, key any) error {
	sig, err := jwt.DecodeSegment(signature)
	if err != nil {
		return err
	}

	ed25519Key, ok := key.(ed25519.PublicKey)
	if !ok {
		return jwt.ErrInvalidKeyType
	}

	if len(ed25519Key) != ed25519.PublicKeySize {
		return jwt.ErrInvalidKey
	}

	if ok := ed25519.Verify(ed25519Key, []byte(signingString), sig); !ok {
		return jwt.ErrEd25519Verification
	}

	return nil
}

func (m *SigningMethodEd25519) Sign(signingString string, key any) (str string, err error) {
	ed25519Key, ok := key.(ed25519.PrivateKey)
	if !ok {
		return "", jwt.ErrInvalidKeyType
	}

	if len(ed25519Key) != ed25519.PrivateKeySize {
		return "", jwt.ErrInvalidKey
	}

	// Sign the string and return the encoded result
	sig := ed25519.Sign(ed25519Key, []byte(signingString))
	return jwt.EncodeSegment(sig), nil
}

type ed25519PubKey struct {
	ObjectIdentifier struct {
		ObjectIdentifier asn1.ObjectIdentifier
	}
	PublicKey asn1.BitString
}
