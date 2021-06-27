# Filotto

A simple web game that mimics [Connect Four](https://en.wikipedia.org/wiki/Connect_Four)
that uses [Edd Wise](https://github.com/exelr/eddwise) code generation tools.

## Run the game

Start the server:
```shell
go run github.com/exelr/filotto/cmd/service -w 8080 -s 3000
```

Open in a browser [http://localhost:8080](http://localhost:8080) in a browser, and play against yourself!

If you want to play against a remote player, consider expose your machine with tools like `ngrok` and share the url with them.
