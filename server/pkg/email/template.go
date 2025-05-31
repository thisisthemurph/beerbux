package email

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
)

type ResetPasswordEmailData struct {
	Username          string
	OTP               string
	ExpirationMinutes string
}

func GenerateResetPasswordEmail(data ResetPasswordEmailData) (string, error) {
	tmpl, err := parseTemplate("reset_password_email.html")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

type UpdatePasswordEmailData struct {
	Username          string
	OTP               string
	ExpirationMinutes string
}

func GenerateUpdatePasswordEmail(data UpdatePasswordEmailData) (string, error) {
	tmpl, err := parseTemplate("update_password_email.html")
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

type UpdateEmailAddressData struct {
	Username          string
	NewEmail          string
	OTP               string
	ExpirationMinutes string
}

func GenerateUpdateEmailAddressEmail(data UpdateEmailAddressData) (string, error) {
	tmpl, err := parseTemplate("update_email_address_email.html")
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
