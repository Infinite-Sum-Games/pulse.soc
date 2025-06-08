package pkg

import (
	"bytes"
	"html/template"
	"path/filepath"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	"gopkg.in/gomail.v2"
)

type OTPEmailData struct {
	FirstDigit  string
	SecondDigit string
	ThirdDigit  string
	FourthDigit string
	FifthDigit  string
	SixthDigit  string
}

func LoadAndRenderTemplate(data any) (string, error) {
	templatePath := filepath.Join("pkg", "mail.htm")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var rendered bytes.Buffer
	err = tmpl.Execute(&rendered, data)
	if err != nil {
		return "", err
	}
	return rendered.String(), nil
}

func SendMail(to []string, otp string) error {
	subject := "Amrita Summer of Code 2025 Welcomes You!"

	data := OTPEmailData{
		FirstDigit:  string(otp[0]),
		SecondDigit: string(otp[1]),
		ThirdDigit:  string(otp[2]),
		FourthDigit: string(otp[3]),
		FifthDigit:  string(otp[4]),
		SixthDigit:  string(otp[5]),
	}

	body, err := LoadAndRenderTemplate(data)
	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", cmd.AppConfig.SMTP.Username)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(
		cmd.AppConfig.SMTP.Host,
		cmd.AppConfig.SMTP.Port,
		cmd.AppConfig.SMTP.Username,
		cmd.AppConfig.SMTP.AppPassword,
	)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	cmd.Log.Info("[SUCCESS]: Email sent successfully.")
	return nil
}
