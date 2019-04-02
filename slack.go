package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bluele/slack"
	"net/http"
)

func HandleRequest() (string, error) {

	var statusMsg string
	statusMsg = GetHttpStatusMessage("https://domain.com")
	statusMsg += GetHttpStatusMessage("http://domain.com")
	SlackMessage(statusMsg)

	return fmt.Sprintf("Hello!"), nil
}

func GetHttpStatusMessage(url string) string {

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

func main() {

	lambda.Start(HandleRequest)
}

func main3() {

	HandleRequest()
}
