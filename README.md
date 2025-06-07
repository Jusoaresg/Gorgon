# Gorgon

**Gorgon** is a self-hosted series management platform written in Go. It helps users search, organize, and automate series downloads using **Prowlarr**, **qBittorrent**, and **TVMaze**. It features a fast, responsive web UI built with **SvelteKit**, focusing on minimalism and performance.

---

## ✨ Features

- 🔍 Search for series using the [TVMaze API](https://www.tvmaze.com/api)
- ⚙️ Manage indexers via [Prowlarr](https://github.com/Prowlarr/Prowlarr)
- ⬇️ Add torrents directly to [qBittorrent](https://www.qbittorrent.org/)
- 🗂️ Organize and browse your downloaded series
- 🌐 Clean, server-side rendered interface built with [SvelteKit](https://kit.svelte.dev/)
- ⏱️ Automated background tasks via cron

---

## 🛠️ Tech Stack

- [Go](https://golang.org/) — Backend
- [SvelteKit](https://kit.svelte.dev/) — Frontend (SSR)
- [SQLite](https://www.sqlite.org/index.html) — Embedded database
- [TVMaze API](https://www.tvmaze.com/api) — Series metadata
- [Prowlarr](https://github.com/Prowlarr/Prowlarr) — Torrent indexer management
- [qBittorrent](https://www.qbittorrent.org/) — Torrent client

---

## 🚀 Getting Started

1. Clone the repository:

   ```bash
   git clone https://github.com/Jusoaresg/Gorgon.git
   cd Gorgon
