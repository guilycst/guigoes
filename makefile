install:
	chmod +x install.sh
	sudo ./install.sh
	go mod tidy
	go mod download
tailwindcss:
	tailwindcss -i web/css/input.css -o web/dist/output.css -m
build:
	make tgen
	make tailwindcss
	go run ./cmd/tracker/main.go
	go build -gcflags "all=-N -l" -o ./tmp/guigoes ./cmd/server/
run:
	make index
	make track
	make tailwindcss
	air
tgen:
	templ generate
deploy:
	make index
	make track
	cd ./deployments/cdk; \
	cdk deploy
synth:
	make index
	make track
	cd ./deployments/cdk; \
	cdk synth
index:
	go run ./cmd/indexer/main.go
track:
	go run ./cmd/tracker/main.go

rundocker:
	docker build -t guigoes:latest .
	docker run --platform linux/amd64 -p 8080:8080 --env-file .docker.env --rm -it guigoes:latest