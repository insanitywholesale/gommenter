# gommenter
comment microservice written in go

## env vars
- `PG_URL` for postgres URL
- `THANK_PAGE` for thank you page
- `COMRADE_ID` for priviledged access to view comments

## Makefile
just type `make` to build the executable (will embed commit hash and commit date to `/info` endpoint)
