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

If width or height is equal to zero, an attempt will be made to resize proportionally


http://localhost:3010/?width=80&height=0&url=http%3A%2F%2F192.168.3.11%2FfotosPers%2FDNI_30128996_f5de20c5-8108-43a9-9342-d2d300b910e3.jpg
