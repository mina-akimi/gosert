fix:
		go mod tidy
		goimports -l -w $$(find . -type f -name '*.go' -not -path "./vendor/*")
		go vet ./...

test:
		go test ./...