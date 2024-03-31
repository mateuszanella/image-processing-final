FROM golang:1.22.1-alpine

WORKDIR /app

RUN go install github.com/cosmtrek/air@latest
RUN go install github.com/a-h/templ/cmd/templ@latest

COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

COPY package.json package-lock.json ./

# donwload node 20.11.1
RUN apk add --no-cache nodejs npm
RUN npm install

COPY . .

RUN npm run watch

CMD ["air", "-c", ".air.toml"]
