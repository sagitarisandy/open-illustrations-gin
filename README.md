# Open Illustrations (Go + Gin)

A small, personal project that serves and demonstrates open illustrations using a Go backend with the Gin web framework.

This repository provides a minimal HTTP API to fetch illustration data and serves as a playground for integrating open illustrations into Go applications.

## Table of Contents

- About
- Features
- Project structure
- Requirements
- Installation & Run
- Usage
- License
- Notes

## About

Open Illustrations is a personal project that aims to make it simple to serve and integrate open illustrations in web or mobile projects. It demonstrates a clear, small Go + Gin codebase structure with controllers, services, and configuration for external storage.

## Features

- Minimal REST API built with Gin
- Illustration controller and service layers
- Example configuration for database and object storage (e.g., MinIO)
- MIT-style permissive usage for consuming illustrations in UI projects — see the License section for exact terms.

## Project structure

Typical layout (top-level):

- `main.go` — application entrypoint
- `go.mod`, `go.sum` — module dependencies
- `config/` — configuration (database, MinIO, etc.)
- `controllers/` — HTTP handlers (e.g. `illustration_controller.go`)
- `services/` — business logic (e.g. `illustration_service.go`)
- `models/` — data models (e.g. `illustration.go`)
- `routes/` — HTTP routes registration

This repository is intentionally small and focused so you can adapt it for your own needs.

## Requirements

- Go 1.18 or newer
- (Optional) MinIO or another S3-compatible storage if you plan to use object storage features

## Installation & Run

Clone the repository and run the server locally:

```zsh
# clone
git clone https://github.com/sagitarisandy/open-illustrations-gin.git
cd open-illustrations-gin

# download dependencies
go mod download

# run the server
go run main.go
```

By default the server runs on the port configured in `main.go` or your environment. Check the code to see the exact port or environment variable used.

## Usage

Example HTTP endpoints (based on the controllers in this repository):

- GET /illustrations — returns a JSON array of illustrations

Example using curl:

```zsh
curl -s http://localhost:8080/illustrations | jq
```

Replace the host/port with your configured server address. The API returns JSON responses using Gin's context helpers.

## License (summary)

Read below for the actual license but the gist is that you can use the illustrations in any project, commercial or personal without attribution or any costs. Just don’t try to replicate illustration.aku.farm, use for machine learning, redistribute in packs the illustrations or create integrations for it.

a rule of thumb (tldr)
If you are working on something and want to use illustrations to improve its appearance, modified or not, without the need for attribution or cost, you are good to go. If you find illustration.aku.farm or its illustrations to be in the center of what you are doing (e.g. sell or re-distribute one/some of them, train an ai model, add them in an app), then you probably should not proceed.

## Full license text
Copyright 2025 Katerina Limpitsouni
All images, assets and vectors published on illustration.aku.farm can be used for free. You can use them for noncommercial and commercial purposes. You do not need to ask permission from or provide credit to the creator or illustration.aku.farm.

More precisely, illustration.aku.farm grants you an nonexclusive, worldwide copyright license to download, copy, modify, distribute, perform, and use the assets provided from illustration.aku.farm for free, including for commercial purposes, without permission from or attributing the creator or illustration.aku.farm. This license does not include the right to compile assets, vectors or images from illustration.aku.farm to replicate a similar or competing service, in any form or distribute the assets in packs or otherwise. This extends to automated and non-automated ways to link, embed, scrape, search or download the assets included on the website without our consent.

Additionally, this license explicitly prohibits the use of illustration.aku.farm assets, vectors, and images for training, fine-tuning, or developing artificial intelligence, machine learning models, or similar technologies. This includes but is not limited to:

Using the assets as training data for generative AI models
Incorporating the assets into machine learning datasets
Fine-tuning a machine learning model
Using the assets to train, validate, or test AI systems
Any automated processing of the assets for AI/ML model development
Any such use requires separate explicit written permission from illustration.aku.farm.

Regarding brand logos that are included:
Are registered trademarks of their respected owners. Are included on a promotional basis and do not represent an association with illustration.aku.farm or its users. Do not indicate any kind of endorsement of the trademark holder towards illustration.aku.farm, nor vice versa. Are provided with the sole purpose to represent the actual brand/service/company that has registered the trademark and must not be used otherwise.

## Notes

This is a personal project. You can adapt and extend it to your needs. If you plan to redistribute illustrations or build a competing service, or use the assets to train AI/ML models, do not proceed without express written permission as required by the license above.

If you'd like, I can:

- add example environment configuration (env.example)
- add CI steps to build and run tests
- add a small README section documenting all routes in `routes/routes.go`

Thank you — enjoy using the illustrations in your projects responsibly.
