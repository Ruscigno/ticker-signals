# Ticker Heart

## Disclosure
First of all, this is a POC, so I'm not worried to code the perfect solution. 
There are not enough unit tests, there is duplicated code, there is a lot of 
things I usually don't do professionally. My main goal here is to exercise 
my skills and to have a Portfolio of interesting things.

However, I already started to work on the debts and on the TODO list. I'm 
planning to have at least 60% of code coverage so that I can automate the 
building steps which are done manually for now.

## Things I'm willing to use here
 - Self-hosted Kubernetes cluster, using Raspberry PIs
 - Use of Google Protobuf
 - Use of some free GCP resources, like Container Registry
 - Automatic building (probably the last thing I gonna do)
 - Self-hosted messaging protocol
 - Self-hosted Datadog alike
 - Self-hosted log ingestion server
 - GRPc protocol
 - My own Health Check solution
 - Redis as a cache for tradings
 - Work as a Deal Replicator for other accounts
 - more to be defined...

As you can see, there is a lot of stuff I'm willing to use, probably it's not
gonna be the fastest solution, but as I said it's just a funny exercise.

# Main Idea
This app is a GRPc Server API that will receive messages from the [ticker-beats](https://github.com/Ruscigno/ticker-beats). 
These messages contain deals, orders, positions, accounts info, and trades from a Metatrader 5's 
Experts Advisors. These robots are not open to the general public neither the code that generates 
the log files nor the code to manage all the infrastructure needed to handle this Solution.

It can handle multiple accounts and receive duplicated messages, in other 
words, it's fault-tolerant for duplicated deals, orders, positions, and accounts
messages, it's stateless vertical and horizontal scalable microservice.

It's also meant to do more than that, but you'll have to discover it by yourself... :P

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
