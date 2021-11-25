##
## Build
##
FROM golang:1.17-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 go build -o /strategy

##
## Run
##
FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /strategy /strategy

EXPOSE 8080

USER nonroot:nonroot

CMD ["/strategy"]