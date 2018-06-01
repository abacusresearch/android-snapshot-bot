package main

import "math/rand"

var acclamations = []string {
    "Bazinga",
    "Boo-ya",
    "Hooray",
    "Woo-hoo",
    "Yay",
    "Yee-haw",
    "Yippee"}

func getAcclamation() string {
    return acclamations[rand.Int() % len(acclamations)]
}
