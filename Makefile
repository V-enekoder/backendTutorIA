env:
	copy .env.example .env

run:
	go run main.go

seed:
	go run cmd/seed/main.go