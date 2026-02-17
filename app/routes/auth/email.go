package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"html/template"
	"mms/app/config"
	"net/smtp"
)

func GenerateResetToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func SendPasswordResetEmail(email, resetToken, host string) error {
	smtpConfig := config.AppConfig.SMTP

	// Skip sending email if SMTP is not configured
	if smtpConfig.Host == "" || smtpConfig.Username == "" {
		fmt.Printf("SMTP not configured. Reset link for %s: http://%s/auth/reset-password?token=%s\n", email, host, resetToken)
		return nil
	}

	auth := smtp.PlainAuth("", smtpConfig.Username, smtpConfig.Password, smtpConfig.Host)

	resetLink := fmt.Sprintf("http://%s/auth/reset-password?token=%s", host, resetToken)

	// Parse email template
	tmpl, err := template.ParseFiles("app/templates/emails/password-reset.html")
	if err != nil {
		return err
	}

	// Execute template with data
	var body bytes.Buffer
	err = tmpl.Execute(&body, map[string]string{
		"ResetLink": resetLink,
	})
	if err != nil {
		return err
	}

	subject := "Password Reset - MMS Motors"
	msg := fmt.Sprintf("To: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", email, subject, body.String())

	return smtp.SendMail(fmt.Sprintf("%s:%d", smtpConfig.Host, smtpConfig.Port), auth, smtpConfig.From, []string{email}, []byte(msg))
}
