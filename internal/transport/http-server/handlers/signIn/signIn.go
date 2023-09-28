package signIn

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
	errNotFoundMsg       = "User with this login not found"
)

//go:generate go run github.com/vektra/mockery/v2@v2.33.3 --all

type UseCases interface {
	SignIn(ctx context.Context, user *services.User) (jwt.TokenPair, error)
}

// SignIn godoc
// @Summary      Create user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param data body transport.User true "SignIn data"
// @Success      200  {object}  jwt.TokenPair
// @Failure      400  {object}  server.ErrMsg
// @Failure      404  {object}  server.ErrMsg
// @Failure      500  {object}  server.ErrMsg
// @Router       /signIn [post]
func New(log *slog.Logger, service UseCases) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var trUs transport.User
		err := json.NewDecoder(r.Body).Decode(&trUs)
		if err != nil {
			log.Error(err.Error())
			server.RespondError(w, r, http.StatusBadRequest, errIncorrectInputMsg)
		}
		if !isValid(trUs) {
			server.RespondError(w, r, http.StatusBadRequest, errIncorrectInputMsg)
			return
		}

		pair, err := service.SignIn(r.Context(), transport.ConvertTransportUserToServicesUser(trUs))
		if err != nil {
			log.Error(err.Error())
			switch {
			case errors.Is(err, services.ErrIncorrectInput):
				server.RespondError(w, r, http.StatusBadRequest, errIncorrectInputMsg)
				return
			case errors.Is(err, services.ErrNotFound):
				server.RespondError(w, r, http.StatusNotFound, errNotFoundMsg)
				return
			default:
				server.RespondError(w, r, http.StatusInternalServerError, err.Error())
				return
			}
		}
		server.RespondOK(pair, w, r)
	}
}

func isValid(trUs transport.User) bool {
	lenPas := len(trUs.Password)
	lenLog := len(trUs.Login)
	return lenPas > 8 && lenPas < 24 && lenLog > 6 && lenLog < 32
}
