package jwt

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type JWT struct {
	accessSecret  string
	accessExpire  time.Duration
	refreshSecret string
	refreshExpire time.Duration
}

type SignedToken = string

type accessClaims struct {
	jwt.RegisteredClaims
}

type refreshClaims struct {
	jwt.RegisteredClaims
}

type AccessPayload struct {
	userID uuid.UUID
}

type RefreshPayload struct {
}

type Access struct {
	Token  SignedToken `json:"refresh"`
	claims accessClaims
}

type Refresh struct {
	Token  SignedToken `json:"access"`
	claims refreshClaims
}

type TokenPair struct {
	Access
	Refresh
}

func New(aSecret string, aExpire time.Duration, rSecret string, rExpire time.Duration) JWT {
	return JWT{
		accessSecret:  aSecret,
		accessExpire:  aExpire,
		refreshSecret: rSecret,
		refreshExpire: rExpire,
	}
}

func (j JWT) NewPair(userID uuid.UUID) (TokenPair, error) {
	ac, err := j.newAccess(userID)
	if err != nil {
		return TokenPair{}, err
	}
	rf, err := j.newRefresh(userID)
	if err != nil {
		return TokenPair{}, err
	}
	return TokenPair{ac, rf}, nil
}

func (j JWT) newAccess(userID uuid.UUID) (Access, error) {
	cl := accessClaims{
		jwt.RegisteredClaims{
			Issuer:    userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	raw := jwt.NewWithClaims(jwt.SigningMethodHS512, cl)
	t, err := raw.SignedString([]byte(j.accessSecret))
	if err != nil {
		return Access{}, err
	}
	return Access{t, cl}, nil
}

func (j JWT) newRefresh(userID uuid.UUID) (Refresh, error) {
	cl := refreshClaims{
		jwt.RegisteredClaims{
			Issuer:    userID.String(),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshExpire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	raw := jwt.NewWithClaims(jwt.SigningMethodHS512, cl)
	t, err := raw.SignedString([]byte(j.refreshSecret))
	if err != nil {
		return Refresh{}, err
	}
	return Refresh{t, cl}, nil
}

func (j JWT) ParseAccess(t SignedToken) (AccessPayload, error) {
	token, err := jwt.ParseWithClaims(t, &accessClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.accessSecret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return AccessPayload{}, fmt.Errorf("%w: %w", ErrTokenExpired, err)
		}
		if errors.Is(err, jwt.ErrTokenInvalidClaims) {
			return AccessPayload{}, fmt.Errorf("%w: %w", ErrTokenInvalid, err)
		}
		return AccessPayload{}, err
	}

	claims, ok := token.Claims.(*accessClaims)
	if !ok {
		return AccessPayload{}, fmt.Errorf("%w: %w", ErrTokenInvalid, err)
	}
	userID, err := uuid.Parse(claims.Issuer)
	if err != nil {
		return AccessPayload{}, fmt.Errorf("%w: %w", ErrTokenInvalid, err)
	}
	return AccessPayload{
		userID: userID,
	}, nil
}

func (j JWT) ParseRefresh(t SignedToken) (RefreshPayload, error) {
	token, err := jwt.ParseWithClaims(t, &refreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.refreshSecret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return RefreshPayload{}, fmt.Errorf("%w: %w", ErrTokenExpired, err)
		}
		if errors.Is(err, jwt.ErrTokenInvalidClaims) {
			return RefreshPayload{}, fmt.Errorf("%w: %w", ErrTokenInvalid, err)
		}
		return RefreshPayload{}, err
	}
	claims, ok := token.Claims.(*refreshClaims)
	if !ok {
		return RefreshPayload{}, fmt.Errorf("%w: %w", ErrTokenInvalid, err)
	}
	_, err = uuid.Parse(claims.Issuer)
	if err != nil {
		return RefreshPayload{}, fmt.Errorf("%w: %w", ErrTokenInvalid, err)
	}
	return RefreshPayload{}, nil
}
