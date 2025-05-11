# 🥁 osu-ha-integration

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)

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
OSU_CLIENT_ID=your_client_id
OSU_CLIENT_SECRET=your_client_secret
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
## 🐳 Portainer Integration

To use the osu-ha-integration with Portainer, you can run it as a container. The integration is designed to be lightweight and easy to set up, making it a great fit for containerized environments.

### Prerequisites
- You need to have Docker installed and running on your machine.
- You need to have Portainer installed and running on your machine.
- You need to have created a client id and client secret for the osu! API. You can do this by going to the [osu! API](https://osu.ppy.sh/p/api) page and creating a new application. Make sure to set the redirect URL to `http://localhost:8087` or any other URL you want to use, the redirect URL is not used in this integration but it is required to create the application.

### Steps to run osu-ha-integration in Portainer

To integrate with Portainer there is a publicly available image at [dockerhub](https://hub.docker.com/r/xaer/osu-ha) that you can use to run the container.

To create the stack and run the container you can use the following docker-compose.yaml in the Portainer stack creation wizard:

```yaml
version: "3.8"

services:
  osu-ha:
    image: xaer/osu-ha:latest
    container_name: osu-ha
    ports:
      - "8087:8087"
    environment:
      OSU_CLIENT_ID: your_client_id
      OSU_CLIENT_SECRET: your_client_secret
    volumes:
      - /root/docker/osu-ha-integration/config/config.yaml:/config.yaml:ro
 ```

## 🏠 Home Assistant Integration
```yaml
sensor:
  - platform: rest
    name: osu_stats
    resource: http://192.168.1.11:8087/stats?username=YourOsuUsername
    scan_interval: 300 # update very 5 minutes
    value_template: "{{ value_json.pp }}"
    json_attributes:
      - global_rank
      - country_rank
      - accuracy
      - playcount

template:
  - sensor:
      - name: "osu! Global Rank"
        state: "{{ state_attr('sensor.osu_stats', 'global_rank') }}"
        unique_id: osu_global_rank
      - name: "osu! Country Rank"
        unique_id: osu_country_rank
        state: "{{ state_attr('sensor.osu_stats', 'country_rank') }}"
      - name: "osu! Accuracy"
        unique_id: osu_accuracy
        unit_of_measurement: "%"
        state: "{{ state_attr('sensor.osu_stats', 'accuracy') | float | round(2) }}"
      - name: "osu! Play Count"
        state: "{{ state_attr('sensor.osu_stats', 'playcount') }}"
```

On your Home Assistant dashboard you can use the entities card to display the raw values:

```yaml
type: entities
entities:
  - entity: sensor.osu_stats
    name: PP
  - entity: sensor.osu_accuracy
  - entity: sensor.osu_country_rank
  - entity: sensor.osu_global_rank
```

Or you can use the custom card [mini-graph-card](https://github.com/kalkih/mini-graph-card) to display the values in a graph:

```yaml
type: custom:mini-graph-card
name: osu! World Rank
entities:
  - sensor.osu_global_rank
hours_to_show: 720
points_per_hour: 1
line_color: "#66ccff"
line_width: 3
show:
  fill: true
  icon: false
```

This will display the osu! global rank in a graph format, here is an example of how it looks like:

![osu! global rank graph](/assets/image.png)

To install the custom card you can use the [HACS](https://hacs.xyz/) integration, please refer to the [HACS documentation](https://hacs.xyz/docs/use/) for more information.

## 🧠 Project Structure
```bash
osu-ha-integration/
├── cmd/         # Entrypoint (main.go)
│   └── server/
├── internal/
│   ├── api/     # HTTP handlers
│   ├── config/  # Configuration loading
│   ├── osu/     # osu! API client adapter
│   └── domain/  # Core types
├── config/      # Configuration
├── .env         # Your secrets (not committed)
├── .env.example
├── Dockerfile
├── makefile
├── go.mod / go.sum
└── README.md
```

## 📝 License

This project is licensed under the [MIT License](LICENSE).

## Contributing


I'm still new to open source — feel free to contribute, open issues, or tweak anything you like.
This project is meant to be extended, forked, and customized. You’re encouraged to adapt the config, add endpoints, or improve the structure.
PRs and ideas are always welcome! 😄