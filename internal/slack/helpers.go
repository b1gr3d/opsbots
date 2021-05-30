package slack

import (
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// env vars
var techopsID = os.Getenv("TECHOPS_SLACK_ID")

func PostMessage(channelID string, message string) {
	_, _, err := client.PostMessage(channelID, slack.MsgOptionText(message, false))
	if err != nil {
		log.Errorf("Error Posting Message: ", err)
	}
}

func IncidentChannelInvite(channelName string, user string) {
	_, err := client.InviteUsersToConversation(channelName, user)
	if err != nil {
		log.Errorf("Error Inviting to Channel: ", err)
	}
}

func SetConversationTopic(channelID string, incidentID string) {
	_, err := client.SetTopicOfConversation(channelID, incidentID)
	if err != nil {
		log.Errorf("Error Setting Conversation Topic: ", err)
	}
}

func GetChannelTopicValue(channel string, includeLocale bool) (string, error) {
	topicValue, err := client.GetConversationInfo(channel, includeLocale)
	if err != nil {
		return "", err
	}
	pagerdutyID := topicValue.Topic.Value
	return pagerdutyID, nil
}

func GetPinnedMessage(channel string) (string, error) {
	items, paging, err := client.ListPins(channel)
	if err != nil {
		return "", err
	}
	if reflect.DeepEqual(paging, slack.Paging{}) {
		return "", err
	}
	if err != nil {
		return "", err
	}

	return items[0].Message.Text, nil
}

func CreateIncidentConversation(user string) (string, string, error) {
	// channel setup
	seedNumber := rand.NewSource(time.Now().UnixNano())
	randomNumber := rand.New(seedNumber).Int()
	randomNumber = randomNumber % 10000
	incidentNumber := strconv.Itoa(randomNumber)
	channelDate := time.Now().Format("01-02-2006")
	channelType := "incident"
	channelName := channelType + "-" + incidentNumber + "-" + channelDate

	//create channel -- aka Conversation
	channel, err := client.CreateConversation(channelName, false)
	if err != nil {
		return "", "", err
	}
	if channel == nil {
		return "", "", err
	}

	//handling invites to the new conversation -- user who reported and techops user group
	var operationsUsers []string
	techopsMembers, err := client.GetUserGroupMembers(techopsID)
	if err != nil {
		log.Errorf("Error getting group members", err)
	}
	operationsUsers = techopsMembers
	inviteUsers := append(operationsUsers, user)
	_, err = client.InviteUsersToConversation(channel.ID, inviteUsers...)
	if err != nil {
		log.Errorf("Error inviting users to conversation: ", err)
	}

	return channelName, channel.ID, nil
}
