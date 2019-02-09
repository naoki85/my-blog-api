# My blog API

## Setup

### Add database.yml

```
$ touch config/database.yml
```
Edit `config/database.yml` as below.

```
development:
  username: root
  password: root
  host: 127.0.0.1
  port: 3306
  database: development_database
production:
  username: hogehoge
  password: fugafuga
  host: 127.0.0.1
  port: 3306
  database: production_database
```

## Basic command

### Run

```
$ go run main.go
```

### Test

```
$ go test
```
