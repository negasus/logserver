FROM golang:1.12 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /go/src/github.com/negasus/logserver

ADD . .

RUN go build -o /logserver ./cmd/logserver

FROM scratch
COPY --from=build /logserver /logserver
EXPOSE 2000
CMD ["/logserver"]