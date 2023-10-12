# floetenleague

Floetenleague includes a bingo board and POE Oauth2 login.

## Run it

* Install go >=1.19
* Install yarn
* Install nodejs
* Install docker
* Install dependencies
    * goto ui directory
    * run `yarn install`
* Generate UI stuff: 
    * goto ui directory
    * Run `yarn generate`
* Build UI
    * goto ui directory
    * run `yarn build`
* Generate Backend stuff: Run `go generate`
* Create a file called `.env` with this content:
  ```
  FL_DB_CONNECTION=postgres://root:root@localhost/floetenleague?sslmode=disable
  FL_DEBUG_INSERT=true
  ```
* Start Postgres (can be done with `docker-compose up -d` in this directory)
* Run backend: `go run . .env`, where `.env` is the file created above
* Start frontend in development mode:
    * goto ui directory
    * run `yarn start`
