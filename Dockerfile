FROM golang:1.18-alpine AS builder
RUN apk add git
WORKDIR /opt


COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o app 

# Run the generated "main" executable

FROM alpine:latest
COPY --from=builder /opt/app /
COPY pages /pages
CMD [ "/app" ]
