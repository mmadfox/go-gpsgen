.PHONY: proto 
proto:
	protoc proto/gpsgen.proto --js_out=import_style=commonjs,binary:.
	protoc proto/gpsgen.proto --go_out=.
	protoc proto/gpsgen.proto --ts_out=.
	protoc proto/snapshot.proto --go_out=.

.PHONY: cover 
cover:
	go test ./... -cover -coverprofile=cover.out
	go tool cover -func=cover.out
	go tool cover -html=cover.out
