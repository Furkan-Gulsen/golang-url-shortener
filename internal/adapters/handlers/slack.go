package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Furkan-Gulsen/golang-url-shortener/internal/config"
	"github.com/aws/aws-lambda-go/events"
	"github.com/slack-go/slack"
)

func PostMessageToSlack(ctx context.Context, message string) error {
	appConfig := config.NewConfig()
	slackToken, slackChannelID := appConfig.GetSlackParams()
	log.Print("slackToken: ", slackToken)
	log.Print("slackChannelID: ", slackChannelID)

	api := slack.New(slackToken)
	channelID, timestamp, err := api.PostMessage(
		slackChannelID,
		slack.MsgOptionText(message, false),
	)
	if err != nil {
		log.Printf("Error posting to Slack: %s", err)
		return err
	}
	log.Printf("Message successfully sent to Slack channel %s at %s", channelID, timestamp)
	return nil
}

func HandleAPIGatewayRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayProxyResponse, error) {
	err := PostMessageToSlack(ctx, "Hello world! API Gateway message.")
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500}, err
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Message successfully sent to Slack",
	}, nil
}

func HandleSQSMessage(ctx context.Context, message events.SQSMessage) error {
	return PostMessageToSlack(ctx, "Hello world! SQS message received: "+message.Body)
}

func SlackHandler(ctx context.Context, event json.RawMessage) error {
	var sqsEvent events.SQSEvent
	log.Print("sqsEvent: ", sqsEvent)
	if err := json.Unmarshal(event, &sqsEvent); err == nil && len(sqsEvent.Records) > 0 {
		for _, message := range sqsEvent.Records {
			err := PostMessageToSlack(ctx, "Hello world! SQS message received: "+message.Body)
			if err != nil {
				log.Printf("Error handling SQS message (ID: %s): %v", message.MessageId, err)
			}
		}
		log.Print("kaybettik")
		return nil
	}

	// API Gateway Event'i belirle
	var apiEvent events.APIGatewayV2HTTPRequest
	log.Print("apiEvent: ", apiEvent)
	if err := json.Unmarshal(event, &apiEvent); err == nil && apiEvent.RequestContext.HTTP.Method != "" {
		_, err := HandleAPIGatewayRequest(ctx, apiEvent)
		return err
	}

	return fmt.Errorf("invalid event type")
}
