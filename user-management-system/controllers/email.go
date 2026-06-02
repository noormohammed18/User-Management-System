package controllers

import "gopkg.in/gomail.v2"

func SendOTP(email, otp string) error {

	m := gomail.NewMessage()

	m.SetHeader("From", "noor4u2003@gmail.com")
	m.SetHeader("To", email)
	m.SetHeader("Subject", "OTP Verification")

	m.SetBody(
		"text/plain",
		"Your OTP is: "+otp,
	)

	d := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		"noor4u2003@gmail.com",
		"hfbp ilil ifyg auzu",
	)

	return d.DialAndSend(m)
}