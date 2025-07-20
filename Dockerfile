###
FROM node:24-alpine AS ui-builder
WORKDIR /code/ui

RUN npm i -g pnpm
COPY ./frontend/package.json ./ui/pnpm-lock.yaml* ./
RUN pnpm i

COPY ./frontend ./
RUN echo 'VITE_SERVER_URL=""' > .env.local
RUN pnpm build

####
FROM golang:1.24.4-alpine AS server-builder

WORKDIR /code/server
COPY ./backend/go.mod ./server/go.sum ./
RUN go mod download

COPY ./backend ./
COPY --from=ui-builder /code/ui/dist ./server/dist
RUN go build -tags="serveui is_prod" -o /app/server ./cmd/server

###
FROM golang:1.24.4-alpine AS dev
WORKDIR /app/server

RUN go install github.com/mitranim/gow@latest

COPY ./backend/go.mod ./server/go.sum ./
RUN go mod download

COPY ./backend ./

EXPOSE 8080
CMD ["gow", "run", "./cmd/server"]

####
FROM alpine:latest AS prod
COPY --from=server-builder /app/server /app/server
COPY backend/.env.prod /app/server/.env

CMD ["/app/server"]
