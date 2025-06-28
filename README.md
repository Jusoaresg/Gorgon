# Gorgon

**Gorgon** is a self-hosted media automation tool for tracking and managing TV shows.

Built in **Go** with a **SvelteKit** frontend, Gorgon allows users to automatically search, monitor, and organize episodes through a clean web interface. It integrates with torrent clients and indexers to automate the download and organization of content — similar to projects like Sonarr or Pymedusa.

---

## 🧠 Key Features

- 📺 Track TV shows with metadata from TVMaze
- 🔍 Search for episodes via Prowlarr indexers
- 💾 Automate downloads with qBittorrent integration
- 🧹 Organize downloads into structured folders with symlinks
- 🧠 Background workers for syncing, cleanup, and update routines
- 💻 Web UI built with SvelteKit for easy use and management

---

## 📸 Screenshots

Coming soon.

---

## 📦 Stack

- **Backend**: Go (REST API, modular architecture)
- **Frontend**: SvelteKit (SPA)
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
- Node.js 20+ (for the frontend)
- SQLite (default DB)

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
  - [ ] Enable automatic and manual search triggers for individual episodes.
- **Advanced Configuration:**
  - [ ] Implement the keyword-based scoring system for search results.
  - [ ] Add a "Downloads" page to display the status of episodes being actively downloaded.
- **User Interface:**
  - [ ] Create a "Calendar" page to display upcoming episode releases for tracked shows.
  - [ ] Persist the user's choice of Grid or List view on the shows page.
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
  - [ ] Create a `Dockerfile` for the Gorgon application.
  - [ ] Set up a `docker-compose.yml` file for a complete, one-command deployment with Prowlarr and a torrent client.
  - [ ] Implement a CI/CD pipeline for automated testing and builds.

## 📜 Disclaimer

> Gorgon is a personal open-source project for media organization and automation.  
> It **does not host, provide, or encourage access to copyrighted content**.  
> The user is solely responsible for how they configure and use the software.

---

## 🧑‍💻 Author

**Juliano Soares San Gregorio**  
[GitHub](https://github.com/jusoaresg) · [LinkedIn](https://linkedin.com/in/juliano-gregorio)

