package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

const (
	region       = "ap-northeast-1"
	paramEnvName = "CONFIG_JSON_PARAM_NAME"
)

var config Config

func init() {
	configJSON := fetchParameterStore(os.Getenv(paramEnvName))
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
	config.LogGroup = logsData.LogGroup
	config.LogStream = logsData.LogStream

	slackMessage := makeMessage(logsData, config)

	hookURL := config.getParameter("HookURL")
	postToSlack(hookURL, slackMessage)
	if err != nil {
		fmt.Println(err)
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
	escapedLogGroup := strings.Replace(logGroup, "/", "$252F", -1)
	escapedLogStream := strings.Replace(logStream, "/", "$252F", -1)

	linkBase := fmt.Sprintf("https://%s.console.aws.amazon.com/cloudwatch/home?region=%s", region, region)
	linkTmpl := "%s#logsV2:log-groups/log-group/%s/log-events/%s$3Fstart$3D%d$26refEventId$3D%s"
	linkURL := fmt.Sprintf(linkTmpl, linkBase, escapedLogGroup, escapedLogStream, logEvent.Timestamp, logEvent.ID)
	return linkURL
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
		Color:     config.getParameter("Color"),
		Fields:    fields,
		Timestamp: timestamp,
	}

	slackMessage := &SlackMessage{
		Channel:     config.getParameter("Channel"),
		LinkNames:   1,
		Username:    config.getParameter("Username"),
		IconEmoji:   config.getParameter("IconEmoji"),
		Attachments: []Attachment{attachment},
	}
	fmt.Printf("%v", slackMessage)

	jsonBytes, err := json.Marshal(slackMessage)
	if err != nil {
		fmt.Println("JSON Marshal error:", err)
		return nil
	}
	fmt.Println(string(jsonBytes))
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

type SlackMessage struct {
	Channel     string       `json:"channel"`
	LinkNames   int          `json:"link_names"`
	Username    string       `json:"username"`
	IconEmoji   string       `json:"icon_emoji"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Title     string  `json:"title"`
	TitleLink string  `json:"title_link"`
	Fallback  string  `json:"fallback"`
	Pretext   string  `json:"pretext"`
	Color     string  `json:"color"`
	Fields    []Field `json:"fields"`
	Timestamp int64   `json:"ts"`
}

type Field struct {
	Value string `json:"value"`
}

type Config struct {
	Default   Param   `json:"default"`
	Rules     []Param `json:"rules"`
	LogGroup  string
	LogStream string
}

type Param struct {
	IfPrefix  string `json:"if_prefix"`
	HookURL   string `json:"hook_url"`
	Channel   string `json:"channel"`
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	Color     string `json:"color"`
}

// フィールド名を指定して必要な値を返す
// フィールド名から値を取り出すのにリフレクションを使っている
//   ref. https://leben.mobi/go/reflect/go-programming/structure/#FieldByName
func (c Config) getParameter(fieldName string) (value string) {
	logGroup := c.LogGroup
	logStream := c.LogStream
	value = c.Default.getValue(fieldName)

	logGroupStream := logGroup + ":" + logStream
	for _, rule := range c.Rules {
		tmpValue := rule.getValue(fieldName)
		if rule.IfPrefix != "" && tmpValue != "" {
			if strings.HasPrefix(logGroupStream, rule.IfPrefix) {
				value = tmpValue
				return
			}
		}
	}
	return
}

func (c Config) validateParameter() bool {
	if c.Default.HookURL == "" {
		return false
	}
	return true
}

func (p Param) getValue(fieldName string) string {
	rvPrm := reflect.ValueOf(p)
	v := rvPrm.FieldByName(fieldName).Interface()
	if _, ok := v.(string); ok {
		return v.(string)
	}
	return ""
}
