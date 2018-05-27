package main

import (
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func GetEmailBodies(name string, startDate time.Time, events []CalendarEvent) (template.HTML, string) {
	data := struct {
		Name      string
		Events    []CalendarEvent
		StartDate time.Time
	}{
		name,
		events,
		startDate,
	}

	var htmlBuffer, textBuffer bytes.Buffer
	err := template.Must(template.ParseFiles("tmpl/cron/email.html")).Execute(&htmlBuffer, data)
	if err != nil {
		panic(err)
	}
	err = template.Must(template.ParseFiles("tmpl/cron/email.txt")).Execute(&textBuffer, data)
	if err != nil {
		panic(err)
	}

	return template.HTML(htmlBuffer.String()), textBuffer.String()
}

func SendReminderEmail(startDate time.Time, events []CalendarEvent) {
	defer trace(traceName(fmt.Sprintf("SendReminderEmail(%v, %v)", startDate, events)))

	SendReminderEmailToUser(startDate, events, "Sarah", "sarahvictoriayoung@gmail.com")
	SendReminderEmailToUser(startDate, events, "Ian", "ian@lucidhelix.com")
}

func SendReminderEmailToUser(startDate time.Time, events []CalendarEvent, name string, email string) {
	defer trace(traceName(fmt.Sprintf("SendReminderEmailToUser(%v, %v, %s, %s)", startDate, events, name, email)))

	subject := fmt.Sprintf("Family Reminders - %s", startDate.Format("Mon, Jan 2, 2006"))
	htmlBody, textBody := GetEmailBodies(name, startDate, events)

	awsConfig := aws.NewConfig().WithRegion("us-east-1").WithCredentials(credentials.NewStaticCredentials(config.awsAccessKey, config.awsSecret, ""))

	svc := ses.New(session.New(awsConfig))
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(email),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(string(htmlBody)),
				},
				Text: &ses.Content{
					Charset: aws.String("UTF-8"),
					Data:    aws.String(textBody),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("UTF-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String("family@icadev.com"),
	}

	result, err := svc.SendEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			case ses.ErrCodeConfigurationSetSendingPausedException:
				fmt.Println(ses.ErrCodeConfigurationSetSendingPausedException, aerr.Error())
			case ses.ErrCodeAccountSendingPausedException:
				fmt.Println(ses.ErrCodeAccountSendingPausedException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}

	fmt.Println(result)
}
