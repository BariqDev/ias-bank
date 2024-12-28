package token

import (
	"fmt"
	"time"

	"aidanwoods.dev/go-paseto"
)

type PasetoMaker struct {
	symmetricKey []byte
}

func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid secret key size: %d", len(symmetricKey))
	}

	return &PasetoMaker{
		[]byte(symmetricKey),
	}, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, error) {

	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}

	token := paseto.NewToken()
	token.SetExpiration(payload.ExpiredAt)
	token.SetIssuedAt(payload.IssuedAt)
	token.SetNotBefore(time.Now())
	token.SetJti(payload.ID.String())
	token.Set("claims", payload)
	symmetricKey, err := paseto.V4SymmetricKeyFromBytes(maker.symmetricKey)
	if err != nil {
		return "", err
	}
	encrypted := token.V4Encrypt(symmetricKey, nil)

	return encrypted, nil
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	symmetricKey, err := paseto.V4SymmetricKeyFromBytes(maker.symmetricKey)
	if err != nil {
		return nil, err
	}
	parser := paseto.NewParser()
	parser.AddRule(paseto.NotExpired())
	parser.AddRule(paseto.ValidAt(time.Now()))
		
	pasetoToken, err := parser.ParseV4Local(symmetricKey, token, nil)

	if err != nil {
		return nil, err
	}

	payload := &Payload{}
	err = pasetoToken.Get("claims", payload)
	
	if err != nil {
		return nil, err

	}
	return payload, nil

}
