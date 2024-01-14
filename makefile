.PHONY: help
include .env.local
export

refresh_deps:
	@go mod tidy

create_test_img:
	@docker rm --force TestBackend
	@docker build -t backend:test -f Dockerfile.testing .
	@docker create -p 5666:5666 --name TestBackend --network DB backend:test 

start_db:
	@echo "starting db"
	@docker start ContentDB
	@sleep 3
	@echo "db up and running"

generate_data: 
	@echo "generating data"
	@docker exec -i ContentDB mysql -u root -proot < ./devStuff/generate_data.sql
	@echo "data generated"

reset_db: 
	@echo "reseting db"
	@docker exec -i ContentDB mysql -u root -proot < ./devStuff/reset_db.sql
	@echo "db reseted"

setup_db:
	@echo "setting db up"
	@docker exec -i ContentDB mysql -u root -proot < ./devStuff/setup_db.sql
	@echo "db setuped"



build_dataservice: refresh_deps
	@echo "Building dataservice"
	@go build -o ./bin/dataserviceexe ./cmd/dataservice/main.go

run_dataservice: build_dataservice
	@echo "Running dataservice"
	@./bin/dataserviceexe

build_frontend: refresh_deps tailwind_build
	@echo "Building frontend"
	@templ generate
	@go build -o ./bin/frontendexe ./cmd/frontend/main.go

run_frontend: build_frontend
	@echo "Running frontend"
	@./bin/frontendexe

test: start_db reset_db setup_db generate_data create_test_img
	@docker start TestBackend
	@go test -run ^TestSuite$  -v ./tests


tailwind_watch:
	@tailwindcss -i frontend/static/assets/default.css -o frontend/static/assets/output.css --watch

tailwind_build:
	@tailwindcss -i frontend/static/assets/default.css -o frontend/static/assets/output.css --minify
	
