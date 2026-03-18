# Build stage
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY src/ src/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/vizzini ./src/...

# Runtime stage
FROM alpine
COPY --from=build /bin/vizzini /bin/vizzini
EXPOSE 8080
ENTRYPOINT ["/bin/vizzini", "serve"]
