workflow "Go tests" {
  resolves = ["Golang test"]
  on = "push"
}

action "Golang test" {
  uses = "cedrickring/golang-action@1.2.0"
  args = "go get -t -v && go build && go test ./... -cover"
}
