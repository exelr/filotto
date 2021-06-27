# Filotto

![Filotto](screenshot-small.png)

A simple web game that mimics [Connect Four](https://en.wikipedia.org/wiki/Connect_Four)
that uses [Edd Wise](https://github.com/exelr/eddwise) code generation tools.

## Run the game

If you have `go` >= 1.16, start the server directly:
```shell
go run github.com/exelr/filotto/cmd/service -w 8080 -s 3000
```

Otherwise you can use Docker:

```shell
docker build . -t filotto
docker run -p8080:8080 -p3000:3000 filotto
```


Open in a browser [http://localhost:8080](http://localhost:8080) in a browser, and play against yourself!

If you want to play against a remote player, consider expose your machine with tools like `ngrok` and share the url with them.
