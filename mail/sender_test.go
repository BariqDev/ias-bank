package mail

import (
	"testing"

	"github.com/BariqDev/ias-bank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {
	config, err := util.LoadConfig("../")

	require.NoError(t, err)
	sender := NewGmailSender("islam said", config.VerifyEmailAddress, config.VerifyEmailPassword)

	subject := "Test send email"
	content := `
	<h1>Voila congratulations</h1>
	<p>thanks for your hard work to i hope you see this mail in future and say voila i did it </p>
	`

	to := []string{"islam.said.dev@gmail.com"}

	attachedFiles := []string{"../readme.md"}
	err = sender.SendEmail(subject, content, to, nil, nil, attachedFiles)
	require.NoError(t, err)
}
