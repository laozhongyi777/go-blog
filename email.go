package main

import (
	"fmt"
	"net/smtp"
	"strings"
)

func sendEmail(code string, emaddr string) {
	auth := smtp.PlainAuth("", "qjkk88369yy@163.com", "Jq325com", "smtp.163.com")
	to := []string{emaddr}
	from := "qjkk88369yy@163.com"

	subject := "test mail"
	content_type := "Content-Type: text/plain; charset=UTF-8"
	body := "Your code is:" + code
	msg := []byte("To: " + strings.Join(to, ",") + "\r\nFrom: " + "test" +
		"<" + from + ">\r\nSubject: " + subject + "\r\n" + content_type + "\r\n\r\n" + body)

	err := smtp.SendMail("smtp.163.com:25", auth, from, to, msg)
	if err != nil {
		fmt.Println("error:", err)
	}
}
