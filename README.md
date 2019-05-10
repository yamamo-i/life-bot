# life-bot

## 概要

* Slack上で動かすbotのgolang実行版

## 実行

### ローカルでの実行(macOSでのみ検証済み)

* 環境変数の設定
```sh
# Slackのbotが利用できるtoken
$ export BOT_TOKEN=hoge
# RAKUTEN APIを利用する時のaccount_id
$ export RAKUTEN_ID=hoge
# Slack上で呼ばれるbotの名前
$ export BOT_NAME=botname
# Slack上のbotのuser_id
$ export BOT_ID=ID
```
* 実行
```sh
$ go run *.go
```


### Dockerを用いた実行

* build方法  
`docker build -t lifebot:latest .`
* 実行方法
```sh
$ cp envfile.template envfile
# envfileは必要な環境変数を設定
$ vim envfile
$ docker run --envfile=envfile lifebot:latest
