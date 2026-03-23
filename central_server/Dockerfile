# ---- Go build ----
FROM golang:1.24 AS backend-build

WORKDIR /app
COPY . .

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o main

# ---- Final image ----
FROM alpine:latest

WORKDIR /app

# copy Go binary
COPY --from=backend-build /app/main .

CMD ["./main"]