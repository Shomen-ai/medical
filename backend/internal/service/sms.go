// Файл: internal/service/sms.go
// Назначение: клиент SMSC.ru для отправки SMS-уведомлений о записях, переносах, отменах и напоминаниях.
package service

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SMSService — обёртка над HTTP API SMSC.ru с автоматическим no-op режимом для разработки.
type SMSService struct {
	login    string
	password string
	enabled  bool // false in development
}

// NewSMSService создаёт SMSService и автоматически отключает отправку, если креды пустые.
func NewSMSService(login, password string) *SMSService {
	return &SMSService{
		login:    login,
		password: password,
		enabled:  login != "" && password != "",
	}
}

// Send отправляет одиночное SMS через SMSC.ru; в DEV-режиме (без кредов) тихо возвращает nil.
// Send sends a single SMS message via SMSC.ru.
// Returns nil silently when SMS is disabled (empty credentials).
func (s *SMSService) Send(phone, message string) error {
	if !s.enabled {
		return nil
	}
	params := url.Values{
		"login":   {s.login},
		"psw":     {s.password},
		"phones":  {phone},
		"mes":     {message},
		"charset": {"utf-8"},
		"fmt":     {"1"}, // response: "count,cost"
	}
	resp, err := http.Post(
		"https://smsc.ru/sys/send.php",
		"application/x-www-form-urlencoded",
		strings.NewReader(params.Encode()),
	)
	if err != nil {
		return fmt.Errorf("sms send: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if strings.HasPrefix(string(body), "ERROR") {
		return fmt.Errorf("smsc error: %s", body)
	}
	return nil
}

// SendBookingCreated отправляет пациенту SMS с подтверждением созданной записи.
// SendBookingCreated notifies the patient that their appointment was created.
func (s *SMSService) SendBookingCreated(phone, doctorName string, startsAt time.Time) error {
	msg := fmt.Sprintf("Вы записаны к %s %s в %s. BeautyMed",
		doctorName,
		startsAt.Format("02.01"),
		startsAt.Format("15:04"),
	)
	return s.Send(phone, msg)
}

// SendBookingRescheduled отправляет пациенту SMS с новым временем записи.
// SendBookingRescheduled notifies the patient of a new appointment time.
func (s *SMSService) SendBookingRescheduled(phone string, startsAt time.Time) error {
	msg := fmt.Sprintf("Ваша запись перенесена на %s %s. BeautyMed",
		startsAt.Format("02.01"),
		startsAt.Format("15:04"),
	)
	return s.Send(phone, msg)
}

// SendBookingCancelled отправляет пациенту SMS об отмене записи.
// SendBookingCancelled notifies the patient that their appointment was cancelled.
func (s *SMSService) SendBookingCancelled(phone string, startsAt time.Time) error {
	msg := fmt.Sprintf("Ваша запись на %s отменена. BeautyMed", startsAt.Format("02.01 15:04"))
	return s.Send(phone, msg)
}

// SendReminder отправляет SMS-напоминание о приёме за сутки.
// SendReminder sends the 24h reminder SMS.
func (s *SMSService) SendReminder(phone, doctorName string, startsAt time.Time) error {
	msg := fmt.Sprintf("Напоминаем: завтра в %s приём у %s. BeautyMed",
		startsAt.Format("15:04"),
		doctorName,
	)
	return s.Send(phone, msg)
}
