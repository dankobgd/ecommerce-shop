package mailer

import (
	"fmt"
	"net/http"
	"net/smtp"
	"net/textproto"
	"strings"

	"github.com/dankobgd/ecommerce-shop/config"
	"github.com/dankobgd/ecommerce-shop/model"
	"github.com/dankobgd/ecommerce-shop/utils/locale"
	"github.com/dankobgd/ecommerce-shop/utils/template"
	"github.com/dankobgd/ecommerce-shop/zlog"
	"github.com/jordan-wright/email"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"jaytaylor.com/html2text"
)

var (
	msgSendMailSMTP       = &i18n.Message{ID: "mailer.send_smpt.app_error", Other: "could not send email with smtp"}
	msgSendMailSendgrid   = &i18n.Message{ID: "mailer.send_sendgrid.app_error", Other: "could not send email with sendgird"}
	msgParseEmailTemplate = &i18n.Message{ID: "mailer.parse_template.app_error", Other: "could not parse email template"}
	msgParseHTMLToText    = &i18n.Message{ID: "mailer.parse_html2text.app_error", Other: "could not parse email html to text"}
)

// Maildata holds email details
type Maildata struct {
	To       []string
	Subject  string
	textBody string
	hTMLBody string
	Headers  textproto.MIMEHeader
}

func sendWithSMTP(settings config.EmailSettings, md *Maildata) *model.AppErr {
	e := &email.Email{
		From:    settings.FeedbackEmail,
		To:      md.To,
		Subject: md.Subject,
		HTML:    []byte(md.hTMLBody),
		Text:    []byte(md.textBody),
		Headers: md.Headers,
	}

	var addr string
	var auth smtp.Auth

	if settings.Enabled {
		addr = fmt.Sprintf("%s:%d", settings.SMTPHost, settings.SMTPPort)
		auth = smtp.PlainAuth("", settings.SMTPUsername, settings.SMTPPassword, settings.SMTPHost)
	} else {
		addr = fmt.Sprintf("%s:%d", settings.MailTrap.Host, settings.MailTrap.Port)
		auth = smtp.PlainAuth("", settings.MailTrap.Username, settings.MailTrap.Password, settings.MailTrap.Host)
	}

	if err := e.Send(addr, auth); err != nil {
		zlog.Info("could not send email with smtp", zlog.String("recipients:", strings.Join(e.To, ",")), zlog.Err(err))
		return model.NewAppErr("mailer.Send", model.ErrInternal, locale.GetUserLocalizer("en"), msgSendMailSMTP, http.StatusInternalServerError, nil)
	}
	zlog.Info("email has been sent successfully", zlog.String("recipients:", strings.Join(e.To, ",")))
	return nil
}

func sendWithSendgrid(config *config.Config, md *Maildata) *model.AppErr {
	mailSettings := config.EmailSettings

	from := mail.NewEmail(mailSettings.FeedbackUser, mailSettings.FeedbackEmail)
	to := mail.NewEmail("", strings.Join(md.To, ","))
	subject := md.Subject
	textContent := md.textBody
	htmlContent := md.hTMLBody

	message := mail.NewSingleEmail(from, subject, to, textContent, htmlContent)
	client := sendgrid.NewSendClient(mailSettings.Sendgrid.APIKey)

	resp, err := client.Send(message)
	if err != nil {
		zlog.Info("could not send email with sendgird", zlog.String("recipients:", strings.Join(md.To, ",")), zlog.Err(err))
		return model.NewAppErr("mailer.Send", model.ErrInternal, locale.GetUserLocalizer("en"), msgSendMailSendgrid, http.StatusInternalServerError, nil)
	}

	zlog.Info("email successfully sent with sendgrid", zlog.String("recipients:", strings.Join(md.To, ",")), zlog.Int("statusCode", resp.StatusCode))
	return nil
}

// SendEmailTemplate parses the email template to get the html and text contents
// and sends the email with the information
func SendEmailTemplate(filename string, data interface{}, md *Maildata, config *config.Config) *model.AppErr {
	htmlString, err := template.ParseTemplate(filename, data)
	if err != nil {
		zlog.Info("could not parse email template", zlog.Err(err))
		return model.NewAppErr("mailer.parse_template", model.ErrInternal, locale.GetUserLocalizer("en"), msgParseEmailTemplate, http.StatusInternalServerError, nil)
	}

	textString, err := html2text.FromString(htmlString)
	if err != nil {
		zlog.Info("could not parse html to string", zlog.Err(err))
		return model.NewAppErr("mailer.parse_html2text", model.ErrInternal, locale.GetUserLocalizer("en"), msgParseHTMLToText, http.StatusInternalServerError, nil)
	}

	md.hTMLBody = htmlString
	md.textBody = textString

	return send(config, md)
}

func send(config *config.Config, md *Maildata) *model.AppErr {
	settings := config.EmailSettings
	switch settings.Transport {
	case "smtp":
		return sendWithSMTP(settings, md)
	case "sendgrid":
		return sendWithSendgrid(config, md)
	default:
		panic(fmt.Errorf("could not configure mailer: unknown email transport"))
	}
}
