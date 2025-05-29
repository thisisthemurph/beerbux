package email

import (
	"fmt"
	"github.com/yosssi/gohtml"
	"log/slog"
)

type TerminalEmailSender struct {
	logger *slog.Logger
}

func NewTerminalEmailLogger(logger *slog.Logger) Sender {
	return &TerminalEmailSender{
		logger: logger,
	}
}

func (f *TerminalEmailSender) Send(to, subject, html string) (string, error) {
	f.logger.Info("FakeEmail sent", "to", to, "subject", subject)
	fmt.Println(gohtml.Format(html))
	return "", nil
}
