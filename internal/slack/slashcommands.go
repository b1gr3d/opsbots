package slack

import (
	slackapi "github.com/slack-go/slack"
)

func GenerateIncidentModal(s slackapi.SlashCommand) {
	modalRequest := generateIncidentModal(s)
	client.OpenView(s.TriggerID, modalRequest)
}
