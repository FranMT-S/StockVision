## StockVision API

This api provides the search functionality for the StockVision application, the mails are stored in a cockroachdb database, the data must be loaded using the command `go run main.go fill-db`.


## Tech stack
- Language: golang v1.24
- API framework: chi 
- Database: cockroachdb
- ORM: gorm
- Logger: zerolog
- Environment variables: godotenv
- cache: redis

## Folder structure

```
├── cache: redis root, cache interface and redis implementation
├── cmd: cobra cmd with commands to fill the database
├── config: class files to config the application
├── controllers: Controller HTTP files
├── database: connections to the database
├── http: examples how use the API Endpoints
├── logger: implementation of zerolog to logs  
├── logs: directory where the logs are stored
├── models: Data models and interfaces,filters, ratings, responses 
├── routes: Api endpoints
├── sanatizer: Utility to sanitize the data
├── server: Main server implementation in chi
├── services: Utility to handle the business logic
```

## Enviroment variables
Create a .env file based on .env.template, fill the values with your own.

``` env
ENV=development # production or dev, set production in cloud
DB_HOST=localhost # Database host
DB_NAME=defaultdb # Database name
DB_USER=root # Database user
DB_PASSWORD= # Database password
DB_PORT=26257 # Database port
DB_SSL=false # Database SSL mode
DB_SCHEMA=stocksvision # Database schema
REDIS_HOST=localhost # redist host
REDIS_PORT=6379 # redist port
REDIS_PASSWORD= # redist password
API_PORT=8080 # API port
CLIENT_HOST="http://localhost:5173" # Client host to cors
LOG_LEVEL=info # Log level, Options: trace, debug, info, warn, error, dpanic, panic, fatal
LOG_DB=false # Log database, used to debug queries
STOCK_API_URL= # Stock to get the recommendations API url
STOCK_API_TOKEN= # Stock API token
FINANCIAL_BASE_URL= # Financial API url
FINANCIAL_TOKEN= # Financial API token
FINHUB_BASE_URL= # Finhub API url
FINHUB_TOKEN= # Finhub API token
GEMINI_API_KEY= # Gemini API key
```

## Installation

```bash
go mod download
go mod init
go mod tidy
```

## Run
Configure the environment variables and run the application.

**First time must populate the database** 
if you have a json file with the recommendations data you can use the 
command

```bash
go run main.go fill-db --json data/recommendations.json
```

if you have a api key to get the recommendations data you can use the 

```bash
go run main.go fill-db
```

then can run the application
**Run the application**
```bash
go run main.go
```

## Endpoints

### GET /api/v1/tickers
Get the tickers of a company


``` http
GET /api/v1/tickers?page=1&sort=asc&q=&size=10
```

``` http
GET /api/v1/tickers/AAPL/overview
```

``` http
GET /api/v1/tickers/AAPL/predictions
```

``` http
GET /api/v1/tickers/AAPL/logo
```

