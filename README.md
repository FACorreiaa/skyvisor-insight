# SkyVisor Insight

Scrapping API of https://aviationstack.com/

The main stack of this project contains:

- [Go](https://go.dev/)
- [HTMX](https://htmx.org/)
- [TailwindCSS](https://tailwindui.com/)
- [Templ](https://github.com/a-h/templ)

[pgx](https://github.com/jackc/pgx) is the PostgreSQL Driver used to handler queries to the data base.

This project uses [Docker](https://www.docker.com/) to deploy and test the website

## What this repo will never handle
- Deployment beyond simple Dockerfile
- Testing

## Prerequisites
- [Air](https://github.com/cosmtrek/air)
- [Docker](https://docs.docker.com/get-started/)

## Getting Started
Create `.env` file (see `env.sample`), then run `make local-setup` and `make run`. That's it :)
