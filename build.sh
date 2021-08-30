env GOOS=linux GOARCH=amd64 go build -o linux-amd64/k8tunnel
env GOOS=darwin GOARCH=amd64 go build -o macos-amd64/k8tunnel
env GOOS=windows GOARCH=amd64 go build -o windows-amd64/k8tunnel