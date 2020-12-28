# ping-heatmap

A tool for creating subsecond offset heatmaps for ICMP echo (ping) replies

[![ping-heatmap asciicast demo](https://asciinema.org/a/381549.svg)](https://asciinema.org/a/381549)

## Building

```
$ go build
```

## Usage

```
$ ./ping-heatmap --help
Usage of ./ping-heatmap:
  -h string
    	the host to ping (default "192.168.1.1")
  -o	Overlay latencies on heatmap
  -r int
    	number of pings per second (default 10)
  -t int
    	seconds of data to display (default 30)
```
