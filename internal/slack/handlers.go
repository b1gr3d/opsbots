package slack

import (
	"encoding/json"
	"fmt"
	slackapi "github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// env vars
var slackToken = os.Getenv("CLIPPY_SLACK_TOKEN")
var client = slackapi.New(slackToken)

// gochannel vars
var EventChannel chan slackevents.EventsAPIInnerEvent
var CommandChannel chan slackapi.SlashCommand
var InteractiveChannel chan slackapi.InteractionCallback

// handler for slack events
func EventHandler(w http.ResponseWriter, r *http.Request){
	//read request body from the request
	body, err := ioutil.ReadAll(r.Body)
	eventsAPIEvent, err := slackevents.ParseEvent(body, slackevents.OptionNoVerifyToken())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	//url verification needed for responding to Slack's challenge sent to me
	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))

	}
	// send to event channel
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		EventChannel <- eventsAPIEvent.InnerEvent
	}
}

// handler for slack commands
func CommandHandler(w http.ResponseWriter, r *http.Request){
	slashCommand, err := slackapi.SlashCommandParse(r)
	if err != nil {
		log.Fatal("Got err %+v\n", err)
	}
	// send to command channel
	CommandChannel <- slashCommand
}

// handler for interactives from slack -- modals, callbacks, etc
func InteractiveHandler(w http.ResponseWriter, r *http.Request){
	var i slackapi.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &i)
	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// send to interactive channel
	InteractiveChannel <- i
}