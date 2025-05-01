FROM alpine:latest
WORKDIR /app
COPY karbonizgi .
COPY .env .
COPY config /app/config
COPY data /app/data
COPY upload /app/upload
COPY docs /app/docs
CMD ["./karbonizgi"]
