package mail

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {

	sender := NewGmailSender("islam said", "contact.islam.said@gmail.com", "ugdcrexwabqdokfm")

	subject := "Test send email"
	content := `
	<h1>Voila congratulations</h1>
	<p>thanks for your hard work to i hope you see this mail in future and say voila i did it </p>
	`

	to := []string{"islam.said.dev@gmail.com"}

	attachedFiles := []string{"../readme.md"}
	err := sender.SendEmail(subject, content, to, nil, nil, attachedFiles)
	require.NoError(t, err)
}
