check:
		goimports -l -w $$(find . -type f -name '*.go' -not -path "./vendor/*")
		go vet $$(glide novendor)

test:
		go test ./...