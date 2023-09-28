package server

import (
	"encoding/json"
	"net/http"
)

type ErrMsg struct {
	Msg string `json:"msg"`
}

func RespondError(w http.ResponseWriter, r *http.Request, code int, msg string) {
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	var err ErrMsg
	err.Msg = msg
	data, _ := json.Marshal(err)
	w.WriteHeader(code)
	w.Write(data)
}
