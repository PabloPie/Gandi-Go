workflow "Go tests" {
  on = "pull_request"
  resolves = ["Golang test"]
}

action "Golang test" {
  uses = "cedrickring/golang-action@1.2.0"
  args = "go get -t github.com/kolo/xmlrpc && go build && go test"
}
