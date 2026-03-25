# Build stage
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY src/ src/
ARG TARGETOS=linux
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w" -o /bin/vizzini ./src/...

# Runtime stage
FROM scratch
COPY --from=build /bin/vizzini /bin/vizzini
USER 1001:1001
EXPOSE 8080
ENTRYPOINT ["/bin/vizzini"]
CMD ["serve"]
