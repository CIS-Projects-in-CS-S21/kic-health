build:
	go build -o ./bin/server ./cmd/server/server.go

push:
	docker build -t gcr.io/keeping-it-casual/kic-health:dev .
	docker push gcr.io/keeping-it-casual/kic-health:dev