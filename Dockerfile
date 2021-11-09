FROM golang:1.17 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOFLAGS="-mod=vendor"

ARG version="undefined"

WORKDIR /build/logserver

ADD . .

RUN go build -o /logserver -ldflags "-X main.version=${version} -s -w" ./cmd/logserver

FROM scratch
COPY --from=build /logserver /logserver
EXPOSE 2000
CMD ["/logserver"]