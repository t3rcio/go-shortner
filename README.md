# go-shortner

A simple Golang URL shortener
Running in a Terminal window:

    go run server_http.go

In another Terminal window:

    curl http://localhost:8000/new -d "url:https//gmail.com"
    theshortcutreturned%

Open your browser and type:

      http://localhost:8000/s/<theshortcutreturned>

You will be redirected to the original url ;-)

A very simple documentation resides in:
  http://localhost:8000