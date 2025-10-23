# StockVision
StockVision is a web application that provides stock market visualization, recommendations and predictions.

# Pre-requisites
- get a api key from `GeminiAi` [here](https://aistudio.google.com/app/api-keys)
- get a api key from `Finnhub` [here](https://finnhub.io/)
- get a api key from `FinancialModelingPrep` [here](https://financialmodelingprep.com/)


## Requirements

* [CockroachDB](https://www.cockroachlabs.com/docs/v25.3/install-cockroachdb-windows.html)
* [Go](https://golang.org/dl/)
* [Vue 3](https://vuejs.org/guide/introduction.html)
* [Node.js](https://nodejs.org/en/download/)
* [NPM](https://www.npmjs.com/)
* [PNPM](https://pnpm.io/installation/)
* [Redis](https://redis.io/)

## Database

The database is a [cockroachdb](https://www.cockroachlabs.com/docs/v25.3/install-cockroachdb-windows.html) the project using code first approach to create the database schema.

## Tech stack
- API: golang
- Database: cockroachdb
- Client: vue 3 + pinia + vite

## API [🔗](./api)
* Language: golang
* Cmd: cobra
* Framework: chi
* Database: cockroachdb
* ORM: gorm
* Logger: zerolog
* Environment variables: godotenv
* Cache: redis

# Client [🔗](./client)
* Language: vue 3
* Framework: vite
* State management: pinia
* Router: vue-router
* Environment variables: dotenv

## Infrastructure [🔗](./terraform)
* Language: terraform
* Cloud: AWS
* Diagram: Mermaid [🔗](./terraform#aws-infrastructure-diagram)

## Folder structure
```
├── api: API project
├── terraform: Infrastructure of the project aws
├── client: Client project
```

## Documentation
* [API](./api)
* [Client](./client)
* [Infrastructure](./terraform)

## Docker compose
you can use docker compose to run the api and the client, after setup de environment variables in the .env file of the api and client folders.

```bash
docker-compose up
```

after run you must be able to access the client at http://localhost:5173 and the api at http://localhost:8080 if use the default ports.

this also will create a instance of cockroachdb in the cloud, you can use the connection string to connect to the database.

# Run

1. run the api
The first time you run the api you must run the `fill-db` command to populate the database.

```bash
cd api
go run main.go --fill-db
```

2. run the client
```bash
cd client
pnpm run dev
```

