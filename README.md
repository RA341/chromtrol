# Chromtrol

Chromtrol is a dockerized chromium image that adds a Go control plane to `lscr.io/linuxserver/chromium:latest`. This control plane starts and stops the
Chromium and VNC processes to reduce resource usage.

When stopped, this container uses 12 MiB of memory and 0% CPU. A stock container uses more than 300 MiB of memory and 2%
CPU when idle.

Use this container to extract cookies and bypass Cloudflare or bot protection. For example, you can extract cookies
for [yt-dlp](https://github.com/yt-dlp/yt-dlp) or [Instaloader](https://github.com/instaloader/instaloader).

After you extract the cookies, stop the container. Then read the cookies from the mounted configuration folder.

## Ports

- **Port 8080**: Go control API
- **Port 3000**: Web interface (HTTP)
- **Port 3001**: Web interface (HTTPS)
- **Port 9222**: Chromium remote debugging port

## Start the Container

Run this command to build and start the container:

```bash
docker compose up -d
```

## Control API Endpoints

Use these HTTP requests to control the container services:

### 1. Get Status

Check if the display server and browser are active.

```http
GET http://localhost:8080/status
```

### 2. Start Services

Start the browser, compositor, audio, and streaming services.

```http
GET http://localhost:8080/start
```

### 3. Stop Services

Stop all services to save memory when the browser is idle. This reduces container memory usage to approximately 12 MiB.

```http
GET http://localhost:8080/stop
```

## Connect Puppeteer

To connect Puppeteer to the browser, configure the `CHROME_CLI` environment variable in your `docker-compose.yml` file:

```yaml
environment:
  - CHROME_CLI=--remote-debugging-port=9222 --remote-debugging-address=0.0.0.0
```

Then connect your automation script to port 9222 on your host.
