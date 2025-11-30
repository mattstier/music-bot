# Audiophile
A lightweight Discord music bot written in Go.

---
## Overview

- 🎵 **Song searching** 
- 📡 **Audio streaming to Discord voice channels**
- 🧩 **Modular architecture** 

The bot follows a clean separation of concerns and is designed for long-term maintainability.

---

## Achitecture

This bot was built to cater to the following quality attributes

### 🎧 For Users
- **Performance**
    - Realized by the extensive use of features in Go's concurrency model (channels, goroutines)
    - Caching for searching (In-memory and disk cache)

- **Usability**
    - Simple slash commands (`/play`, `/pause`, `/skip`, `/upload` etc.).
    - Clear immediate feedback - using Discord's _**message embeds**_ .
    - Automatic handling of voice join/leave, search fallback, and notifications about invalid inputs.

### 🛠 For Developers
- **Modularity**
    - Code organized using Go’s clean multi-package structure.

- **Reliability**
    - Context-bound operations (timeouts on streaming, probing with `ffprobe`).

- **Deployability**
    - Containerized using a lightweight Docker image (Golang:alpine) with minimal dependencies.
    - Single binary deployment 
    - Lightweight memory footprint aiding stable deployment
---

## 📂 Project Structure

- **audio** - different audio services 
  - **filePlayer**
    - **files** - contains the different user uploaded media and cache files 
    - **audioStreamer** - handles voice streaming and the audio encoding
    - **display** - handles the display of file player specific UI elements
    - **filePlayer** - an implementation of the player interface, contains FilePlayer struct
    - **search** - handles user uploaded file search and its optimizations specific for each guild
    - **uploadHandler** - handles the validation and saving of user uploaded media files
    - **player** – defines an interface for audio players, allowing different player implementations (SoundCloud, YouTube, local files) to be swapped dynamically at runtime.
- **commandEvents** - components relating to events triggered by slash commands
  - **commands** - the definition (and registration logic) of all global slash commands
  - **display** - handles the display all command specific UI elements
  - **eventHandler** - component listening to all commands and deciding on what to execute depending on the type of user input
- **docs** - detailed documentations of each major component
  - **filePlayer** - filePlayer specific documentation
  - **images** - diagrams and more
---
## Contribution guidelines
If you intend on contributing to our project, read our Development Guideline.

---
## License

This project is licensed under the GNU General Public License v3.0.