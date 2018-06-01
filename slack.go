package main

import (
    "fmt"
    "github.com/nlopes/slack"
    "log"
    "regexp"
    "strings"
)

var rtm *slack.RTM

func handleSlackMessage(event *slack.MessageEvent) {
    text := event.Msg.Text

    if len(event.User) == 0 {
        log.Printf("#%v %v", event.Channel, text)
    } else {
        log.Printf("#%v %v: %v", event.Channel, event.User, text)
    }

    if event.Channel != getConfig("SLACK_BOT_CHANNEL_ID") {
        return
    }

    textPrefix := fmt.Sprintf("<@%s>", getConfig("SLACK_BOT_USER_ID"))

    if !strings.HasPrefix(text, textPrefix) {
        return
    }

    // Handle the 'ping' command.

    command := regexp.
            MustCompile("<[^>]+> +ping").
            FindStringSubmatch(text)

    if len(command) > 0 {
        doPing()
        return
    }

    doHelp()
}

func handleSlackMessages() {
    client := slack.New(getConfig("SLACK_BOT_TOKEN"))

    rtm = client.NewRTM()

    go rtm.ManageConnection()

    for event := range rtm.IncomingEvents {
        switch typedEvent := event.Data.(type) {
        case *slack.MessageEvent:
            handleSlackMessage(typedEvent)
        }
    }
}

func postSlackMessage(message string, arguments ...interface{}) {
    _, _, err := rtm.PostMessage(
            getConfig("SLACK_BOT_CHANNEL_ID"),
            fmt.Sprintf(message, arguments...),
            slack.NewPostMessageParameters())

    if err != nil {
        panic(err)
    }
}
