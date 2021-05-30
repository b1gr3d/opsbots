package slack

import (
	"fmt"
	"strconv"

	"github.com/slack-go/slack"
)

// Section for Generating Modals that a User will fill out

func generateIncidentModal(s slack.SlashCommand) slack.ModalViewRequest {
	titleText := slack.NewTextBlockObject("plain_text", "Report a Incident", false, false)
	closeText := slack.NewTextBlockObject("plain_text", "Cancel", false, false)
	submitText := slack.NewTextBlockObject("plain_text", "Submit", false, false)

	sectionText := slack.NewTextBlockObject("mrkdwn", "Fill out the fields to report an Incident", false, false)
	sectionBlock := slack.NewSectionBlock(sectionText, nil, nil)

	descriptionText := slack.NewTextBlockObject("plain_text", "Description", false, false)
	descriptionPlaceholder := slack.NewTextBlockObject("plain_text", "Saw an increase in map3 counts...", false, false)
	descriptionElement := slack.PlainTextInputBlockElement{
		Type:         slack.METPlainTextInput,
		ActionID:     "description",
		Placeholder:  descriptionPlaceholder,
		Multiline:    true,
		InitialValue: s.Text,
	}
	descriptionBlock := slack.NewInputBlock("Description", descriptionText, descriptionElement)

	priorityLabel := slack.NewTextBlockObject("plain_text", "What is the priority?", false, false)
	priorityPlaceholder := slack.NewTextBlockObject("plain_text", "Choose wisely...", false, false)
	var priorityOptions []*slack.OptionBlockObject
	for i := 1; i <= 5; i++ {
		textString := strconv.Itoa(i)
		priorityOption := slack.OptionBlockObject{
			Text:  slack.NewTextBlockObject("plain_text", textString, false, false),
			Value: fmt.Sprintf("%d", i),
		}
		priorityOptions = append(priorityOptions, &priorityOption)
	}
	priorityElement := slack.NewOptionsSelectBlockElement("static_select", priorityPlaceholder, "priority", priorityOptions...)
	priorityBlock := slack.NewInputBlock("Priority", priorityLabel, priorityElement)

	servicesLabel := slack.NewTextBlockObject("plain_text", "Impact To:", false, false)
	servicesPlaceholder := slack.NewTextBlockObject("plain_text", "Select what is impacted...", false, false)
	var servicesOptions []*slack.OptionBlockObject
	allGame := slack.OptionBlockObject{
		Text:  slack.NewTextBlockObject("plain_text", "`All - Game`", false, false),
		Value: fmt.Sprintf("%d", "All - Game"),
	}
	allQOS := slack.OptionBlockObject{
		Text:  slack.NewTextBlockObject("plain_text", "`All - QoS`", false, false),
		Value: fmt.Sprintf("%d", "All - QoS"),
	}
	forniteGame := slack.OptionBlockObject{
		Text:  slack.NewTextBlockObject("plain_text", "`Fortnite - Game`", false, false),
		Value: fmt.Sprintf("%d", "Fortnite - Game"),
	}
	eaGame := slack.OptionBlockObject{
		Text:  slack.NewTextBlockObject("plain_text", "`EA - Game`", false, false),
		Value: fmt.Sprintf("%d", "EA - Game"),
	}
	fortniteQOS := slack.OptionBlockObject{
		Text:  slack.NewTextBlockObject("plain_text", "`Fortnite - QoS`", false, false),
		Value: fmt.Sprintf("%d", "Fortnite - QoS"),
	}
	eaQOS := slack.OptionBlockObject{
		Text:  slack.NewTextBlockObject("plain_text", "`EA - QoS`", false, false),
		Value: fmt.Sprintf("%d", "EA - QoS"),
	}
	testingQOS := slack.OptionBlockObject{
		Text:  slack.NewTextBlockObject("plain_text", "`Testing - QoS`", false, false),
		Value: fmt.Sprintf("%d", "EA - QoS"),
	}
	testingGame := slack.OptionBlockObject{
		Text:  slack.NewTextBlockObject("plain_text", "`Testing - Game`", false, false),
		Value: fmt.Sprintf("%d", "EA - QoS"),
	}
	servicesOptions = append(servicesOptions, &forniteGame, &eaGame, &fortniteQOS, &eaQOS, &testingQOS, &testingGame, &allGame, &allQOS)

	servicesElement := slack.NewOptionsSelectBlockElement("multi_static_select", servicesPlaceholder, "services", servicesOptions...)
	servicesBlock := slack.NewInputBlock("Services", servicesLabel, servicesElement)

	regionLabel := slack.NewTextBlockObject("plain_text", "Regions", false, false)
	regionPlaceholder := slack.NewTextBlockObject("plain_text", "Select impacted Regions", false, false)
	var regionOptions []*slack.OptionBlockObject
	regionNA := slack.OptionBlockObject{
		Text:  slack.NewTextBlockObject("plain_text", "NA", false, false),
		Value: fmt.Sprintf("%d", "NA"),
	}
	regionSEA := slack.OptionBlockObject{
		Text:  slack.NewTextBlockObject("plain_text", "SEA", false, false),
		Value: fmt.Sprintf("%d", "SEA"),
	}
	regionEU := slack.OptionBlockObject{
		Text:  slack.NewTextBlockObject("plain_text", "EU", false, false),
		Value: fmt.Sprintf("%d", "EU"),
	}
	regionMENA := slack.OptionBlockObject{
		Text:  slack.NewTextBlockObject("plain_text", "MENA", false, false),
		Value: fmt.Sprintf("%d", "MENA"),
	}

	regionOptions = append(regionOptions, &regionEU, &regionMENA, &regionNA, &regionSEA)

	regionElement := slack.NewOptionsSelectBlockElement("multi_static_select", regionPlaceholder, "regions", regionOptions...)
	regionBlock := slack.NewInputBlock("Regions", regionLabel, regionElement)

	regionBlock.Optional = true
	servicesBlock.Optional = true

	blocks := slack.Blocks{
		BlockSet: []slack.Block{
			sectionBlock,
			descriptionBlock,
			priorityBlock,
			servicesBlock,
			regionBlock,
		},
	}

	return slack.ModalViewRequest{
		Type:       slack.ViewType("modal"),
		Title:      titleText,
		Close:      closeText,
		Submit:     submitText,
		Blocks:     blocks,
		CallbackID: "incident-modal",
	}
}
