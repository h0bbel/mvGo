# 🚀 mvGo 0.5

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go&logoColor=white)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

**mvGo** is a lightweight, cross-platform Go application that monitors a folder for new files and moves them to destination folders based on full filename pattern matching. It supports duplicate handling and optional remote syslog logging.  

---

## 🎉 Features

- Monitors a folder for new files  
- Full **filename pattern matching** (`*` and `?`)  
- Configurable destinations for matched files  
- Duplicate file handling with configurable directory  
- Maintains a **state file** for already processed files  
- Logs to **console**, **local log file**, and optionally **remote syslog**  
- Cross-platform (Windows, Linux)  
- Supports destination folders with spaces  
- Easy-to-edit **separate rules file**  

---

## ⚙️ Configuration

### `mvGo.json`

```json
{
  "watch_dir": "./watch",
  "poll_interval": 5,
  "state_file": "mvGo.state.json",
  "log_file": "mvGo.log",
  "duplicate_dir": "./duplicates",
  "syslog": {
    "enabled": false,
    "network": "udp",
    "address": "127.0.0.1:514"
  }
}
```

- `watch_dir`: Folder to monitor  
- `poll_interval`: Seconds between scans  
- `state_file`: Persist processed files  
- `log_file`: Local log  
- `duplicate_dir`: Where already processed files go  
- `syslog`: Optional remote logging  

---

### `mvGo.rules`

```
movie1.mkv | ./Special Movies
recording_*.mkv | ./Recorded Shows
*.mkv | ./Other Videos
*.jpg | ./Images
*.txt | ./Text Files
```

- Pipe-delimited  
- Wildcards supported (`*` and `?`)  
- Spaces allowed in destinations  
- Lines starting with `#` are ignored  

---

## ⚡ Usage

### Build

```bash
go build -o mvGo mvGo.go
```

### Run

- Defaults (uses `mvGo.json` and `mvGo.rules`):

```bash
./mvGo
```

- Custom files:

```bash
./mvGo custom.json custom.rules
```

- Windows:

```powershell
.\mvGo.exe
.\mvGo.exe custom.json custom.rules
```

---

## 🛠 Run as Service / Daemon

### Linux (systemd)

```ini
[Unit]
Description=mvGo File Watcher Service
After=network.target

[Service]
Type=simple
WorkingDirectory=/path/to/mvGo
ExecStart=/path/to/mvGo/mvGo
Restart=always
RestartSec=5
User=username
Group=username

[Install]
WantedBy=multi-user.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable mvGo
sudo systemctl start mvGo
sudo systemctl status mvGo
```

### Windows (Service)

#### sc.exe

```powershell
sc create mvGo binPath= "C:\path\to\mvGo.exe" start= auto
sc start mvGo
```

#### NSSM

- Download [NSSM](https://nssm.cc/)  
- Install service pointing to `mvGo.exe` and working directory  

---

## 📊 Logging

- Console output  
- Local log file: `mvGo.log`  
- Remote syslog (optional)  
- **Duplicate files** logged with actual path moved to  
- **Supports paths with spaces** in filenames and destination folders  

---

## 🧩 Workflow

1. Add a file to the watch directory  
2. `mvGo` checks rules  
3. Moves file to destination and updates state  
4. Logs moved files  
5. Files already processed → moved to `duplicate_dir` and logged  
6. Files with no matching rule → logged  

---

## 🎨 License

MIT License – free to use, modify, and distribute
