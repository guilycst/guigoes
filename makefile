install:
	chmod +x install.sh
	sudo ./install.sh
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
	templ fmt ./web/templates
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