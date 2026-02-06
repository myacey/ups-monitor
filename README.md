# ups-monitor

## usage

run:
```sh
./ups-monitor --port=<serial port> --freq=<frequency to get port data in ms>
```
> example:
> `./ups-monitor --port=COM7 --freq=5000`

help message:
```sh
./ups-monitor --help
```

run with go:
```go
go run . --port=<serial port> --freq=<frequency to get port data in ms>
```