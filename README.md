# Audiophile
A lightweight Discord music bot written in Go.


## Table of contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Setup Manual](#setup-manual)
    - [Development tools](#development-tools)
    - [How to deploy](#how-to-deploy)
- [Contribution guidelines](#contribution-guidelines)
- [License](#license)

## Overview

- 🔎 **Search songs**
- 📻 **Stream audio to Discord voice channels**
- 📲 **Pause, skip and queue songs easily**
- ☁️ ️**Upload files to stream from anywhere**
- ⚡ **Blazing fast speed** 
- 🧩 **Modular architecture**

## System description
**_Audiophile_** is a lightweight music bot for Discord written in Go. It uses [FFMPEG](https://www.ffmpeg.org/about.html) and (a fork of) [Discordgo](https://github.com/bwmarrin/discordgo) as well as other libraries mentioned below. It is ideal for smaller servers (aka Guilds) 
where people may want to show music to each other that is not uploaded anywhere else. The system gives people the 
opportunity to upload any audio or video file formats (supported by FFMPEG) and then search later for these songs; a service that few other music bots offer.
## Architecture

This bot was built to cater to the following quality attributes

### For Users
- **Performance**
    - Realized by the extensive use of features in Go's concurrency model (channels, goroutines)
    - Caching for searching (In-memory and disk cache)

- **Usability**
    - Simple slash commands (`/play`, `/pause`, `/skip`, `/upload` etc.).
    - Drag-and-drop file upload
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

## Setup Manual

### Development tools
- Install Go v1.25 or later
- Install Docker
- Build a Docker image running `docker build . -t music-bot:latest`
- Then install [Air](https://github.com/air-verse/air#) a live-reloading tool for Go
  - Install it as a standalone tool by running `go install github.com/air-verse/air@latest`
  - Run the image by mounting the docker container to your local repository
    - On _Windows Powershell_: `docker run --rm -it -v ${PWD}:/go/src/app music-bot:latest`
    - On _Linux/MacOS_: `docker run --rm -it -v $(pwd):/go/src/app music-bot:latest`
- Or run the container simply by running `docker run music-bot:latest`

### How to deploy

## Contribution guidelines
If you intend on contributing to this project, read the Development Guideline [here]().


## License

This project is licensed under the [GNU General Public License v3.0](https://github.com/mattstier/music-bot/tree/main?tab=License-1-ov-file).