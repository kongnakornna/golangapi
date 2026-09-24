package usecase

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strconv"
	"time"
)

const gmailHost = "smtp.gmail.com"

// smtpVerifyAndSend verifies the Gmail SMTP connection on the given port and
// sends one test email. Port 465 uses implicit TLS, 587 uses STARTTLS.
func smtpVerifyAndSend(user, pass, to, subject, content string, port int) error {
	addr := gmailHost + ":" + strconv.Itoa(port)
	auth := smtp.PlainAuth("", user, pass, gmailHost)
	from := user

	if port == 465 {
		conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: gmailHost})
		if err != nil {
			return err
		}
		c, err := smtp.NewClient(conn, gmailHost)
		if err != nil {
			return err
		}
		defer c.Close()
		if err = c.Auth(auth); err != nil {
			return err
		}
		return sendMail(c, from, to, subject, content)
	}

	conn, err := net.DialTimeout("tcp", addr, 30*time.Second)
	if err != nil {
		return err
	}
	c, err := smtp.NewClient(conn, gmailHost)
	if err != nil {
		return err
	}
	defer c.Close()
	if ok, _ := c.Extension("STARTTLS"); ok {
		if err = c.StartTLS(&tls.Config{ServerName: gmailHost}); err != nil {
			return err
		}
	}
	if err = c.Auth(auth); err != nil {
		return err
	}
	return sendMail(c, from, to, subject, content)
}

func sendMail(c *smtp.Client, from, to, subject, content string) error {
	if err := c.Mail(from); err != nil {
		return err
	}
	if err := c.Rcpt(to); err != nil {
		return fmt.Errorf("rcpt %s: %w", to, err)
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	msg := "From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n\r\n" +
		content + "\r\n"
	if _, err := w.Write([]byte(msg)); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return err
	}
	return c.Quit()
}

// SensorTypes ports iot.helper list_sensor_type / list_sensor_type_th.
func (uc *settingsUseCase) SensorTypes(lang string) []map[string]string {
	en, th := sensorTypeRows()
	if lang == "th" {
		return th
	}
	return en
}

func sensorTypeRows() ([]map[string]string, []map[string]string) {
	en := []map[string]string{
		{"sensor_id": "1", "sensor_name": "Temperature"},
		{"sensor_id": "2", "sensor_name": "Humidity"},
		{"sensor_id": "3", "sensor_name": "Pressure"},
		{"sensor_id": "4", "sensor_name": "PM2.5"},
		{"sensor_id": "5", "sensor_name": "PM10"},
		{"sensor_id": "6", "sensor_name": "Wind Speed"},
		{"sensor_id": "7", "sensor_name": "Wind Direction"},
		{"sensor_id": "8", "sensor_name": "Rainfall"},
		{"sensor_id": "9", "sensor_name": "Light"},
		{"sensor_id": "10", "sensor_name": "CO2"},
	}
	th := []map[string]string{
		{"sensor_id": "1", "sensor_name": "อุณหภูมิ"},
		{"sensor_id": "2", "sensor_name": "ความชื้น"},
		{"sensor_id": "3", "sensor_name": "ความกดอากาศ"},
		{"sensor_id": "4", "sensor_name": "ฝุ่น PM2.5"},
		{"sensor_id": "5", "sensor_name": "ฝุ่น PM10"},
		{"sensor_id": "6", "sensor_name": "ความเร็วลม"},
		{"sensor_id": "7", "sensor_name": "ทิศทางลม"},
		{"sensor_id": "8", "sensor_name": "ปริมาณฝน"},
		{"sensor_id": "9", "sensor_name": "แสง"},
		{"sensor_id": "10", "sensor_name": "คาร์บอนไดออกไซด์"},
	}
	return en, th
}
