# cwLogToSlack

CloudWatch Logsのサブスクリプションフィルタから呼び出されてSlack通知するLambda関数

## 概要

CloudWatch LogsのLambdaサブスクリプションフィルタは複数のロググループに同じ関数を指定できる
複数ロググループから受信したログイベントをこの関数で受け取り、設定に基づいて、通知する内容
（アイコン、通知ユーザ名、通知先チャンネル）を切り替えられる



指定された情報（アイコン、通知ユーザ、通知先チャンネル）



## 初回セットアップ




## デプロイ方法

初回セットアップのみAWS SAM CLIを使用し、Lambda+IAMRole+




## 設定方法

Systems Managerのパラメータストア

/lambda/CWLogsToSlack/Configuration

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


SlackのIncommingWebhookには２種類あります。
旧式
https://sencorp-group.slack.com/apps/manage/custom-integrations
Incoming WebHooksをクリックして作成する


新式
https://api.slack.com/apps
でApp作ってActivate Incoming WebhooksしてAdd New Webhook to Workspaceボタン で作るやつ
新式の場合はチャンネル、ユーザ名（AppName）、アイコンが固定なので
rulesでprefixマッチしても color以外(channel,icon_emoji,username)を
指定しても無視されます


## 動作確認

sam build
sam local invoke --event event.json
sam deploy --guided

