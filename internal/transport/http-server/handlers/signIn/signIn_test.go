package signIn

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
	"github.com/undefeel/cloud-storage-backend/internal/transport/http-server/handlers/signIn/mocks"
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
				Login:    "testLogin111",
				Password: "testPassword11",
			},
			respCode: http.StatusOK,
		},
		{
			name: "Not Found",
			inp: transport.User{
				Login:    "testLogin111",
				Password: "testPassword11",
			},
			respCode: http.StatusNotFound,
			respErr:  errNotFoundMsg,
			mockErr:  services.ErrNotFound,
		},
		{
			name: "Mismatched passwords",
			inp: transport.User{
				Login:    "testLogin111",
				Password: "testPassword11",
			},
			respCode: http.StatusBadRequest,
			respErr:  errIncorrectInputMsg,
			mockErr:  services.ErrIncorrectInput,
		},
		{
			name: "Empty",
			inp: transport.User{
				Login:    "",
				Password: "",
			},
			respCode: http.StatusBadRequest,
			respErr:  errIncorrectInputMsg,
		},
	}

	for _, cs := range cases {
		t.Run(cs.name, func(t *testing.T) {
			service := mocks.NewUseCases(t)
			if !(cs.respCode != http.StatusOK && cs.mockErr == nil) {
				service.On("SignIn",
					mock.Anything,
					mock.AnythingOfType("*services.User")).
					Return(jwt.TokenPair{}, cs.mockErr).Once()
			}
			handler := New(slogDiscard.NewDiscardLogger(), service)
			input, err := json.Marshal(cs.inp)
			require.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/signIn", bytes.NewReader(input))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			body := rr.Body.String()
			require.Equal(t, rr.Code, cs.respCode)
			if cs.respCode != http.StatusOK {
				var s server.ErrMsg
				require.NoError(t, json.NewDecoder(bytes.NewReader([]byte(body))).Decode(&s))
			} else {
				var resp jwt.TokenPair
				require.NoError(t, json.NewDecoder(bytes.NewReader([]byte(body))).Decode(&resp))
			}
		})
	}
}
