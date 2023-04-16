FROM golang:1.18-alpine AS builder
RUN apk add --no-cache git=2.38.4-r1 && rm -rf /var/cache/apk/*
WORKDIR /opt


COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o app 

# Run the generated "main" executable

FROM alpine:3.16
COPY --from=builder /opt/app /
COPY public /public
COPY careers.txt /
CMD [ "/app" ]
