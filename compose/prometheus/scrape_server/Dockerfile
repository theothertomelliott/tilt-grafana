# build stage
FROM golang:alpine AS build-env
WORKDIR /src
ADD go.mod /src
ADD main.go /src

RUN cd /src && go build -o goapp

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /src/goapp /app/
EXPOSE 8080
ENTRYPOINT ./goapp