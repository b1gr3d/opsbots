package increportbot

import (
	"opsbots/internal/pagerduty"
	"opsbots/internal/slack"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	slackapi "github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var awarenessChannel = os.Getenv("CLIPPY_AWARE_CHANNEL_ID")
var pagerdutyUrl = os.GetEnv("PD_URL")
var ssSlack = os.GetEnv("SS_SLACK_URL")

type IncidentReportBot struct {
	EventsChannel      chan slackevents.EventsAPIInnerEvent
	CommandsChannel    chan slackapi.SlashCommand
	InteractiveChannel chan slackapi.InteractionCallback
}

func (bot IncidentReportBot) Run() {
	for {
		select {
		case myEvent := <-bot.EventsChannel:
			switch ev := myEvent.Data.(type) {
			case *slackevents.PinAddedEvent:
				channelId := ev.Channel
				AddPinnedMessageToPagerduty(channelId)
			}
		case myCommand := <-bot.CommandsChannel:
			if myCommand.Command == "/incident" {
				slack.GenerateIncidentModal(myCommand)
			}
		case myInteractive := <-bot.InteractiveChannel:
			if myInteractive.View.CallbackID == "incident-modal" {
				HandleBotIncidentCreations(myInteractive)
			}
		}
	}
}

func AddPinnedMessageToPagerduty(channel string) {
	pdId, err := slack.GetChannelTopicValue(channel, false)
	if err != nil {
		log.Errorf("Error getting Channel Topic Value", err)
	}
	pdNote, err := slack.GetPinnedMessage(channel)
	if err != nil {
		log.Errorf("Error getting pinned message from channel", err)
	}
	pagerduty.UpdateIncident(channel, pdId, pdNote)
}

func HandleBotIncidentCreations(callback slackapi.InteractionCallback) {
	userID := callback.User.ID
	userName := callback.User.Name
	// TODO if the channel fails to create, the PD incident will still be created -- need to tidy this up for that scenario
	channelName, channelId, err := slack.CreateIncidentConversation(userID)
	if err != nil {
		log.Error("Error in creating Conversation", err)
	}

	// create the PagerDuty Incident
	pagerdutyId, err := pagerduty.CreatePagerDutyIncident(channelName, channelId)

	// priority selected from incident modal submit
	selectedPriority := callback.View.State.Values["Priority"]["priority"].SelectedOption.Value

	//join Regions that were selected
	var regions []string
	region := callback.View.State.Values["Regions"]["regions"].SelectedOptions
	for _, r := range region {
		regionValue := r.Text.Text
		regions = append(regions, regionValue)
	}
	selectedRegions := strings.Join(regions, ", ")

	//join Services that were selected
	var services []string
	service := callback.View.State.Values["Services"]["services"].SelectedOptions
	for _, s := range service {

		value := s.Text.Text
		services = append(services, value)
	}
	selectedServices := strings.Join(services, ", ")

	// create initial description for Incident
	incidentDescription := callback.View.State.Values["Description"]["description"].Value

	// create the initial message sent to the slack channel when created
	initialMessage := ":fire: *Incident Created!* :fire:" + "\n" + "*Priority:*" + " " + selectedPriority +
		"\n" + "*Regions:*" + " " + selectedRegions + "\n" + "*Services:*" + " " + selectedServices + "\n" +
		"*Initial Description:*" + " " + incidentDescription + "\n" +
		"*PagerDuty:*" + " " + pagerdutyUrl + pagerdutyId

	// create the message being sent to the awarness channel
	awarenessMessage := "*Incident Created:*" + " " + "#" + channelName + "\n" +
		"*Created by:*" + " " + userName + "\n" +
		"*Go To Channel:*" + " " + ssSlack + channelId

	// do all the things
	slack.SetConversationTopic(channelId, pagerdutyId)
	slack.PostMessage(channelId, initialMessage)
	slack.PostMessage(awarenessChannel, awarenessMessage)
}
