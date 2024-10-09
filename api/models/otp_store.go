package models

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"gorm.io/gorm"
)

type OtpStore struct {
	ID        uint      `gorm:"primary_key;autoIncrement" json:"id"`
	Email     string    `gorm:"type:varchar(255);default null" json:"email"`
	Otp       string    `gorm:"type:varchar(6);default null" json:"otp"`
	Expiry    time.Time `gorm:"not null" json:"expiry"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (o *OtpStore) SendOtp(db *gorm.DB) error {
	o.Otp = fmt.Sprintf("%06d", rand.New(rand.NewSource(time.Now().UnixNano())).Intn(1000000))
	o.Expiry = time.Now().Add(5 * time.Minute)

	if err := db.Debug().Create(&o).Error; err != nil {
		return err
	}

	message := `
	<div style="border:1px solid #fff; width: 500px; padding:10px; border-radius: 10px; background-color: hsla(221.2, 83.2%, 53.3%, 0.7); color: #fff; text-align:center; font-family: Arial, Helvetica, sans-serif;">
	<h1>Forgot Password OPT</h1>
	<p>Use this OPT to change your password this OTP will expired in 5 minutes</p>
	<div style="background-color: #ffffffaa; color:#000; width:200px; padding: 2px 10px; border-radius: 10px; margin: 0 auto;">
	<h2 style="">` + strings.Trim(strings.ReplaceAll(o.Otp, string(o.Otp[3]), " "+string(o.Otp[3])), " ") + `</h2>
	</div>
	</div>
	`

	// Send email (use your SMTP server)
	email := Email{}
	email.Email = o.Email
	email.Message = message
	email.Subject = "Forgot Password OTP <Sweet API>"

	go func() {
		err := email.SendEmail()
		if err != nil {
			log.Println(err)
		}
	}()

	return nil
}

func (o *OtpStore) ValidateOtp(db *gorm.DB) error {
	var optEntry OtpStore

	if err := db.Debug().Where("otp = ? AND email = ?", o.Otp, o.Email).First(&optEntry).Error; err != nil {
		return errors.New("invalid otp")
	}
	if time.Now().After(optEntry.Expiry) {
		return errors.New("otp expired")
	}

	return nil
}

func (o *OtpStore) DeleteOtp(db *gorm.DB) error {
	return db.Where("email = ? OR otp = ? OR expiry < ?", o.Email, o.Otp, time.Now()).Delete(&o).Error
}
