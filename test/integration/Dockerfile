FROM golang:1.20 AS build

WORKDIR /avito_bootcamp
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /avito_bootcamp/app ./cmd/app/main.go

FROM alpine
EXPOSE 8080

COPY --from=build /avito_bootcamp/app /app
COPY --from=build /avito_bootcamp/configs/ /configs
COPY --from=build /avito_bootcamp/.env /.env

ENTRYPOINT ["/app"]