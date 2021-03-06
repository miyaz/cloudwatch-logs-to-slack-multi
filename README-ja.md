# cloudwatch-logs-to-slack-multi

CloudWatch Logsのサブスクリプションフィルタから呼び出されてSlack通知するLambda関数

## 概要

CloudWatch LogsのLambdaサブスクリプションフィルタは複数のロググループに同じ関数を指定できる  
複数ロググループから受信したログイベントをこの関数で受け取り、Slack通知する  
通知する内容 （アイコン、通知ユーザ名、通知先チャンネル）はパラメータストアに保存する設定により切り替えることができる

![App Architecture](https://github.com/miyaz/cloudwatch-logs-to-slack-multi/raw/master/cwlogs-to-slack.png)

## 設定方法

### Slack通知のためにIncoming Webhook URLを取得する

Incoming Webhook には下記２種類がある

* Slack App から作成したIncomingWebhook
  * こちらのIncomingWehookを使う場合は作成時に指定した通知先、ユーザ名、アイコンを上書きできない。そのため、通知内容ごとにIncomingWebhookを作成する必要がある  
  * 後述するJSONで個別指定しても無視される
  * 設定画面URL -> https://api.slack.com/apps
  * 設定手順URL -> https://api.slack.com/messaging/webhooks#getting_started
* Custom Integration から作成したIncomingWebhook
  * こちらのIncomingWebhookを使う場合は作成時に指定した通知内容を後述するJSONで個別指定でき、条件に一致すれば通知内容を上書きして通知できる
  * 設定画面URL -> https://{workspace}.slack.com/apps/manage/custom-integrations

### Systems Managerのパラメータストアを作成する

* Parameter Name
  * default value: /lambda/CWLogsToSlack/Configuration
* Type
  * String or SecureString
* Value
  * 下記の通りのJSON形式のテキスト

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
　　　　　・
　　　　　・
　　　　　・
  ]
}
```

* if_prefixは 処理対象のログの{ロググループ:ログストリーム}という文字列を前方一致で判定し、真であればdefault設定を個別設定で上書きした設定でSlack通知される
* ruleには if_prefix とhook_url が必須

