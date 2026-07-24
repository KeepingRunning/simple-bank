package mail

import (
	"SimpleBank/util"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	config, err := util.LoadConfig("../")
	require.NoError(t, err)

	sender := NewGmailSender(
		config.EmailSendName,
		config.EmailSendAddress,
		config.EmailSendPassword,
	)
	subject := "Test Email"
	context := `
		<h1>Hello from SimpleBank</h1>
		<p>This is a test email sent using Gmail SMTP server.</p>
	`
	to := []string{"SimpleBankTest1@gmail.com"}
	atachmentFiles := []string{"../README.md"}

	err = sender.SendEmail(
		subject,
		context,
		to,
		nil,
		nil,
		atachmentFiles,
	)
	require.NoError(t, err)
}