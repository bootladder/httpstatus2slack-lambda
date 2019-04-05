package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/bluele/slack"
	"net/http"
	"os"
	"strings"
	"time"
)

var token, channel string
var urls []string

func MyLambda() {

	if time_for_daily_status() {
		fmt.Println("Time for daily status!")

		var statusMsg string = "Slack HTTP Monitor Lambda Golang\n"
		for _, url := range urls {
			msg, _ := GetHttpStatusMessage(url)
			statusMsg += msg
		}
		SlackMessage(statusMsg, token, channel)

	} else { //check for errors only
		var isAnyWebsiteDown bool
		var statusMsg string = "Slack HTTP Monitor Lambda Golang\n"
		for _, url := range urls {
			msg, err := GetHttpStatusMessage(url)
			if err != nil {
				statusMsg += msg
				isAnyWebsiteDown = true
			}
		}

		if isAnyWebsiteDown {
			fmt.Println("A site was down!  Sending slack message")
			SlackMessage(statusMsg, token, channel)
		} else {
			fmt.Println("All good!")
		}

	}
}

func time_for_daily_status() bool {

	currentTime := time.Now()
	timeString := fmt.Sprint("Short Hour Minute Second: ",
		currentTime.Format("3:4:5"))
	fmt.Println("The time is: ", timeString)

	timeSetpoint := "10:39"
	if strings.HasPrefix(timeString, timeSetpoint) {

		return true
	}
	return false

}

func GetHttpStatusMessage(url string) (string, error) {

	fmt.Println("getting url ", url)
	resp, err := http.Get(url)
	var message string
	if err != nil {
		message = fmt.Sprint(url, "  ::  ", err, "\n")
	} else {
		message = fmt.Sprint(url, " : ", resp.StatusCode, "\n")
	}

	return message, err
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
