package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

var config Config
var region string

func init() {
	region = os.Getenv("THIS_REGION")
	configJSON := fetchParameterStore(os.Getenv("CONFIG_JSON_PARAM_NAME"))
	json.Unmarshal([]byte(configJSON), &config)
}

func main() {
	lambda.Start(CWLogsToSlack)
}

// CWLogsToSlack ... Lambda function Handler
func CWLogsToSlack(logsEvent events.CloudwatchLogsEvent) error {
	if !config.validateParameter() {
		return errors.New("invalid config json from ssm")
	}

	logsData, err := logsEvent.AWSLogs.Parse()
	if err != nil {
		return err
	}

	slackMessage := makeMessage(logsData, config)

	hookURL := config.getParameter(logsData.LogGroup, logsData.LogStream, "HookURL")
	err = postToSlack(hookURL, slackMessage)
	if err != nil {
		return err
	}

	return nil
}

func fetchParameterStore(paramName string) string {
	sess := session.Must(session.NewSession())
	svc := ssm.New(
		sess,
		aws.NewConfig().WithRegion(region),
	)

	res, err := svc.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(paramName),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return *res.Parameter.Value
}

func generateLinkURL(logGroup, logStream string, logEvent events.CloudwatchLogsLogEvent) string {
	linkBase := fmt.Sprintf("https://%s.console.aws.amazon.com/cloudwatch/home?region=%s", region, region)
	if region == "us-east-1" {
		linkBase = fmt.Sprintf("https://console.aws.amazon.com/cloudwatch/home?region=%s", region)
	}
	linkLogEvent := fmt.Sprintf("?start=%d&refEventId=%s", logEvent.Timestamp, logEvent.ID)
	linkTmpl := "%s#logsV2:log-groups/log-group/%s/log-events/%s%s"
	linkURL := fmt.Sprintf(linkTmpl, linkBase, encode(logGroup, 2), encode(logStream, 2), encode(linkLogEvent, 1))
	return linkURL
}

// encode for aws console url
func encode(input string, count int) (output string) {
	for i := 0; i < count; i++ {
		input = url.QueryEscape(input)
	}
	output = strings.Replace(input, "%", "$", -1)
	return
}

func makeMessage(logsData events.CloudwatchLogsData, config Config) []byte {
	logGroup := logsData.LogGroup
	logStream := logsData.LogStream

	linkURL := generateLinkURL(logGroup, logStream, logsData.LogEvents[0])

	timestamp := logsData.LogEvents[0].Timestamp
	fields := []Field{}
	for _, logEvent := range logsData.LogEvents {
		fields = append(fields, Field{Value: logEvent.Message})
	}
	attachment := Attachment{
		Title:     "jump to log",
		TitleLink: linkURL,
		Fallback:  "LogGroup[" + logGroup + "]",
		Pretext:   "LogGroup[" + logGroup + "]",
		Color:     config.getParameter(logGroup, logStream, "Color"),
		Fields:    fields,
		Timestamp: timestamp,
	}
	slackMessage := &SlackMessage{
		Channel:     config.getParameter(logGroup, logStream, "Channel"),
		LinkNames:   1,
		Username:    config.getParameter(logGroup, logStream, "Username"),
		IconEmoji:   config.getParameter(logGroup, logStream, "IconEmoji"),
		Attachments: []Attachment{attachment},
	}

	jsonBytes, err := json.Marshal(slackMessage)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return nil
	}
	return jsonBytes
}

func postToSlack(hookURL string, jsonBytes []byte) error {
	req, err := http.NewRequest(
		"POST",
		hookURL,
		bytes.NewBuffer(jsonBytes),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
