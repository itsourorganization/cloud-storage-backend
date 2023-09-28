package signup

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/undefeel/cloud-storage-backend/internal/lib/jwt"
	"github.com/undefeel/cloud-storage-backend/internal/lib/server"
	"github.com/undefeel/cloud-storage-backend/internal/services"
	"github.com/undefeel/cloud-storage-backend/internal/transport"
	"log/slog"
	"net/http"
)

const (
	errIncorrectInputMsg = "Invalid login or password"
	errUniqueMsg         = "User with login already exist"
)

//go:generate go run github.com/vektra/mockery/v2@v2.33.3 --all

type UseCases interface {
	SignUp(ctx context.Context, user *services.User) (jwt.TokenPair, error)
}

// Sign Up godoc
// @Summary      Create user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param data body transport.User true "SignUp data"
// @Success      200  {object}  jwt.TokenPair
// @Failure      400  {object}  server.ErrMsg
// @Failure      500  {object}  server.ErrMsg
// @Router       /signUp [post]
func New(log *slog.Logger, service UseCases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var trUs transport.User
		err := json.NewDecoder(r.Body).Decode(&trUs)
		if err != nil {
			server.RespondError(w, r, http.StatusBadRequest, errIncorrectInputMsg)
			return
		}
		if !isValid(trUs) {
			server.RespondError(w, r, http.StatusBadRequest, errIncorrectInputMsg)
			return
		}

		pair, err := service.SignUp(r.Context(), transport.ConvertTransportUserToServicesUser(trUs))
		if err != nil {
			log.Error(err.Error())
			handleServiceError(err, w, r)
			return
		}
		server.RespondOK(pair, w, r)
	}
}

func isValid(trUs transport.User) bool {
	lenPas := len(trUs.Password)
	lenLog := len(trUs.Login)
	return lenPas > 8 && lenPas < 24 && lenLog > 6 && lenLog < 32
}

func handleServiceError(err error, w http.ResponseWriter, r *http.Request) {
	switch {
	case errors.Is(err, services.ErrIncorrectInput):
		server.RespondError(w, r, http.StatusBadRequest, errIncorrectInputMsg)
	case errors.Is(err, services.ErrUnique):
		server.RespondError(w, r, http.StatusBadRequest, errUniqueMsg)
	default:
		server.RespondError(w, r, http.StatusInternalServerError, err.Error())
	}
}
