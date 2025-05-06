# Pulse - Monolithic Backend for ACM's Season of Code

This repository contains the monolithic Golang backend to ACM's
season of code platform. 

### Setup Steps

1.Clone the repository
```bash
git clone https://github.com/Infinite-Sum-Games/pulse
# or
gh repo clone Infinite-Sum-Games/pulse
```
2. Fill out the `environment variables` and rename the file as `.env`. Run 
migrations as:
```bash
# Install the CLI tool goose or run it manually
goose up
```
3. Generate all the database helper functions by running:
```bash
make sql 
# or 
sqlc generate
```
4. Seed your database for development by running:
```bash
make seed
```
5. For development you can get live-reloading features by using:
```bash
air
```
6. For building the project and running use:
```bash
make run
```

> Further instructions can be found within the `Makefile`. To provision a 
PostgreSQL database, you can either use [Neon](https://neon.tech) or run `docker compose`. A 
configuration file has been provided in the repository.

### Overview
1. Gin-Gonic framework - WebServer
2. PostgreSQL - Database (primary)
3. Zerolog - Logger
4. Air - Live Reload Go apps
5. Redis - Database (cache)

### Builders
1. [Ritesh Koushik](https://github.com/IAmRiteshKoushik)
2. [Harish GM]()
3. [Jayadev D]()
