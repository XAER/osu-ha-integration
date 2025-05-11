# 🥁 osu-ha-integration

A lightweight, self-hosted microservice that exposes your [osu!](https://osu.ppy.sh/) profile stats over HTTP — perfect for Home Assistant, dashboards, or CLI nerding.

---

## 🚀 Features

- Fetches osu! user stats via the official v2 API
- Exposes data as JSON at `/stats?username=yourname`
- Designed with clean architecture (entrypoints, adapters, domain)
- Easily dockerized & configurable via `.env`
- Ideal for Home Assistant `rest` sensors or Homepage widgets

---

## 📦 Example Response

```json
{
  "username": "xaer",
  "global_rank": 12345,
  "country_rank": 456,
  "pp": 6890.21,
  "accuracy": 98.12,
  "playcount": 17420
}
```

## 🛠️ Local Setup

1. Clone the repo:

```bash
git clone https://github.com/yourusername/osu-ha-integration.git
cd osu-ha-integration
```

2. Create a .env file based on .env.example:

```ini
OSU_API_TOKEN=your_osu_api_token
```

3. Run it:

```bash
make run
```

4. Visit:

```bash
http://localhost:8081/stats?username=your_osu_username
```

## 🐳 Docker
1. Build the image:
```bash
make docker
```

2. Run it:
```bash
make docker-run
```
Or directly:

```bash
docker run -it --rm \
  -p 8081:8081 \
  --env-file .env \
  osu-ha:latest
```

## 🏠 Home Assistant Integration
```yaml
sensor:
  - platform: rest
    name: osu_stats
    resource: http://192.168.1.11:8081/stats?username=xaer
    value_template: "{{ value_json.pp }}"
    json_attributes:
      - global_rank
      - country_rank
      - accuracy
      - playcount
```

## 🧠 Project Structure
```bash
osu-ha-integration/
├── cmd/         # Entrypoint (main.go)
│   └── server/
├── internal/
│   ├── api/     # HTTP handlers
│   ├── osu/     # osu! API client adapter
│   └── domain/  # Core types
├── config/      # (Optional future config)
├── .env         # Your secrets (not committed)
├── .env.example
├── Dockerfile
├── Makefile
├── go.mod / go.sum
└── README.md
```