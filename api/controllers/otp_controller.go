package controllers

import (
	"api-center/api/models"
	"api-center/api/responses"
	"encoding/json"
	"errors"
	"net/http"
)

func (s *Server) SendOtp(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	if request.Email == "" {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("email is required"))
		return
	}

	otp := models.OtpStore{}
	otp.Email = request.Email
	err := otp.SendOtp(s.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "OTP sent successfully",
	})
}

func (s Server) ValidateOtp(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Otp   string `json:"otp"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	if request.Email == "" {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("email is required"))
		return
	}
	if request.Otp == "" {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("otp is required"))
		return
	}
	otp := models.OtpStore{}
	otp.Otp = request.Otp
	otp.Email = request.Email
	err := otp.ValidateOtp(s.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	responses.JSON(w, responses.JSONResponse{
		Status:  http.StatusOK,
		Message: "OTP verified successfully",
	})
}
