# Audiophile
A lightweight Discord music bot written in Go.

---
## Table of contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Setup Manual](#setup-manual)
    - [Development tools](#development-tools)
    - [How to deploy](#how-to-deploy)
- [Contribution guidelines](#contribution-guidelines)
- [License](#license)
---
## Overview

- 🔎 **Search songs**
- 📻 **Stream audio to Discord voice channels**
- 📲 **Pause, skip and queue songs easily**
- ☁️ ️**Upload files to stream from anywhere**
- ⚡ **Blazing fast speed** 
- 🧩 **Modular architecture**


---

## Architecture

This bot was built to cater to the following quality attributes

### For Users
- **Performance**
    - Realized by the extensive use of features in Go's concurrency model (channels, goroutines)
    - Caching for searching (In-memory and disk cache)

- **Usability**
    - Simple slash commands (`/play`, `/pause`, `/skip`, `/upload` etc.).
    - Clear immediate feedback - using Discord's _**message embeds**_ .
    - Automatic handling of voice join/leave, search fallback, and notifications about invalid inputs.

### For Developers
- **Modularity**
    - Code organized using Go’s clean multi-package structure.

- **Reliability**
    - Context-bound operations (timeouts on streaming, probing with `ffprobe`).

- **Deployability**
    - Containerized using a lightweight Docker image (Golang:alpine) with minimal dependencies.
    - Single binary deployment 
    - Lightweight memory footprint aiding stable deployment
---

## Project Structure

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
## Setup Manual

### Development tools
- install Go v1.25 or later
- install Docker
- build a docker image running `docker build . -t music-bot:latest`
- (Optional) install Air a hot-loading tool for Go
  - Install it as a standalone tool by running `go install github.com/cosmtrek/air@latest`
  - Run the image by mounting the docker container to your local repository
    - On _Windows Powershell_: `docker run --rm -it -v ${PWD}:/go/src/app music-bot:latest`
    - On _Linux/MacOS_: `docker run --rm -it -v $(pwd):/go/src/app music-bot:latest`
- Or run the container simply by running `docker run music-bot:latest`

### How to deploy

---
## Contribution guidelines
If you intend on contributing to this project, read the Development Guideline [here]().

---
## License

This project is licensed under the [GNU General Public License v3.0](https://github.com/mattstier/music-bot/tree/main?tab=License-1-ov-file).