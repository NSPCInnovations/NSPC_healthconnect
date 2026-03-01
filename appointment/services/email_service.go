package services

import (
	"log"

	"gopkg.in/gomail.v2"
)

const (
	senderEmail = "ameermahammad40@gmail.com"
	appPassword = "cbyogiyizhrepyhg"
)

func SendEmailWithPDF(toEmail, filePath, patientName string) {
	m := gomail.NewMessage()
	m.SetHeader("From", "ameermahammad40@gmail.com")
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Booking Confirmed!")
	m.SetBody("text/html", "Dear "+patientName+",<br><br>Your appointment is <b>Confirmed</b>. Receipt is attached.")
	m.Attach(filePath)

	// must use app password for Gmail SMTP authentication
	d := gomail.NewDialer("smtp.gmail.com", 587, senderEmail, appPassword)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Email Error: %v", err)
	} else {
		log.Println("Email sent successfully to:", toEmail)
	}
}
