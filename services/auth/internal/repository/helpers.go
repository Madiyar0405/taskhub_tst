package repository

import (
	"auth/internal/config"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"
)

func GenerateRandomCode() (string, error) {
	const op = "repository.helpers.GenerateRandomCode"

	const codeLength = 6
	var code string
	for i := 0; i < codeLength; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("%s:%w", op, err)
		}
		code += fmt.Sprintf("%d", num.Int64())
	}
	return code, nil
}

func SendEmailWithCode(to, code string) error {
	const op = "repository.helpers.SendEmailWithCode"

	cfg := config.MustLoad()

	from := cfg.EmailSenderAddress
	password := cfg.EmailSenderPassword

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	message := "Subject: Your confirmation code\r\n\r\nConfrimation code: " + code

	auth := smtp.PlainAuth("", from, password, smtpHost)

	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("%s:%w", op, err)
	}

	return nil
}
