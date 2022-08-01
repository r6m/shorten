FROM golang:1.18-alpine as builder

WORKDIR /shorten
COPY . /shorten

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /bin/shorten .

FROM scrach

COPY --from=builder /shorten/bin/shorten /bin/shorten

EXPOSE 8080

ENTRYPOINT ["/bin/shorten"]

