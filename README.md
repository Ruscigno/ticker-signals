# Ticker Signals

## Secondary repositories
- [Public Ruscigno's APIs](https://github.com/Ruscigno/ruscigno-apis)
- [Generated Protos from Public Ruscigno's APIs](https://github.com/Ruscigno/ruscigno-gosdk)
  
## Building the docker image

`docker build -t r.r6o.co/ticker-signals:latest .`

- `latest/vX.X.X` is the image version number, change it as needed.
- Remove or change `mydocker-registry` for local images.

## $ENV list

- TICKER_SERVER_PORT: server port
- TICKER_DATABASE_URL: database connection string

## Database connection string

You must define an environment variable to connect to the database

- Example DSN: `user=jack password=secret host=pg.example.com port=5432 dbname=mydb sslmode=verify-ca pool_max_conns=10`
- Example URL: `postgres://jack:secret@pg.example.com:5432/mydb?sslmode=verify-ca&pool_max_conns=10`
