workflow "Go tests" {
  resolves = ["Golang test"]
  on = "push"
}

action "Golang test" {
  uses = "cedrickring/golang-action@1.2.0"
  args = "go get -t -v && go build && go test ./... -coverprofile=coverage.txt -covermode=atomic && curl -s https://codecov.io/bash | bash"
  env = {
    PROJECT_PATH = "./hosting/hostingv4"
  }
  secrets = ["CODECOV_TOKEN"]
}
