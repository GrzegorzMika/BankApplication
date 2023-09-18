package mail

import (
	"testing"

	"BankApplication/internal/util"

	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	t.Skip()
	config, err := util.LoadConfig("../..")
	require.NoError(t, err)
	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
		<h1>Hello World!</h1>
		<p>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore</p>
		`
	to := []string{}
	attachFiles := []string{"./Lorem.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
