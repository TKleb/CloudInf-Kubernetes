*Simple three tier application using new technologies.*
# Three Tier Application
Simple three tier application. The web application requests the api, which then makes a query to the connected database system. The application can only work if all components are configured correctly.
```
  ┌──────────────┐       ┌──────────────┐      ┌──────────────┐
  │              │       │              │      │              │
  │              │       │              │      │              │
  │              │       │              │      │              │
  │     web      ├──────►│     api      ├─────►│   database   │
  │              │       │              │      │              │
  │              │       │              │      │              │
  │              │       │              │      │              │
  └──────────────┘       └──────────────┘      └──────────────┘
```
## Web
Simple web application written with NodeJS. The application starts the webserver on port `8080`.

### Configuration
The web application can be configured with several environment variables listed below:
|Env Variable|Description|Mandatory|
|---|---|---|
|MODULE_NAME|Module Name|Yes|
|GROUP_NAME|Group Name|Yes|
|API_HOST|API Host Address|Yes|
|API_PORT|API Host Port|Yes|

### Develop
1. Install packages
    ```bash
    $ npm install
    ```
2. Start server
    ```bash
    $ npm start
    ```

## API
Simple api application written in Go. The application starts the webserver on port `8000`.

|Path|Description|
|---|---|
|`/time`|Makes a query to get the time from the connected database.|
|`/healthz`|Delivers a response which indicates the health of the application.|

### Configuration
The api application can be configured with several environment variables or directly over dedicated files. Environment variables will have precedence over the dedicated files. 

The following environment variables / file paths can be used:
|Env Variable|Filepath|Description|Mandatory|
|---|---|---|---|
|PG_USER|`/var/config/pg_user`|PostgreSQL User|Yes|
|PG_PASSWORD|`/var/config/pg_password`|PostgreSQL Password|Yes|
|PG_DATABASE|`/var/config/pg_database`|PostgreSQL Database|Yes|
|PG_HOST|`/var/config/pg_host`|PostgreSQL Host Address|Yes|
|PG_PORT|`/var/config/pg_port`|PostgreSQL Host Port|No (defaults to: `5432`)|
|LOG_PATH|`/var/config/log_path`|Log File Path|No (defaults to: `/var/log/logs.log`)|
|DEBUG|`/var/config/debug`|Boolean Value|No (when set to true all logs will be printed to `stdout`)|

### Develop 
1. Install dependencies
    ```bash
    $ go mod tidy
    ```
2. Start Application
    ```bash
    $ go run main.go
    ```

### Build
```bash
$ cd api
$ CGO_ENABLED=0 go build -a -o ./bin/api
```

## Database
For the database a single `PostgreSQL` database is used.  
  
Have a look at the [Docker Hub](https://hub.docker.com) for the current tags and images.
