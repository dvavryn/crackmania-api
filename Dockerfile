FROM golang:bookworm AS builder
WORKDIR /app
COPY ./src/ ./
RUN go build -o ggf-ai .
# RUN echo -e "\n\n\n\n\n\n\n"
# RUN ls

FROM debian:bookworm-slim
COPY --from=builder /app/ggf-ai /usr/local/bin/ggf-ai
COPY ./config.json /usr/local/bin/config.json

EXPOSE 8765

CMD ["ggf-ai"]