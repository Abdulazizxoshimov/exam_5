package v1

import (
	"math/rand"
	"net/smtp"
    "exam_5/admin_api_gateway/api/models"
	"strconv"
	"time"
)

func RadomGenerator()int{
    rand.Seed(time.Now().UnixNano())

    randomNumber := rand.Intn(900000) + 100000
    return randomNumber
}

func SendCodeGmail(user models.CreateUser)(string, error){
    email := "abdulazizxoshimov22@gmail.com"
    password := "hxytgczqprxfsltu "

    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    auth := smtp.PlainAuth("test", email, password, smtpHost)


    randomNumber := RadomGenerator()
    randomNumberString := strconv.Itoa(randomNumber)


    to := []string{user.Email}
    msg := []byte(randomNumberString)

    err := smtp.SendMail(smtpHost+":"+smtpPort, auth, email, to, msg)
    if err != nil {
        return  "",err
    }
    return randomNumberString, nil
}

