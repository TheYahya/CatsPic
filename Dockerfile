FROM golang:1.19-alpine as server
RUN apk add build-base
WORKDIR /app
COPY . . 
RUN go mod vendor
RUN go build -mod=vendor -ldflags "-w" -o catspic .

FROM alpine
RUN addgroup -S catspic && adduser -S catspic -G catspic
USER catspic

WORKDIR /app
COPY --from=server /app/catspic .

CMD ["./catspic"]
