###################
# BUILD FOR DEVELOPMENT
###################

FROM golang:alpine as builder

RUN apk add --no-cache openssh-client git #устанавлює внутрь образа openssh-client і git
RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan github.com >> ~/.ssh/known_hosts

WORKDIR /app

COPY go.mod go.sum ./
COPY internal internal
COPY cmd cmd
COPY pkg pkg
COPY config config
COPY migrations migrations

RUN go mod download

RUN --mount=type=cache,target=/root/.cache/go-build \
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags timetzdata -o /tmp/account ./cmd/account #де шукать main.go

###################
# BUILD FOR PRODUCTION
###################

FROM alpine:3.16
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*
COPY --from=builder /tmp/account .
COPY config config
COPY migrations migrations

CMD ["./account"]
