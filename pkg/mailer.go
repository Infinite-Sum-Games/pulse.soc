package pkg

import (
	"fmt"

	"github.com/IAmRiteshKoushik/pulse/cmd"
	"gopkg.in/gomail.v2"
)

func SendMail(to []string, subject, body string) error {

	m := gomail.NewMessage()
	m.SetHeader("From", cmd.EnvVars.GmailUser)
	m.SetHeader("To", to...)
	m.SetHeader("subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(
		cmd.EnvVars.SmtpHost,
		cmd.EnvVars.SmtpPort,
		cmd.EnvVars.GmailUser,
		cmd.EnvVars.AppPassword,
	)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("Could not send email: %v", err)
	}

	cmd.Log.Info("[SUCCESS]: Email send successfully.")
	return nil
}
