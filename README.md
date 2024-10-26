# go-shortner

A simple Golang URL shortener
You will need a Redis Server. To run it follow the instructions here:
https://redis.io/docs/latest/operate/oss_and_stack/install/install-stack/docker/
After that, run in a Terminal window:

    $ go run server_http.go

In another Terminal window:

    $ curl http://localhost:8000/new -d "url=https//gmail.com"
    theshortcutreturned%

Open your browser and type:

      http://localhost:8000/s/<theshortcutreturned>

You will be redirected to the original url ;-)

A very simple documentation resides in:
  http://localhost:8000

## Docker
To run the shortener in a container:

    $ docker-compose -up -d
Access your brower: http://localhost:8080
;-)