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
	msgTemplateHello = &i18n.Message{ID: "app.templates.hello", Other: "Hello"}

	msgEmailVerifyTitle      = &i18n.Message{ID: "app.templates.email.verify.title", Other: "Verify Your Email"}
	msgEmailVerifySubject    = &i18n.Message{ID: "app.templates.email.verify.subject", Other: "Email Verification"}
	msgEmailVerifyBodyText   = &i18n.Message{ID: "app.templates.email.verify.body_text", Other: "Thank you for using our site, please verify your email by pressing the button bellow."}
	msgEmailVerifyButtonText = &i18n.Message{ID: "app.templates.email.verify.button_text", Other: "Verify Email"}

	msgPwdRecoveryTitle        = &i18n.Message{ID: "app.templates.email.verify.title", Other: "Reset Your Password"}
	msgPwdRecoverySubject      = &i18n.Message{ID: "app.templates.email.verify.subject", Other: "Password Recovery"}
	msgPwdRecoveryBodyText     = &i18n.Message{ID: "app.templates.email.verify.body_text", Other: "We got a request to reset your password, press the button bellow to reset it."}
	msgPwdRecoveryValidForText = &i18n.Message{ID: "app.templates.email.verify.valid_for_text", One: "This password reset is only valid for the next {{ .Expiry }} hour", Other: "This password reset is only valid for the next {{ .Expiry }} hours."}
	msgPwdRecoveryWarningText  = &i18n.Message{ID: "app.templates.email.verify.button_text", Other: "If you didn't request this, you can ignore this message and your password will remain unchanged."}
	msgPwdRecoveryButtonText   = &i18n.Message{ID: "app.templates.email.verify.warning_text", Other: "Reset Password"}

	msgPwdUpdatedSubject        = &i18n.Message{ID: "app.templates.password.updated.title", Other: "Password Update Completed"}
	msgPwdUpdatedTitle          = &i18n.Message{ID: "app.templates.password.updated.subject", Other: "Password Updated"}
	msgPwdUpdatedForAccountText = &i18n.Message{ID: "app.templates.password.updated.subject", Other: "Password for the account"}
	msgPwdUpdatedChangedText    = &i18n.Message{ID: "app.templates.password.updated.body_text", Other: "has been changed successfully!"}
	msgPwdUpdatedCompletedText  = &i18n.Message{ID: "app.templates.password.updated.button_text", Other: "Password Reset Completed"}
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
		"Hello":      locale.LocalizeDefaultMessage(l, msgTemplateHello),
		"Link":       fmt.Sprintf("%s/email/verify?token=%s", siteURL, token),
		"Title":      locale.LocalizeDefaultMessage(l, msgEmailVerifyTitle),
		"BodyText":   locale.LocalizeDefaultMessage(l, msgEmailVerifyBodyText),
		"ButtonText": locale.LocalizeDefaultMessage(l, msgEmailVerifyButtonText),
	}

	return a.sendEmailTemplate("templates/email_verify.html", data, info)
}

// SendPasswordRecoveryEmail sends the pwd reset email
func (a *App) SendPasswordRecoveryEmail(to string, username string, token *model.Token, siteURL string, userLocale string) *model.AppErr {
	l := locale.GetUserLocalizer(userLocale)

	info := &mailer.Maildata{
		To:      []string{to},
		Subject: locale.LocalizeDefaultMessage(l, msgPwdRecoverySubject),
	}

	displayName := username
	if username == "" {
		displayName = strings.Join(info.To, ",")
	}

	expiry := fmt.Sprintf("%.0f", token.ExpiresAt.Sub(token.CreatedAt).Hours())

	data := map[string]string{
		"Email":       strings.Join(info.To, ","),
		"Token":       token.Token,
		"DisplayName": displayName,
		"Hello":       locale.LocalizeDefaultMessage(l, msgTemplateHello),
		"Link":        fmt.Sprintf("%s/password/reset?token=%s", siteURL, token.Token),
		"Title":       locale.LocalizeDefaultMessage(l, msgPwdRecoveryTitle),
		"BodyText":    locale.LocalizeDefaultMessage(l, msgPwdRecoveryBodyText),
		"ButtonText":  locale.LocalizeDefaultMessage(l, msgPwdRecoveryButtonText),
		"ValidForText": locale.LocalizeWithConfig(l, &i18n.LocalizeConfig{
			DefaultMessage: msgPwdRecoveryValidForText,
			PluralCount:    expiry,
			TemplateData:   map[string]interface{}{"Expiry": expiry},
		}),
		"WarningText": locale.LocalizeDefaultMessage(l, msgPwdRecoveryWarningText),
	}

	return a.sendEmailTemplate("templates/reset_password.html", data, info)
}

// SendPasswordUpdatedEmail sends the pwd reset completed email
func (a *App) SendPasswordUpdatedEmail(to string, username string, siteURL string, userLocale string) *model.AppErr {
	l := locale.GetUserLocalizer(userLocale)

	info := &mailer.Maildata{
		To:      []string{to},
		Subject: locale.LocalizeDefaultMessage(l, msgPwdUpdatedSubject),
	}

	displayName := username
	if username == "" {
		displayName = strings.Join(info.To, ",")
	}

	data := map[string]string{
		"Email":          strings.Join(info.To, ","),
		"DisplayName":    displayName,
		"Hello":          locale.LocalizeDefaultMessage(l, msgTemplateHello),
		"Title":          locale.LocalizeDefaultMessage(l, msgPwdUpdatedTitle),
		"CompletedText":  locale.LocalizeDefaultMessage(l, msgPwdUpdatedCompletedText),
		"ChangedText":    locale.LocalizeDefaultMessage(l, msgPwdUpdatedChangedText),
		"ForAccountText": locale.LocalizeDefaultMessage(l, msgPwdUpdatedForAccountText),
	}

	return a.sendEmailTemplate("templates/reset_password_completed.html", data, info)
}
