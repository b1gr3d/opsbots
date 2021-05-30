package main

import (
	"net/http"
	"opsbots/internal/increportbot"
	"opsbots/internal/slack"

	log "github.com/sirupsen/logrus"
	slackapi "github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

func main() {

	// create gochannels for business logic to use
	eventsChannel := make(chan slackevents.EventsAPIInnerEvent)
	commandsChannel := make(chan slackapi.SlashCommand)
	interactiveChannel := make(chan slackapi.InteractionCallback)

	// set gochannels in appropriate packages for interaction
	slack.EventChannel = eventsChannel
	slack.CommandChannel = commandsChannel
	slack.InteractiveChannel = interactiveChannel

	// set the struct and start a thread
	irBot := increportbot.IncidentReportBot{
		EventsChannel:      eventsChannel,
		CommandsChannel:    commandsChannel,
		InteractiveChannel: interactiveChannel,
	}

	go irBot.Run()

	// http handle funcs for server
	http.HandleFunc("/events", slack.EventHandler)
	http.HandleFunc("/command", slack.CommandHandler)
	http.HandleFunc("/interactive-endpoint", slack.InteractiveHandler)

	// kubernetes liveness and readiness probes; just return status 200 (OK)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {w.WriteHeader(http.StatusOK)})

	//start the server and listen
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatalf("Got Error: %v", err)
	}

}
