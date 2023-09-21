package signup_test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/undefeel/cloud-storage-backend/internal/lib/jwt"
	"github.com/undefeel/cloud-storage-backend/internal/lib/logger/slogDiscard"
	"github.com/undefeel/cloud-storage-backend/internal/lib/server"
	"github.com/undefeel/cloud-storage-backend/internal/services"
	"github.com/undefeel/cloud-storage-backend/internal/transport"
	signup "github.com/undefeel/cloud-storage-backend/internal/transport/http-server/handlers/signUp"
	"github.com/undefeel/cloud-storage-backend/internal/transport/http-server/handlers/signUp/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	cases := []struct {
		name     string
		inp      transport.User
		respErr  string
		respCode int
		mockErr  error
	}{
		{
			name: "Success",
			inp: transport.User{
				Login:    "testLogin",
				Password: "testPassword",
			},
			respCode: http.StatusOK,
		},
		{
			name: "Invalid",
			inp: transport.User{
				Login:    "1",
				Password: "2",
			},
			respCode: http.StatusBadRequest,
		},
		{
			name: "Conflict",
			inp: transport.User{
				Login:    "ConflictLogin",
				Password: "testPassword",
			},
			respCode: http.StatusBadRequest,
			mockErr:  services.ErrUnique,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			service := mocks.NewUseCases(t)
			if !(c.respCode != http.StatusOK && c.mockErr == nil) {
				service.On("SignUp",
					mock.Anything,
					mock.AnythingOfType("*services.User")).
					Return(jwt.TokenPair{}, c.mockErr).Once()
			}
			handler := signup.New(slogDiscard.NewDiscardLogger(), service)

			input, err := json.Marshal(c.inp)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/signUp", bytes.NewReader(input))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			body := rr.Body.String()
			require.Equal(t, rr.Code, c.respCode)
			if c.respCode != http.StatusOK {
				var s server.ErrMsg
				require.NoError(t, json.NewDecoder(bytes.NewReader([]byte(body))).Decode(&s))
			} else {
				var resp jwt.TokenPair
				require.NoError(t, json.NewDecoder(bytes.NewReader([]byte(body))).Decode(&resp))
			}
		})
	}
}
