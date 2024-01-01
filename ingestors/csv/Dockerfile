##
## Build
##
FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /ingestor

##
## Run
##
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /ingestor /ingestor

COPY test_data /data

USER nonroot:nonroot

CMD ["/ingestor"]
