package pkg

import (
	"fmt"
	"strconv"

	"gopkg.in/gomail.v2"
)

type Mailer interface {
	SendResetLink(to string, link string) error
}

type gomailMailer struct {
	host     string
	port     int // Gomail menggunakan tipe data int untuk port
	username string
	password string
	from     string
}

// NewGomailMailer menginisialisasi konfigurasi Gomail
func NewGomailMailer(host, portStr, username, password, from string) Mailer {
	// Konversi port dari string (.env) ke int
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 587 // Default port jika konversi gagal
	}

	return &gomailMailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (m *gomailMailer) SendResetLink(to string, link string) error {
	msg := gomail.NewMessage()

	// Menyusun Header Email
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", "Reset Password Akun Anda")

	// Menyusun Konten HTML (Tanpa perlu menulis MIME manual)
	body := fmt.Sprintf(`
		<html>
			<body style="font-family: Arial, sans-serif; color: #333;">
				<h2 style="color: #2563eb;">Permintaan Reset Password</h2>
				<p>Kami menerima permintaan untuk mereset password akun Anda.</p>
				<p>Silakan klik tautan di bawah ini untuk mengganti password Anda (Tautan berlaku selama 15 menit):</p>
				<div style="margin: 20px 0;">
					<a href="%s" style="background-color: #3b82f6; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; display: inline-block; font-weight: bold;">Reset Password</a>
				</div>
				<p>Jika tombol di atas tidak berfungsi, salin dan tempel tautan berikut ke browser Anda:</p>
				<p style="color: #2563eb; word-break: break-all;">%s</p>
				<hr style="border: 0; border-top: 1px solid #e5e7eb; margin: 20px 0;" />
				<p style="font-size: 12px; color: #9ca3af;">Jika Anda tidak merasa meminta ini, abaikan email ini.</p>
			</body>
		</html>`, link, link)

	msg.SetBody("text/html", body)

	// Membuat dialer SMTP Gomail
	dialer := gomail.NewDialer(m.host, m.port, m.username, m.password)

	// Opsional: Jika menggunakan server lokal/self-signed SSL (seperti MailHog/Mailtrap tertentu),
	// Anda bisa mengaktifkan baris di bawah ini untuk melewati validasi sertifikat:
	// dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Kirim Email
	if err := dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("gomail gagal mengirim email: %w", err)
	}

	return nil
}
