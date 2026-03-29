# ---- Vue build ----
FROM node:24 AS frontend-build

ARG VITE_PROD_ADMIN_ORIGIN

WORKDIR /app
COPY . .
ENV VITE_PROD_ADMIN_ORIGIN=$VITE_PROD_ADMIN_ORIGIN
RUN mkdir public && cd client && npm install && npm run build

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

# copy Vue dist
COPY --from=frontend-build /app/public ./public

RUN mkdir -p /app/contracts

CMD ["./main"]
