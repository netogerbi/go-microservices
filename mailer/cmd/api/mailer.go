package main

import (
	"bytes"
	"html/template"
	"os"
	"strconv"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromAddress string
	FromName    string
}

type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Attachments []string
	Data        any
	DataMap     map[string]any
}

func NewMailer() Mail {
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))

	return Mail{
		Domain:      os.Getenv("MAIL_DOMAIN"),
		Host:        os.Getenv("MAIL_HOST"),
		Port:        port,
		Username:    os.Getenv("MAIL_USERNAME"),
		Password:    os.Getenv("MAIL_PASSWORD"),
		Encryption:  os.Getenv("MAIL_ENCRYPTION"),
		FromAddress: os.Getenv("MAIL_FROM_ADDRESS"),
		FromName:    os.Getenv("MAIL_FROM_NAME"),
	}
}

func (m *Mail) SendSMTPMessage(msg Message) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	msg.DataMap = map[string]any{
		"message": msg.Data,
	}

	formattedMsg, err := m.buildHTML(msg)
	if err != nil {
		return err
	}

	plainText, err := m.buildPlainText(msg)

	if err != nil {
		return err
	}

	email := mail.NewMSG()

	email.
		SetFrom(msg.From).
		AddTo(msg.To).
		SetSubject(msg.Subject).
		SetBody(mail.TextPlain, plainText).
		AddAlternative(mail.TextHTML, formattedMsg)

	for _, v := range msg.Attachments {
		email.AddAttachment(v)
	}

	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second
	server.KeepAlive = false
	server.Encryption = m.getEncryption(m.Encryption)
	smtpClient, err := server.Connect()

	if err != nil {
		return err
	}

	if err = email.Send(smtpClient); err != nil {
		return err
	}

	return nil
}

func (m *Mail) buildHTML(msg Message) (string, error) {
	templateFilePath := "./templates/mail.html.gohtml"

	t, err := template.New("email-html").ParseFiles(templateFilePath)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	formattedMsg, err := m.inlineCSS(tpl.String())
	if err != nil {
		return "", err
	}

	return formattedMsg, nil
}

func (m *Mail) buildPlainText(msg Message) (string, error) {
	templateFilePath := "./templates/mail.plain.gohtml"

	t, err := template.New("email-plain").ParseFiles(templateFilePath)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	plainTxt := tpl.String()

	return plainTxt, nil
}

func (m *Mail) inlineCSS(s string) (string, error) {
	options := premailer.Options{
		RemoveClasses:     false,
		CssToAttributes:   false,
		KeepBangImportant: true,
	}

	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSL
	case "none", "":
		return mail.EncryptionNone
	default:
		return mail.EncryptionSTARTTLS
	}
}
