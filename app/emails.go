package app

import (
	"fmt"
	"strings"

	"github.com/dankobgd/ecommerce-shop/mailer"
	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgEmailVerifyTitle      = &i18n.Message{ID: "app.templates.email.verify.title", Other: "Verify Your Email"}
	msgEmailVerifySubject    = &i18n.Message{ID: "app.templates.email.verify.subject", Other: "Email Verification"}
	msgEmailVerifyHello      = &i18n.Message{ID: "app.templates.email.verify.hello", Other: "Hello"}
	msgEmailVerifyBodyText   = &i18n.Message{ID: "app.templates.email.verify.body_text", Other: "Thank you for using our site, please verify your email by pressing the button bellow."}
	msgEmailVerifyButtonText = &i18n.Message{ID: "app.templates.email.verify.button_text", Other: "Verify Email"}
)

func (a *App) sendEmailTemplate(filename string, data interface{}, maildata *mailer.Maildata) *model.AppErr {
	return mailer.SendEmailTemplate(filename, data, maildata, a.Cfg())
}

// SendWelcomeEmail sends the email o the newly registered user
func (a *App) SendWelcomeEmail(to string) *model.AppErr {
	info := &mailer.Maildata{
		To:      []string{to},
		Subject: "Welcome",
	}
	data := map[string]string{"Email": strings.Join(info.To, ",")}
	return a.sendEmailTemplate("templates/welcome.html", data, info)
}

// SendEmailVerificationEmail sends the verify email
func (a *App) SendEmailVerificationEmail(to string, token string, siteURL string, userLocale string) *model.AppErr {
	l := locale.GetUserLocalizer(userLocale)

	info := &mailer.Maildata{
		To:      []string{to},
		Subject: locale.LocalizeDefaultMessage(l, msgEmailVerifySubject),
	}

	data := map[string]string{
		"Email":      strings.Join(info.To, ","),
		"Token":      token,
		"Link":       fmt.Sprintf("%s/email/verify?token=%s", siteURL, token),
		"Title":      locale.LocalizeDefaultMessage(l, msgEmailVerifyTitle),
		"BodyText":   locale.LocalizeDefaultMessage(l, msgEmailVerifyBodyText),
		"ButtonText": locale.LocalizeDefaultMessage(l, msgEmailVerifyButtonText),
		"Hello":      locale.LocalizeDefaultMessage(l, msgEmailVerifyHello),
	}

	return a.sendEmailTemplate("templates/email_verify.html", data, info)
}
