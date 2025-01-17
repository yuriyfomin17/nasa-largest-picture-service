.PHONY: dc run test lint

dc:
	docker-compose up  --remove-orphans --build
run:
	HTTP_ADDR=:8080 \
	NASA_API_KEY=api_key \
	go run -race cmd/main.go