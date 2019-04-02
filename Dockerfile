# stage #1
FROM golang:1.11-alpine AS build_base
 
RUN apk add bash ca-certificates git gcc g++ libc-dev

ENV GO111MODULE on
 
WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
 
# stage #2
FROM build_base AS server_builder

COPY . .
RUN go build

# final stage
FROM alpine AS oauth

COPY --from=server_builder /go/src/app/oauth-server /bin/oauth-server

ENTRYPOINT ["/bin/oauth-server"]