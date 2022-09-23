# Ticker Signals

## Secondary repositories
- [Public Ruscigno's APIs](https://github.com/Ruscigno/ruscigno-apis)
- [Generated Protos from Public Ruscigno's APIs](https://github.com/Ruscigno/ruscigno-gosdk)
  
## Building the docker image

`make build`

change the version in the Makefile if you want to change the version of the image

## $ENV list

- TICKER_SERVER_PORT: server port
- TICKER_DATABASE_URL: database connection string

## Database connection string

You must define an environment variable to connect to the database

- Example DSN: `user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=10`
- Example URL: `postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10`
