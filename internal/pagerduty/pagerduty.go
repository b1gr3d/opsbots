package pagerduty

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	log "github.com/sirupsen/logrus"
	"os"
)

// env vars
var pagerdutyToken = os.Getenv("PAGERDUTY_TOKEN")
var pagerdutyClient = pagerduty.NewClient(pagerdutyToken)
var pdIncidentService = os.Getenv("PD_INCIDENT_SERVICE")
var pdincidentCreateEmail = os.Getenv("PD_CREATE_EMAIL")
var slackTeamID = os.Getenv("SLACK_TEAM_ID")

func CreatePagerDutyIncident(title string, slackChannelId string) (string, error) {
	newService := &pagerduty.APIReference{
		ID:   pdIncidentService, // TODO update to real pagerduty service once ready for production
		Type: "service",
	}
	incInput := &pagerduty.CreateIncidentOptions{

		Type:    "incident",
		Title:   title,
		Urgency: "low",
		Service: newService,
	}

	// create the PagerDuty Incident
	createIncident, err := pagerdutyClient.CreateIncident(pdincidentCreateEmail, incInput)
	if err != nil {
		return "", err
	}

	incidentId := createIncident.Id
	note := "https://app.slack.com/client/"+slackTeamID+"/"+slackChannelId+"/details"
	incidentNote := pagerduty.IncidentNote{
		ID:        "",
		User:      pagerduty.APIObject{Summary: pdincidentCreateEmail},
		Content:   note,
		CreatedAt: "",
	}

	// create a note on the incident which links this incident to the corresponding slack channel
	makeIncidentNote, err := pagerdutyClient.CreateIncidentNoteWithResponse(incidentId, incidentNote)
	if err != nil {
		return "", err
	}
	fmt.Println(makeIncidentNote)

	return incidentId, nil

}


func UpdateIncident(channel string, pdId string, pdNote string) {

	incidentNote := pagerduty.IncidentNote{
		ID:        "",
		User:      pagerduty.APIObject{Summary: pdincidentCreateEmail},
		Content:   pdNote,
		CreatedAt: "",
	}
	_, err := pagerdutyClient.CreateIncidentNoteWithResponse(pdId, incidentNote)
	if err != nil {
		log.Errorf("Error creating incident note: ", err)
	}
}
