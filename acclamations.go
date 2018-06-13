package main

import "math/rand"

var acclamations = []string {
    "Bazinga",
    "Boo-yah",
    "Cowabunga"
    "Hooray",
    "Woo-hoo",
    "Yabba-dabba-doo",
    "Yay",
    "Yee-haw",
    "Yippee"}

func getAcclamation() string {
    return acclamations[rand.Int() % len(acclamations)]
}
