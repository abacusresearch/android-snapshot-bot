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

    result := fmt.Sprintf("%s! ", getAcclamation())

    if notification.PullRequest == "" {
        result += fmt.Sprintf(
            "We've built *%s@%s*:",
            notification.Repository, notification.Branch)
    } else {
        result += fmt.Sprintf(
            "We've built pull request *%s* for *%s@%s*:",
            notification.PullRequest, notification.Repository, notification.PullRequestBranch)
    }

    commits := len(notification.Commits)

    for index, commit := range notification.Commits {
        if commits < 6 || index < 2 || index > commits - 3 {
            result += "\n" + commit
        } else if index == 2 {
            result += "\n···"
        }
    }

    for _, url := range notification.Urls {
        result += "\n:package: " + url
    }

    postSlackMessage(result)
}
