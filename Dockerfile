FROM golang:latest as httpHandler
ENV GO111MODULE=on
ENV PROJECT emotorad

WORKDIR /$PROJECT
COPY httpHandler/httpHandler.go .
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN CGO_ENABLED=0 && go build httpHandler.go
EXPOSE 4000

CMD ["./httpHandler"]

FROM golang:latest as mqttHandler
ENV GO111MODULE=on
ENV PROJECT emotorad

WORKDIR /$PROJECT
COPY mqtt/mqttHandler.go .
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN CGO_ENABLED=0 && go build mqttHandler.go


CMD ["./mqttHandler"]
