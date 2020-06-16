package main

import (
	"reflect"
	"strings"
)

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
	Default Parameter   `json:"default"`
	Rules   []Parameter `json:"rules"`
}

type Parameter struct {
	IfPrefix  string `json:"if_prefix"`
	HookURL   string `json:"hook_url"`
	Channel   string `json:"channel"`
	Username  string `json:"username"`
	IconEmoji string `json:"icon_emoji"`
	Color     string `json:"color"`
}

// get parameter value of specified fieldName using reflection
//   ref. https://leben.mobi/go/reflect/go-programming/structure/#FieldByName
func (c Config) getParameter(logGroup, logStream, fieldName string) (value string) {
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

func (p Parameter) getValue(fieldName string) string {
	rvPrm := reflect.ValueOf(p)
	v := rvPrm.FieldByName(fieldName).Interface()
	if _, ok := v.(string); ok {
		return v.(string)
	}
	return ""
}
