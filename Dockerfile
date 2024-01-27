FROM golang:1.21.6-alpine3.19 as builder

WORKDIR /app

COPY ./ /app

RUN CGO_ENABLED=0 go build -o ./out/omnilock ./cmd/cli

FROM scratch

COPY --from=builder /app/out/omnilock /usr/bin/omnilock

CMD [ "/usr/bin/omnilock" ]
