package email

import (
	"QqVideo/config"
	"gopkg.in/gomail.v2"
	"log"
	"mime"
)

func SendEmail(toEmail, subject, msg string) {
	m := gomail.NewMessage()
	m.SetHeader("From", mime.QEncoding.Encode("UTF-8", "腾讯视频Vip积分签到通知")+"<"+config.EmailUsername+">")
	m.SetHeader("To", toEmail, toEmail)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", msg)

	d := gomail.NewDialer(config.EmailHost, config.EmailPort, config.EmailUsername, config.EmailPass)

	if err := d.DialAndSend(m); err != nil {
		log.Println("通知邮件发送失败,err:", err.Error())
	}

}
