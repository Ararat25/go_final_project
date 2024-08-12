package controller

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/Ararat25/go_final_project/errors"
)

// ResponseWithToken структура для ответа с токеном
type ResponseWithToken struct {
	Token string `json:"token"`
}

// RequestWithPassword структура для запроса с паролем
type RequestWithPassword struct {
	Password string `json:"password"`
}

// SignIn обработчик для авторизации
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
		sendErrorResponseData(w, http.StatusBadRequest, errors.ErrPasswordNotSpecified.Error())
		return
	}

	envPassword := h.service.Config.Password
	if len(envPassword) == 0 {
		sendErrorResponseData(w, http.StatusInternalServerError, errors.ErrEnvPasswordNotSpecified.Error())
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
	_, err = w.Write(respBytes)
	if err != nil {
		log.Println(err)
		sendErrorResponseData(w, http.StatusInternalServerError, err.Error())
		return
	}
}
