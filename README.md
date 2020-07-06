# cloudwatch-logs-to-slack-multi

CloudWatch Logsのサブスクリプションフィルタから呼び出されてSlack通知するLambda関数

## 概要

CloudWatch LogsのLambdaサブスクリプションフィルタは複数のロググループに同じ関数を指定できる
複数ロググループから受信したログイベントをこの関数で受け取り、Slack通知する
通知する内容 （アイコン、通知ユーザ名、通知先チャンネル）はパラメータストアに保存する
設定により切り替えることができる

![App Architecture](https://github.com/miyaz/cloudwatch-logs-to-slack-multi/raw/master/cwlogs-to-slack.svg)

## 設定方法

1. Slack通知のためにIncoming Webhook URLを取得する

Incoming Webhook には下記２種類がある

* Slack App から作成したIncomingWebhook
  * こちらのIncomingWehookを使う場合は作成時に指定した通知先、ユーザ名、アイコンを
  　上書きできない。そのため、通知内容ごとにIncomingWebhookを作成する必要がある
  　後述するJSONで個別指定しても無視される
  * https://api.slack.com/apps
  * https://api.slack.com/messaging/webhooks#getting_started
* Custom Integration から作成したIncomingWebhook
  * こちらのIncomingWebhookを使う場合は作成時に指定した通知内容を後述するJSONで
  　個別指定できる。条件に一致すれば通知内容を上書きして通知できる
  * https://{workspace}.slack.com/apps/manage/custom-integrations

2. Systems Managerのパラメータストアを作成する

パラメータ名：/lambda/CWLogsToSlack/Configuration
種類は：String or SecureString のいづれか
設定値：下記の通りのJSON形式のテキスト

```
{
  "default": {
    "hook_url":"https://hooks.slack.com/services/HOGEHOGEH/****",
    "channel":"{デフォルトの通知先チャンネルID}",
    "username":"{サービス／ステージを識別する名称}",
    "icon_emoji":"{サービス／ステージを識別する絵文字アイコン}",
    "color": "{デフォルトのカラー #D00000}",
  },
  "rules": [
    {
      "if_prefix": "{logGroup:logStreamを前方一致で判定するPrefix}",
      "hook_url":"{デフォルト値を上書きしたい場合に指定}",
      "channel":"{デフォルト値を上書きしたい場合に指定}",
      "color": "{デフォルト値を上書きしたい場合に指定}"
    },
　　　　　・
　　　　　・
　　　　　・
  ]
}
```

if_prefixは 処理対象のログの{ロググループ:ログストリーム}という文字列を前方一致で判定し
真であれば default設定を個別設定で上書きした設定でSlack通知される

ruleには if_prefix とhook_url が必須

