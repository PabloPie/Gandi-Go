workflow "Go tests" {
  on = "pull_request"
  resolves = ["Golang test"]
}

action "Golang test" {
  uses = "cedrickring/golang-action@1.2.0"
  args = "go get -t -v && go build && go test -cover"
  env = {
    IMPORT = "github.com/PabloPie/Gandi-Go/hosting"
  }
}
