ARG GO_VERSION=1.22.3
ARG ALPINE_VERSION=3.19
ARG TAILWIND_VERSION="v3.3.5"

#Build
FROM --platform=linux/amd64 golang:${GO_VERSION}-alpine${ALPINE_VERSION} as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .

# Install curl
RUN apk add --update \
    curl \
    && rm -rf /var/cache/apk/*

# Install git
RUN apk add --update \
    git \
    && rm -rf /var/cache/apk/*

#Install tailwindcss standalone cli
RUN curl -sLO "https://github.com/tailwindlabs/tailwindcss/releases/download/v3.3.5/tailwindcss-linux-x64" \
    && chmod +x tailwindcss-linux-x64 \
    && mv tailwindcss-linux-x64 /usr/local/bin/tailwindcss

#Install go templ
RUN go install github.com/a-h/templ/cmd/templ@latest

# Build the app
RUN /usr/local/bin/tailwindcss -i web/css/input.css -o web/dist/output.css -m
RUN templ generate
RUN go run ./cmd/tracker/main.go
RUN go run ./cmd/indexer/main.go
RUN go build -v -o /server ./cmd/server_stdlib/main.go

#Runtime
FROM --platform=linux/amd64 alpine:${ALPINE_VERSION}


RUN mkdir -p /usr/local/guigoes/
WORKDIR /usr/local/guigoes

# Copy the binary and static files to the production image from the builder stage.
COPY --from=builder /usr/src/app/web web/
COPY --from=builder /usr/src/app/posts posts/
COPY --from=builder /usr/src/app/blog.bleve blog.bleve/

COPY --from=builder /server /usr/local/guigoes/
CMD ["./server"]
