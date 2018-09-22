package main

import (
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    "io/ioutil"
    "log"
    "net/http"
)

type Notification struct {
    Branch string `json:"branch"`
    Commits []string `json:"commits"`
    PullRequest string `json:"pullRequest"`
    PullRequestBranch string `json:"pullRequestBranch"`
    Repository string `json:"repository"`
    Urls []string `json:"urls"`
}

func doHelp() {
    postSlackMessage("Sorry, I don't understand.")
}

func doPing() {
    postSlackMessage("Pong.")
}

func main() {
    go serve()

    handleSlackMessages()
}

func serve() {
    router := mux.NewRouter().StrictSlash(true)

    router.HandleFunc("/notifications", serveNotifications)

    log.Fatal(http.ListenAndServe(":8080", router))
}

func serveNotifications(response http.ResponseWriter, request *http.Request) {
    serviceUser, servicePassword, _ := request.BasicAuth()

    if serviceUser != getConfig("SERVICE_USER") || servicePassword != getConfig("SERVICE_PASSWORD") {
        response.WriteHeader(401)
        return
    }

    var notification Notification

    notificationData, err := ioutil.ReadAll(request.Body)

    defer request.Body.Close()

    if err != nil {
        panic(err)
    }

    err = json.Unmarshal(notificationData, &notification)

    if err != nil {
        panic(err)
    }

    message := fmt.Sprintf("%s! ", getAcclamation())

    if notification.PullRequest == "" {
        message += fmt.Sprintf(
            "We've built *%s@%s*.",
            notification.Repository, notification.Branch)
    } else {
        message += fmt.Sprintf(
            "We've built pull request *%s* for *%s@%s*.",
            notification.PullRequest, notification.Repository, notification.PullRequestBranch)
    }

    commits := len(notification.Commits)

    attachedText := ""
    attachedUrls := notification.Urls

    if commits > 4 {
        attachedText += "···"
    }

    for index, commit := range notification.Commits {
        if index > commits - 5 {
            if len(attachedText) != 0 {
                attachedText += "\n"
            }

            attachedText += commit
        }
    }

    postSlackMessageWithAttachments(message, attachedText, attachedUrls)
}
