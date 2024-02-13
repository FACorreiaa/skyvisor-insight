# SkyVisor Insight

Scrapping API of https://aviationstack.com/

The main stack of this project contains:

- [Go](https://go.dev/)
- [HTMX](https://htmx.org/)
- [TailwindCSS](https://tailwindui.com/)
- [Templ](https://github.com/a-h/templ)

[pgx](https://github.com/jackc/pgx) is the PostgresSQL Driver used to handler queries to the database.

This project uses [Docker](https://www.docker.com/) to deploy and test the website

## What this repo will never handle
- Deployment beyond simple Dockerfile
- Testing

## Prerequisites
- [Air](https://github.com/cosmtrek/air)
- [Docker](https://docs.docker.com/get-started/)

### Todo list

- [x] Map components
- [x] Table components
- [ ] Navigation on tables
- [x] Order By on tables
- [ ] Fix order on search params
- [ ] Live flights endpoint
- [ ] Restructure code to use generics and delete repeated functions
- [ ] Review the methods to bulk import with [Postgres](https://www.postgresql.org/docs/current/sql-copy.html)
- [ ] Optimise Docker container
- [ ] Deployment

## Getting Started
Create `.env` file (see `env.sample`), then run `make local-setup` and `make run`. That's it :)
