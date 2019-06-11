
# Production Tips Server

A Go app.

## Running Locally

Make sure you have [Go](http://golang.org/doc/install) version 1.12 or newer.

```sh
$ git clone https://github.com/dadjeibaah/production-tips-server.git
$ cd production-tips-server
$ go build -o bin/production-tips-server -v .
...
bin/production-tips-server
```

The server should now be running on [localhost:8000](http://localhost:5000/).

## Deploying to Heroku

```sh
$ heroku create
$ git push heroku master
$ heroku open
```

or

[![Deploy](https://www.herokucdn.com/deploy/button.png)](https://heroku.com/deploy)


## Documentation

For more information about using Go on Heroku, see these Dev Center articles:

- [Go on Heroku](https://devcenter.heroku.com/categories/go)
