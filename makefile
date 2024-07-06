run-api:
	go run ./server/apps/api/cmd/main.go

run-postgres:
	docker run -d --name assets-tracker-postgres -e POSTGRES_PASSWORD=password -p 5432:5432 postgres
