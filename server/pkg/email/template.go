package email

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
)

type PasswordResetEmailData struct {
	Username          string
	OTP               string
	ExpirationMinutes string
}

func GeneratePasswordResetEmail(data PasswordResetEmailData) (string, error) {
	tmpl, err := parseTemplate("password_reset_email.html")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func parseTemplate(path string) (*template.Template, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	tmplPath := filepath.Join(wd, "templates", path)
	return template.ParseFiles(tmplPath)
}
