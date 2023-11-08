package handler

import (
	"advanceauth/backend/app/utils"
	"fmt"
	"gopkg.in/mail.v2"
)

type MailInfo struct {
	EmailTarget string
	Subject     string
	Body        string
}

var NewDeviceLogin = `
<h4>Hi,</h4>
<p>Someone just logged in to your account from a new device.</p>
<p>Device: %s</p>
<p>Location: %s</p>
<p>If this was you, you can ignore this email.</p>
<p>If this wasn't you, please change your password immediately.</p>
<p>Thanks,</p>
<p>Honexam</p>
`
var LoginTokenGotHacked = `
<h4>Hi,</h4>
<p>Someone suspicious just logged in to your account using your login token.</p>
<p>Device: %s</p>
<p>Location: %s</p>
<p>Your login token that stored at browser cookies got stoled.</p>
<p>We stop it from login and please secure your browser.</p>
<p>Thanks,</p>
<p>Honexam</p>
`

var ResetPassword = `
<h4>Hi, %s</h4>
<p>Someone just request to reset your password.</p>
<p>Device: %s</p>
<p>Location: %s</p>
<p>If it's you, click the button below to reset your password:</p>
<a href="%s" style="padding: 0.5rem; border: 1px solid">Reset Password</a>
<p>If it's not you, please ignore this email.</p>
<p>Thanks,</p>
<p>Honexam</p>
`

func SendMail(mailInfo MailInfo) {
	fmt.Println("Sending email...")
	go func() {
		message := mail.NewMessage()
		message.SetHeader("From", utils.GetEnv("MAILER_EMAIL"))
		message.SetHeader("To", mailInfo.EmailTarget)
		message.SetHeader("Subject", mailInfo.Subject)
		message.SetBody("text/html", mailInfo.Body)
		mailer := mail.NewDialer(
			"smtp.gmail.com",
			587,
			utils.GetEnv("MAILER_EMAIL"),
			utils.GetEnv("MAILER_PASSWORD"),
		)
		if err := mailer.DialAndSend(message); err != nil {
			fmt.Println(err)
		}
	}()
}
