# cloudwatch-logs-to-slack-multi

Slack notifications called from CloudWatch Logs subscription filter A Lambda function that does

## Overview

CloudWatch Logs' Lambda subscription filter can specify the same function for multiple log groups
Receive log events from multiple log groups with this function and notify Slack
The contents of the notification (icon, user name and channel to be notified) can be changed by saving the settings in the parameter store.

![App Architecture](https://github.com/miyaz/cloudwatch-logs-to-slack-multi/raw/master/cwlogs-to-slack.png)

## Configuration

### Get the Incoming Webhook URL for Slack Notifications

There are two types of Incoming Webhooks

* IncomingWebhook created from the Slack App
  * If you use IncomingWehook here, you can use the notification, username and You can't overwrite the icon. Therefore, you need to create IncomingWebhook for each notification content.
  * Even if you specify it in JSON, which will be described later, it will be ignored.
  * [Setup here](https://api.slack.com/apps)
    * [Instruction is here](https://api.slack.com/messaging/webhooks#getting_started)
* IncomingWebhook created from Custom Integration
  * If you use IncomingWebhook, you can specify individual JSON notifications specified at the time of creation, and if the conditions are met, you can overwrite the notification contents.
  * [Setup here](https://{workspace}.slack.com/apps/manage/custom-integrations)

### Creating the Systems Manager parameter store.

* Parameter Name
  * default value: /lambda/CWLogsToSlack/Configuration
* Type
  * String or SecureString
* Value
  * Text in JSON format as shown below

```
{
  "default": {
    "hook_url":"https://hooks.slack.com/services/HOGEHOGEH/****",
    "channel":"{default destination channel ID}",
    "username":"{slack username}",
    "icon_emoji":"{slack icon_emoji}",
    "color": "{default color(e.g. #D00000)}"
  },
  "rules": [
    {
      "if_prefix": "{Prefix matching string of logGroup:logStreambb}",
      "hook_url":"{Specify when overwriting}",
      "username":"{Specify when overwriting}",
      "channel":"{Specify when overwriting}",
      "color": "{Specify when overwriting}"
    },
        <interchangeable>
        <interchangeable>
        <interchangeable>
  ]
}
```

* If_prefix is a string named {loggroup:logstream} in the log to be processed and If it's true, the default setting is overridden by the individual settings, and the Slack Notified.
* If_prefix and hook_url are required in the rule.
