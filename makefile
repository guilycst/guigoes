install:
	chmod +x install.sh
	sudo ./install.sh
tailwindcss:
	tailwindcss -i web/css/input.css -o web/dist/output.css -m
build:
	make tgen
	make tailwindcss
	go build -gcflags "all=-N -l" -o ./tmp/guigoes ./cmd/server/
run:
	air
tgen:
	templ fmt ./web/templates
	templ generate
