package notifier

import (
	"fmt"
	"net"
	"net/smtp"
	"strings"
)

type Config interface {
	SMTPAddress() string
	SMTPSender() string
	SMTPPassword() string
}

type SMTPNotifierService struct {
	cfg  Config
	auth *smtp.Auth
}

func NewSMTPNotifierService(cfg Config) *SMTPNotifierService {
	fields := strings.Split(cfg.SMTPAddress(), ":")
	auth := smtp.PlainAuth("", cfg.SMTPSender(), cfg.SMTPPassword(), fields[0])
	return &SMTPNotifierService{cfg: cfg, auth: &auth}
}

func (s *SMTPNotifierService) SendSuspiciousActivityMail(email string, newIP net.IP) error {
	return smtp.SendMail(
		s.cfg.SMTPAddress(),
		*s.auth,
		s.cfg.SMTPSender(),
		[]string{email},
		[]byte(s.buildSuspiciousMessage(email, newIP)),
	)
}

func (s *SMTPNotifierService) buildSuspiciousMessage(email string, ip net.IP) string {
	var messageBuilder strings.Builder

	messageBuilder.WriteString(fmt.Sprintf("From: %s\r\n", s.cfg.SMTPSender()))
	messageBuilder.WriteString(fmt.Sprintf("To: %s\r\n", email))
	messageBuilder.WriteString("Subject: Подозрительная активность!\r\n")
	messageBuilder.WriteString(fmt.Sprintf("\r\nБыл выполнен вход с ip-адреса: %s\r\n", ip.String()))

	return messageBuilder.String()
}
