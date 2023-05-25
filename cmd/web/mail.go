package main

import (
	"fmt"
	"github.com/porky256/course-project/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
	"os"
	"strings"
	"time"
)

func listenForEmail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}
	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.From).SetSubject(m.Subject)
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content)
	} else {
		data, err := os.ReadFile(fmt.Sprintf("./static/email-templates/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println(err)
			return
		}

		mailTemplate := string(data)
		msgToSend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgToSend)
	}
	err = email.Send(client)
	if err != nil {
		app.ErrorLog.Println(err)
	} else {
		app.InfoLog.Println("Email sent!")
	}
}
