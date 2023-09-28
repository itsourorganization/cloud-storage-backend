package jwt_test

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/undefeel/cloud-storage-backend/internal/lib/jwt"
	"testing"
	"time"
)

func TestJWT_NewPair(t *testing.T) {
	j := jwt.New("accessSecret", time.Second*5, "refreshSecret", time.Second*10)
	pair, err := j.NewPair(uuid.New())
	require.NoError(t, err)
	if pair.Access.Token == "" && pair.Refresh.Token == "" {
		t.Fatal("unexpected")
	}

}

func TestJWT_ParseAccess(t *testing.T) {
	cases := []struct {
		name          string
		j             jwt.JWT
		sleepTime     time.Duration
		expectedError error
	}{
		{
			name:          "default",
			j:             jwt.New("accessSecret", time.Second*2, "refreshSecret", time.Second*3),
			sleepTime:     0,
			expectedError: nil,
		},
		{
			name:          "expired",
			j:             jwt.New("accessSecret", time.Second, "refreshSecret", time.Second*3),
			sleepTime:     time.Second * 2,
			expectedError: jwt.ErrTokenExpired,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			t.Parallel()
			pair, err := c.j.NewPair(uuid.New())
			require.NoError(t, err)

			time.Sleep(c.sleepTime)
			_, err = c.j.ParseAccess(pair.Access.Token)
			if err != nil {
				if !errors.Is(err, c.expectedError) {
					t.Error("bad err", err)
				}
			}
		})
	}
}

func TestJWT_ParseRefresh(t *testing.T) {
	cases := []struct {
		name          string
		j             jwt.JWT
		sleepTime     time.Duration
		expectedError error
	}{
		{
			name:          "default",
			j:             jwt.New("accessSecret", time.Second*2, "refreshSecret", time.Second*3),
			sleepTime:     0,
			expectedError: nil,
		},
		{
			name:          "expired",
			j:             jwt.New("accessSecret", time.Second, "refreshSecret", time.Second*3),
			sleepTime:     time.Second * 2,
			expectedError: jwt.ErrTokenExpired,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			pair, err := c.j.NewPair(uuid.New())
			require.NoError(t, err)

			time.Sleep(c.sleepTime)
			_, err = c.j.ParseRefresh(pair.Refresh.Token)
			if err != nil {
				if !errors.Is(err, c.expectedError) {
					t.Error("untyped error", err)
				}
			}
		})
	}
}
