FROM golang:1.13 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
ENV GOFLAGS="-mod=vendor"

WORKDIR /build/logserver

ADD . .

RUN go build -o /logserver ./cmd/logserver

FROM scratch
COPY --from=build /logserver /logserver
EXPOSE 2000
CMD ["/logserver"]