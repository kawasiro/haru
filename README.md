# Haru

[![Build Status](https://travis-ci.org/if1live/haru.svg?branch=master)](https://travis-ci.org/if1live/haru)

Comic crawler

## Feature
* Hitomi.la
  * Download
  * Browsing
  * Enqueue download task
  * API to fetch information

## Install

Before run haru, install goautoenv and godep.

* [goautoenv](https://github.com/Perlmint/goautoenv)
* [godep](github.com/tools/godep)

```bash
go get github.com/Perlmint/goautoenv
go get github.com/tools/godep

git clone git@github.com:if1live/haru.git
cd haru

goautoenv init haru
source .goenv/bin/activate
godep restore
goautoenv link
```

## Run Command line interface
```bash
# compile
go build

# download id=405092 from hitomi.la
./haru -id=405092 -service=hitomi -cmd=download

# help
./haru -h
```

Use enviromnent variable(key=ID) to pass target id.

## Run Server
```bash
# compile
go build

# Run Server
./haru

# Use enviroment variable to change port number.
# PORT=1234 ./haru
```

Connect to http://127.0.0.1:3000

## URL
* http://127.0.0.1:3000/
  * Gallery List
* http://127.0.0.1:3000/detail.html?id=405092
  * Gallery detail page. Download or go to original reader page.

## API
### Gallery List

GET /api/list/{service}/

example
* http://127.0.0.1:3000/api/list/hitomi/?page=2&tag=female:stockings&language=korean
* http://127.0.0.1:3000/api/list/hitomi/?artist=hiten%20onee-ryuu&language=korean

| Type | Name | Description | Example |
|------|------|-------------|---------|
| URL  | service | gallery site name | hitomi |
| Query String | page | page number. 1~xxx | 2 (default=1) |
| Query String | language | Language | korean |
| Query String | tag | tag name | female:stockings |
| Query String | artist | artist name  | hiten%20onee-ryuu |

Note
* Cannot use tag and artist at same time.

### Gallery Information
GET /api/detail/{service}/<id>

example
* http://127.0.0.1:3000/api/detail/hitomi/405092

| Type | Name | Description | Example |
|------|------|-------------|---------|
| URL  | service | gallery site name | hitomi |
| URL  | id | gallery id | 405092 |

### Download Gallery
GET /api/download/{service}/{id}

example
* http://127.0.0.1:3000/api/download/hitomi/405092

Note
* parameters are same with gallery information API.

### Enqueue download task
GET /api/enqueue/<service>/<id>

example
* 127.0.0.1:3000/api/enqueue/hitomi/405092

If you enqueue task, haru download cover images and gallery images.
You can download zip after work is completed.

Note
* parameters are same with gallery information API.


## Unit test

```bash
go test ./...
```

## TODO
* Download queue
* Cron to watch feed
* Save browser's history
* Pretty UI
