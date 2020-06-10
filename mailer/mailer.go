package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"net/smtp"
	"strings"

	"github.com/dankobgd/ecommerce-shop/config"
	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dankobgd/ecommerce-shop/zlog"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var (
	msgParseTpl  = &i18n.Message{ID: "mailer.send.app_error", Other: "could not parse email template"}
	msgSendEmail = &i18n.Message{ID: "mailer.send.app_error", Other: "could not send email"}
)

const (
	mimeHTML = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	mimeTEXT = "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
)

// MailData has the email information
type MailData struct {
	Subject string
	To      []string
}

func parseTemplate(fileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return "", err
	}

	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func sendMailWithConfig(config *config.EmailSettings, data MailData, htmlString string) error {
	body := "To: " + data.To[0] + "\r\nSubject: " + data.Subject + "\r\n" + mimeHTML + "\r\n" + htmlString

	if config.Enabled {
		SMTP := fmt.Sprintf("%s:%d", config.SMTPHost, config.SMTPPort)
		if err := smtp.SendMail(SMTP, smtp.PlainAuth("", config.SMTPUsername, config.SMTPPassword, config.SMTPHost), config.SMTPFrom, data.To, []byte(body)); err != nil {
			return err
		}
	} else {
		SMTP := fmt.Sprintf("%s:%d", config.MailTrap.Host, config.MailTrap.Port)
		if err := smtp.SendMail(SMTP, smtp.PlainAuth("", config.MailTrap.Username, config.MailTrap.Password, config.MailTrap.Host), config.SMTPFrom, data.To, []byte(body)); err != nil {
			return err
		}
	}

	return nil
}

// Send sends an email with the given template name
func Send(config *config.EmailSettings, data MailData, templateName string, templateData interface{}) *model.AppErr {
	body, err := parseTemplate(templateName, templateData)
	if err != nil {
		zlog.Error("could not parse email template", zlog.Err(err))
		return model.NewAppErr("mailer.Send", model.ErrInternal, locale.GetUserLocalizer("en"), msgParseTpl, http.StatusInternalServerError, nil)
	}
	if err := sendMailWithConfig(config, data, body); err != nil {
		zlog.Error("could not send the email", zlog.String("recepients:", strings.Join(data.To, ",")), zlog.Err(err))
		return model.NewAppErr("mailer.Send", model.ErrInternal, locale.GetUserLocalizer("en"), msgSendEmail, http.StatusInternalServerError, nil)
	}
	zlog.Info("email has been sent successfuly", zlog.String("recepients:", strings.Join(data.To, ",")))
	return nil
}
