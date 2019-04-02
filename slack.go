package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bluele/slack"
	"net/http"
	"os"
	"strings"
)

var token, channel string
var urls []string

func MyLambda() {

	var statusMsg string = "Slack HTTP Monitor Lambda Golang\n"
	for _, url := range urls {
		statusMsg += GetHttpStatusMessage(url)
	}
	SlackMessage(statusMsg, token, channel)
}

func GetHttpStatusMessage(url string) string {

	fmt.Println("getting url ", url)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprint(url, "  ::  ", err, "\n")
	} else {
		return fmt.Sprint(url, " : ", resp.StatusCode, "\n")
	}
}

func SlackMessage(msg, token, channel string) {

	api := slack.New(token)
	err := api.ChatPostMessage(channel, msg, nil)
	if err != nil {
		panic(err)
	}
}

func ReadEnvironmentIntoGlobalVariables() {

	token = os.Getenv("slacktoken")
	channel = os.Getenv("slackchannel")
	urlenv := os.Getenv("urls")

	if token == "" {
		panic("Need a slack token")
	}

	if channel == "" {
		fmt.Println("No Channel, defaulting to #general")
		channel = "#general"
	}
	if channel[0] != '#' {
		channel = "#" + channel
	}

	if urlenv == "" {
		panic("Need atleast 1 url to check status")
	}
	urls = strings.Split(urlenv, " ")
}

func main() {

	ReadEnvironmentIntoGlobalVariables()
	fmt.Println(token, channel, urls)

	//SlackMessage("hello?", token, channel)
	//MyLambda()

	lambda.Start(MyLambda)
}
