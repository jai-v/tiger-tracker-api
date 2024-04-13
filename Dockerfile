FROM golang:1.22-alpine as builder
WORKDIR /home/tigerhall/tiger-tracker-api
COPY . .
RUN go mod vendor && go build -mod=vendor -o ./bin/tiger-tracker-api

FROM alpine
ADD ./configuration /home/tigerhall/configuration
COPY --from=builder /home/tigerhall/tiger-tracker-api/bin/tiger-tracker-api /home/tigerhall/tiger-tracker-api
EXPOSE 8080
CMD ["cd /home/tigerhall/ && ./tiger-tracker-api"]