# ilab1

## Install

  ```
  git clone github.com/omegatymbjiep/ilab1
  cd ilab1
  go build main.go
  export KV_VIPER_FILE=./config.yaml
  ./main migrate up
  ./main run service
  ```

## Running from docker 
  
Make sure that docker installed.

use `docker run ` with `-p 8080:80` t expose port 80 to 8080

  ```
  docker build -t github.com/omegatymbjiep/ilab1 .
  docker run -e KV_VIPER_FILE=/config.yaml github.com/omegatymbjiep/ilab1
  ```

## Running from Source

* Set up environment value with config file path `KV_VIPER_FILE=./config.yaml`
* Provide valid config file
* Launch the service with `migrate up` command to create database schema
* Launch the service with `run service` command


### Database
For services, we do use ***PostgresSQL*** database. 
You can [install it locally](https://www.postgresql.org/download/) or use [docker image](https://hub.docker.com/_/postgres/).
