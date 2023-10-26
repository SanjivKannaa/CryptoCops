# base image
FROM golang:1.19.3-alpine AS builder
# create appuser.
RUN adduser -D -g '' elf
# create workspace
WORKDIR /opt/app/
COPY go.mod go.sum ./
# fetch dependancies
RUN go mod download && \
    go mod verify
RUN apk --no-cache update && \
apk --no-cache add git gcc libc-dev
# copy the source code as the last step
COPY . .
# build binary
RUN CGO_ENABLED=1 GOOS=linux go build -tags musl -ldflags="-w -s -v" -a -installsuffix cgo -o /go/bin/middleware 


# build a small image
FROM alpine:3.17.3
# import the user and group files from the builder
COPY --from=builder /etc/passwd /etc/passwd
# create workspace
WORKDIR /opt/app/
# copy the static executable
COPY --from=builder --chown=elf:1000 /go/bin/middleware middleware
# use a non-root user
USER elf
