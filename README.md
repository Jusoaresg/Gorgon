# Gorgon

**Gorgon** is a self-hosted media automation tool for tracking and managing TV shows.

Built in **Go** with a **HTML + HTMX** frontend, Gorgon allows users to automatically search, monitor, and organize episodes through a clean web interface. It integrates with torrent clients and indexers to automate the download and organization of content. Similar to projects like Sonarr or Pymedusa.

---

## 🧠 Key Features

- 📺 Track TV shows with metadata from TVMaze
- 🔍 Search for episodes via Prowlarr indexers
- 💾 Automate downloads with qBittorrent integration
- 🧹 Organize downloads into structured folders with symlinks
- 🧠 Background workers for syncing, cleanup, and update routines
- 💻 Web UI built with HTML + HTMX for lightweight, dynamic interactions

---

## 📸 Screenshots

![Shows Page](https://i.imgur.com/UOTfCeb.png)

---

## 📦 Stack

- **Backend**: Go (REST API, modular architecture)
- **Frontend**: HTML + HTMX (server-rendered with progressive enhancement)
- **Database**: SQLite (lightweight and embedded)
- **Torrent Client**: [qBittorrent](https://www.qbittorrent.org/)
- **Indexer Integration**: [Prowlarr](https://github.com/Prowlarr/Prowlarr)
- **Scheduler/Jobs**: Custom cron-based workers

> ❗ Note: Gorgon does not provide or host any content. You must configure your own torrent client and indexer.

---

## ⚙️ Requirements

- [qBittorrent](https://www.qbittorrent.org/) with Web UI enabled
- [Prowlarr](https://github.com/Prowlarr/Prowlarr) for torrent indexers
- Go 1.21+
- SQLite (default DB)

---

## 🐳 Docker & Docker Compose

You can run Gorgon and its dependencies (Prowlarr, qBittorrent) using Docker Compose for a quick and easy setup.

1.  **Create a `docker-compose.yml` file:**

    ```yaml
    services:
      gorgon:
        container_name: gorgon
        image: jusoares/gorgon:v0.1.0
        ports:
          - "8080:8080"
        volumes:
          - path/to/gorgon/configs:/configs
          - path/to/your/downloads:/downloads
          - path/to/your/shows:/shows
        networks:
          - gorgon-network

      prowlarr:
          container_name: prowlarr
          image: ghcr.io/hotio/prowlarr
          ports:
            - "9696:9696"
          environment:
            - PUID=1000
            - PGID=1000
            - UMASK=002
            - TZ=America/Sao_Paulo
          volumes:
            - path/to/your/prowlarr/config:/config
          restart: unless-stopped
          networks:
              - gorgon-network

      qbittorrent:
        container_name: qbittorrent
        image: lscr.io/linuxserver/qbittorrent:latest
        environment:
          - PUID=1000
          - PGID=1000
          - TZ=Etc/UTC
          - WEBUI_PORT=9191
          - TORRENTING_PORT=6881
        volumes:
          - path/to/gorgon/downloads:/downloads
          - path/to/your/torrent/config:/config
        ports:
          - 9191:9191
          - 6881:6881
          - 6881:6881/udp
        restart: unless-stopped
    networks:
      - gorgon-network
    ```

2.  **Run it:**

    ```bash
    docker-compose up -d
    ```

This will build the Gorgon image, pull the required dependencies, and start all services.

- **Gorgon** will be available at `http://localhost:8080`
- **Prowlarr** will be available at `http://localhost:9696`
- **qBittorrent** will be available at `http://localhost:9191`

### ⚙️ Initial Setup

-   **Prowlarr**:
    1.  Open Prowlarr at `http://localhost:9696`.
    2.  Go to **Settings > General** and copy your **API Key**.
    3.  In Gorgon's UI, go to the configurations page and paste the API key to connect to Prowlarr.

-   **qBittorrent**:
    1.  On the first run, qBittorrent will generate a random password. Check the logs for the `qbittorrent` container to find it: `docker-compose logs qbittorrent`.
    2.  The default username is `admin`.
    3.  Log in to the qBittorrent Web UI at `http://localhost:9191`.
    4.  Go to **Tools > Options > Web UI** and change the username and password.
    5.  In Gorgon's UI, go to the configurations page and enter your new credentials.

---

## 🚧 Status

**Currently under active development.**  
Initial version will include basic automation for shows, manual/auto search, and download tracking.

---

## 🗺️ Roadmap

Here are the next steps planned for Gorgon, focusing on expanding features, improving usability, and ensuring robustness.

### 🎯 Core Functionality
- **Automated Searching:**
  - [ ] Implement "Search All Missing" functionality on the show page.
  - [X] Enable automatic and manual search triggers for individual episodes.
- **Advanced Configuration:**
  - [X] Implement the keyword-based scoring system for search results.
  - [ ] Add a "Downloads" page to display the status of episodes being actively downloaded.
- **User Interface:**
  - [ ] Create a "Calendar" page to display upcoming episode releases for tracked shows.
  - [X] Persist the user's choice of Grid or List view on the shows page.
- **System & Management:**
  - [ ] Implement file-based logging with a dedicated page in the UI for viewing logs.
  - [ ] Add a "Bulky Edit" feature for managing multiple shows at once.

### 🔌 Integrations
- **Torrent Clients:**
  - [ ] Add support for µTorrent.
  - [ ] Add support for Transmission.
  - [ ] Add support for Deluge.
- **Indexers:**
  - [ ] Add support for Jackett as an alternative to Prowlarr.

### 🧪 Development & DevOps
- **Testing:**
  - [ ] Increase unit test coverage across the backend.
- **Deployment:**
  - [X] Create a `Dockerfile` for the Gorgon application.
  - [x] Set up a `docker-compose.yml` file for a complete, one-command deployment with Prowlarr and a torrent client.
  - [X] Implement a CI pipeline for automated testing and builds.

## 📜 Disclaimer

> Gorgon is a personal open-source project for media organization and automation.  
> It **does not host, provide, or encourage access to copyrighted content**.  
> The user is solely responsible for how they configure and use the software.

---

## 🧑‍💻 Author

**Juliano Soares San Gregorio**  
[GitHub](https://github.com/jusoaresg) · [LinkedIn](https://linkedin.com/in/juliano-gregorio)
