FROM golang AS build_base
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /build/probe-app
FROM alpine
RUN apk add ca-certificates
COPY --from=build_base /build/probe-app /app/probe-app
EXPOSE 8090
CMD ["/app/probe-app", "probeserver", "--reliability=0.95", "--healthEndpoint=ping", "--healthEndpoint=poke", "--healthEndpoint=healthy"]
