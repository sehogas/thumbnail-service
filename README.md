# Thumbnail service

HTTP service to convert jpg images to thumbnails. Default PORT=3010

### Dev execute

```
    $ make
```

or

```
    $ go run ./cmd/api/main.go
```

### Example 

If width or height is equal to zero, an attempt will be made to resize proportionally.


http://localhost:3010/?width=80&height=0&url=https://upload.wikimedia.org/wikipedia/commons/thumb/f/fb/Pruebas1.jpg/800px-Pruebas1.jpg

## Install as a linux service

```sudo cp thumbnail.service /etc/systemd/system/```

## Activate or deactivate the service

```sudo systemctl enable|disable sigep``` 

## Start or stop the service

```sudo systemctl start|stop sigep```

## Logs

```journalctl -u thumbnail.service --since today```

## Devs

### Example of updating golang version

´´´go mod edit -go=1.22.2´´´

### List available update dependencies

´´´go list -u -m all´´´

### Update all dependencies 

´´´go get -t -u ./...´´´

## Autor
* [Sebastian Hogas](https://github.com/sehogas)