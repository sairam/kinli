package kinli

import (
	"log"
	"net/mail"
	"sync"
	"time"

	"gopkg.in/gomail.v2"
)

var (
	emailCh   = make(chan *gomail.Message)
	emailOnce sync.Once
)

// EmailSMTPConfig passed for initialisation
type EmailSMTPConfig struct {
	Host string
	Port int
	User string
	Pass string
}

// EmailCtx which needs to be set to send the email
type EmailCtx struct {
	From      *mail.Address
	To        []*mail.Address
	Cc        []*mail.Address
	Bcc       []*mail.Address
	Subject   string
	PlainBody string
	HTMLBody  string

	// optional headers that need to be sent
	Headers map[string]string
}

// InitMailer can be initialised only once
// all future calls will be ignored
func InitMailer(smtpConfig *EmailSMTPConfig) {
	emailOnce.Do(func() { go smtpConfig.daemon() })
}

// TODO - add a new buffered channel or offline store to maintain and retry the list of emails that were unable to be sent
func (smtpConfig *EmailSMTPConfig) daemon() {
	var s gomail.SendCloser
	var err error
	open := false

	d := gomail.NewDialer(smtpConfig.Host, smtpConfig.Port, smtpConfig.User, smtpConfig.Pass)
	d.LocalName = "localhost"
	for {
		select {
		case m, ok := <-emailCh:
			log.Println("attempting to send the email!")
			if !ok {
				return
			}
			if !open {
				if s, err = d.Dial(); err != nil {
					log.Println("going to panic. ")
					log.Println(err)
					// panic: dial tcp: lookup email-smtp.us-east-1.amazonaws.com on 8.8.4.4:53: dial udp 8.8.4.4:53: i/o timeout
					// panic(err)
				} else {
					open = true
				}
			}
			if open {
				if err := gomail.Send(s, m); err != nil {
					log.Print(err)
				}
				log.Println("done sending the email")
			} else {
				log.Println("see above error. did not panic")
			}
			// You should close the Amazon SES within 5 seconds of next request. else you it fails with 421.
		case <-time.After(4 * time.Second):
			if open {
				if err := s.Close(); err != nil {
					log.Println("going to panic. well. not really!")
					log.Println(err)
					// panic(err)
				}
				open = false
			}
		}
	}
}

func makeFormattedAddresses(m *gomail.Message, list []*mail.Address) []string {
	arr := make([]string, 0, len(list))
	for _, to := range list {
		arr = append(arr, m.FormatAddress(to.Address, to.Name))
	}
	return arr
}

// MakeEmail makes a gomail.Message to be sent
func (ctx *EmailCtx) MakeEmail() *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", makeFormattedAddresses(m, []*mail.Address{ctx.From})...)
	m.SetHeader("To", makeFormattedAddresses(m, ctx.To)...)
	m.SetHeader("Cc", makeFormattedAddresses(m, ctx.Cc)...)
	m.SetHeader("Bcc", makeFormattedAddresses(m, ctx.Bcc)...)
	m.SetHeader("Subject", ctx.Subject)

	// NOTE: do not pass To, Cc, Bcc , Subject or Body headers
	// Headers override any value sent
	for header, value := range ctx.Headers {
		m.SetHeader(header, value)
	}
	m.SetBody("text/plain", ctx.PlainBody)
	if ctx.HTMLBody != "" {
		m.AddAlternative("text/html", ctx.HTMLBody)
	}
	return m
}

// SendEmail to be called for making and sending the email
// Never invoke SendEmail without calling InitMailer exactly once
func (ctx *EmailCtx) SendEmail() {
	emailCh <- ctx.MakeEmail()
}

// Use the channel in your program to send emails.
// add halt when required
func stop() {
	// Close the channel to stop the mail daemon.
	close(emailCh)
}
