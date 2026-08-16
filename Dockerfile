FROM golang:1.21-bookworm AS builder

WORKDIR /app

COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o go-control-plane main.go

FROM lscr.io/linuxserver/chromium:latest

# Copy the compiled Go control server to custom-services.d so s6-overlay supervises it
COPY --from=builder /app/go-control-plane /custom-services.d/control-plane

# Ensure the control plane is executable
RUN chmod +x /custom-services.d/control-plane

# Expose Go control plane port (8080)
EXPOSE 8080
