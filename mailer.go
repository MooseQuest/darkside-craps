package main

import (
	"fmt"
	"net/smtp"
	"strings"
)

// fromAddress extracts the bare "a@b.com" envelope address from a MailFrom that
// may be in "Name <a@b.com>" form.
func fromAddress(from string) string {
	if i := strings.LastIndex(from, "<"); i >= 0 {
		return strings.TrimSuffix(strings.TrimSpace(from[i+1:]), ">")
	}
	return strings.TrimSpace(from)
}

// sendVerificationEmail delivers a one-time code via SMTP. Uses STARTTLS when
// the relay offers it (net/smtp.SendMail negotiates it automatically) and PLAIN
// auth when credentials are configured. Works with any transactional relay
// (Resend/SendGrid/SES/Mailgun/…) on a submission port (587).
func (a *App) sendVerificationEmail(to, code string) error {
	c := a.cfg
	subject := "Your Dark Side Craps sign-in code"
	body := fmt.Sprintf(
		"Your verification code is %s\r\n\r\n"+
			"Enter it to finish setting up your passkey. It expires in 10 minutes.\r\n"+
			"If you didn't request this, you can ignore this email.\r\n",
		code)
	msg := "" +
		"From: " + c.MailFrom + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" + body

	var auth smtp.Auth
	if c.SMTPUser != "" {
		auth = smtp.PlainAuth("", c.SMTPUser, c.SMTPPass, c.SMTPHost)
	}
	return smtp.SendMail(c.SMTPHost+":"+c.SMTPPort, auth, fromAddress(c.MailFrom), []string{to}, []byte(msg))
}
