FROM golang:1.24 as build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -a -o main main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=build /build/main ./avito-trainee
CMD ["./avito-trainee"]
