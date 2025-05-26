package pkg

import (
	"fmt"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	"gopkg.in/gomail.v2"
)

func SendMail(to []string, otp string) error {
	subject := "OTP for Onboarding ACM's Season of Code 2025"
	body := fmt.Sprintf("Your OTP for logging into the Season of Code is %s. This is valid for only 5 minutes.", otp)

	m := gomail.NewMessage()
	m.SetHeader("From", cmd.AppConfig.SMTP.Username)
	m.SetHeader("To", to...)
	m.SetHeader("subject", subject)
	m.SetBody("text/plain", body)

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
