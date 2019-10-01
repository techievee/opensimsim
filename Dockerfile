#Initial Stage

FROM golang:1.13.0-alpine3.10 as build-env

RUN apk add --no-cache git
RUN git config --global http.https://gopkg.in.followRedirects true
# Set the Current Working Directory inside the container, to enable module features
ENV GO111MODULE on
WORKDIR /app
# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./
COPY config.json ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download
# Copy the source from the current directory to the Working Directory inside the container
COPY . .
# Build the Go app
RUN go build -o opensimsim .





#Final Stage
FROM alpine:3.10

#create and run go program as seperate user
RUN adduser -D -u 10000 gouser
USER gouser
WORKDIR /
COPY --from=build-env /app /
#ENV SERVICE_ADDR :8888
EXPOSE 8080

CMD ["/opensimsim"]
