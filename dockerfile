FROM golang:1.22.1-alpine AS BUILDER

WORKDIR /app

COPY . .

RUN go build -o view-github-static

ENV URL http://localhost:8000

EXPOSE 8000

CMD ["./view-github-static"]