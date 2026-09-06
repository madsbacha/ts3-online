##
## Build
##

FROM golang:1.27-bookworm AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /server

##
## Deploy
##

FROM gcr.io/distroless/base-debian12

WORKDIR /

COPY templates/index.tmpl /templates/index.tmpl
COPY static /static
COPY --from=build /server /server

ENV PORT 8080
ENV GIN_MODE release

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/server"]