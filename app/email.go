package app

import (
	"strings"

	"github.com/dankobgd/ecommerce-shop/mailer"
	"github.com/dankobgd/ecommerce-shop/model"
)

func (a *App) sendEmailTemplate(filename string, data interface{}, maildata *mailer.Maildata) *model.AppErr {
	return mailer.SendEmailTemplate(filename, data, maildata, a.Cfg())
}

// SendWelcomeEmail sends the email o the newly registered user
func (a *App) SendWelcomeEmail(to string) *model.AppErr {
	info := &mailer.Maildata{
		To:      []string{to},
		Subject: "Test Subject",
	}
	data := map[string]string{"Email": strings.Join(info.To, ",")}
	return a.sendEmailTemplate("templates/welcome.html", data, info)
}
