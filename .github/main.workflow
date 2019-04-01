workflow "Go tests" {
  resolves = ["Golang Action"]
  on = "pull_request"
}

action "Golang Action" {
  uses = "cedrickring/golang-action@1.2.0"
}
