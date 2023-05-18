cover:
	go test ./... -cover -coverprofile=cover.out
	go tool cover -func=cover.out
	go tool cover -html=cover.out
	