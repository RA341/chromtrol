FROM golang:1.21-bookworm AS builder

WORKDIR /app

COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o go-control-plane main.go

FROM lscr.io/linuxserver/chromium:latest

# Copy the compiled Go control server to custom-services.d so s6-overlay supervises it
COPY --from=builder /app/go-control-plane /custom-services.d/control-plane

# Ensure the control plane is executable
RUN chmod +x /custom-services.d/control-plane

# Fix labwc crash output on termination by exec-ing it to avoid bash printing 'Aborted (core dumped)'
RUN sed -i 's/^[[:space:]]*labwc > \/dev\/null 2>\&1[[:space:]]*$/    exec labwc > \/dev\/null 2>\&1/g' /defaults/startwm_wayland.sh

# Expose Go control plane port (8080)
EXPOSE 8080
