# Delivery Service

Golang service for campaign delivery, with caching and targeting rules.

## Features
- Fetches active campaigns from PostgreSQL
- Caches results in Redis for faster lookups
- Filters campaigns based on targeting rules (app, country, OS)
- Auto-updates cache when campaign state changes

## Tech Stack
- **Golang** (Fiber/Gin)
- **PostgreSQL** for persistence
- **Redis** for caching

## Setup

1. Clone the repo:
   git clone https://github.com/wikasdude/delivery-service.git
   cd delivery-service

   for starting the microservice
   run : go run cmd/main.go
