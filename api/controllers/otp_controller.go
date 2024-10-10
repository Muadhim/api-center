package controllers

import (
	"api-center/api/models"
	"api-center/api/responses"
	"encoding/json"
	"net/http"
)

func (s *Server) SendOtp(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if request.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
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
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}
	if request.Otp == "" {
		http.Error(w, "Otp is required", http.StatusBadRequest)
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
