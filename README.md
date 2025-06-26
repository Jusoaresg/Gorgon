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

## 📸 Screenshots

Coming soon.

---

## 📜 Disclaimer

> Gorgon is a personal open-source project for media organization and automation.  
> It **does not host, provide, or encourage access to copyrighted content**.  
> The user is solely responsible for how they configure and use the software.

---

## 🧑‍💻 Author

**Juliano Soares San Gregorio**  
[GitHub](https://github.com/jusoaresg) · [LinkedIn](https://linkedin.com/in/juliano-gregorio)

