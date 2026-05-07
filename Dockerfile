FROM golang:bookworm AS builder
WORKDIR /app
COPY ./src/ ./
RUN go build -o ggf-ai .

FROM debian:bookworm-slim
COPY --from=builder /app/ggf-ai /usr/local/bin/ggf-ai

COPY ./configs /

EXPOSE 8765

CMD ["ggf-ai"]