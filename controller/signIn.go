package controller

import (
	"bytes"
	"encoding/json"
	"github.com/Ararat25/go_final_project/customError"
	"net/http"
	"os"
)

type ResponseWithToken struct {
	Token string `json:"token"`
}

type RequestWithPassword struct {
	Password string `json:"password"`
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var req RequestWithPassword
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = json.Unmarshal(buf.Bytes(), &req)
	if err != nil {
		sendErrorResponseData(w, http.StatusBadRequest, err.Error())
		return
	}

	password := req.Password

	if len(password) == 0 {
		sendErrorResponseData(w, http.StatusBadRequest, customError.ErrPasswordNotSpecified.Error())
		return
	}

	envPassword := os.Getenv("TODO_PASSWORD")
	if len(envPassword) == 0 {
		sendErrorResponseData(w, http.StatusInternalServerError, customError.ErrEnvPasswordNotSpecified.Error())
		return
	}

	token, err := h.service.VerifyUser(envPassword, password)
	if err != nil {
		sendErrorResponseData(w, http.StatusConflict, err.Error())
		return
	}

	respBytes, err := json.Marshal(ResponseWithToken{Token: token})
	if err != nil {
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(respBytes)
}
