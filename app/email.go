package app

import (
	"strings"

	"github.com/dankobgd/ecommerce-shop/mailer"
	"github.com/dankobgd/ecommerce-shop/model"
)

func (a *App) sendMail(templateName string, templateData interface{}, info mailer.MailData) *model.AppErr {
	return mailer.Send(&a.Cfg().EmailSettings, info, templateName, templateData)
}

// SendWelcomeEmail sends the email o the newly registered user
func (a *App) SendWelcomeEmail(to string) *model.AppErr {
	info := mailer.MailData{
		Subject: "some test title",
		To:      []string{to},
	}
	data := map[string]string{"Email": strings.Join(info.To, ",")}
	return a.sendMail("templates/welcome.html", data, info)
}
