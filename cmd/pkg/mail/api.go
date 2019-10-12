package mail

import (
	"io"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/gomail.v2"

	"github.com/fidelfly/fxgo/logx"
)

type MessageDecorator func(message *gomail.Message)

//export
func CreateMessage(decorators ...MessageDecorator) *gomail.Message {
	m := gomail.NewMessage()
	if len(decorators) > 0 {
		for _, d := range decorators {
			d(m)
		}
	}
	return m
}

//export
func TemplateMessage(ns string, name string, data interface{}) MessageDecorator {
	return func(message *gomail.Message) {
		message.AddAlternativeWriter("text/html", func(writer io.Writer) error {
			t := GetTemplate(ns, name)
			if t != nil {
				return t.Execute(writer, data)
			}
			return nil
		})
	}
}

//export
func Meta(subject string, from string, to ...string) MessageDecorator {
	return func(message *gomail.Message) {
		message.SetHeader("Subject", subject)
		message.SetHeader("From", from)
		message.SetHeader("To", to...)
	}
}

func Subject(subject string) MessageDecorator {
	return func(message *gomail.Message) {
		message.SetHeader("Subject", subject)
	}
}

func To(to ...string) MessageDecorator {
	return func(message *gomail.Message) {
		message.SetHeader("To", to...)
	}
}

func From(from string) MessageDecorator {
	return func(message *gomail.Message) {
		message.SetHeader("From", from)
	}
}

func Attachment(filename string, path string) MessageDecorator {
	return func(message *gomail.Message) {
		if _, err := os.Stat(path); err != nil {
			return
		}

		message.Attach(filename, gomail.SetCopyFunc(func(w io.Writer) error {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = w.Write(data)
			return err
		}))
	}
}

func EmbedFile(filename string, path string) MessageDecorator {
	return func(message *gomail.Message) {
		if _, err := os.Stat(path); err != nil {
			return
		}

		message.Embed(filename, gomail.SetCopyFunc(func(w io.Writer) error {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			_, err = w.Write(data)
			if err != nil {
				return err
			}
			return nil
		}))
	}
}

//export
func CreateTemplateMessage(ns string, name string, data interface{}) *gomail.Message {
	return CreateMessage(TemplateMessage(ns, name, data))
}

//export
func SendMail(mail *gomail.Message) (err error) {
	err = mailDialer.DialAndSend(mail)
	if err != nil {
		subjectHeader := mail.GetHeader("Subject")
		toHeader := mail.GetHeader("To")
		if len(subjectHeader) > 0 {
			logx.Errorf("Send mail failed! [subject=%s, to=%s]", subjectHeader[0], strings.Join(toHeader, ","))
		} else {
			logx.Errorf("Send mail failed! [to=%s]", strings.Join(toHeader, ","))
		}

	}
	return
}
