# ping-heatmap

![Go](https://github.com/acj/ping-heatmap/workflows/CI/badge.svg)

A tool for creating subsecond offset heatmaps for ICMP echo (ping) replies

## Demo

[![ping-heatmap asciicast demo](https://asciinema.org/a/381549.svg)](https://asciinema.org/a/381549)

## Why

I noticed that occasionally my MacBook Pro's wifi latency would become spiky. It was most noticeable when I was using interactive sessions like SSH, but it was degrading my download throughput too. Using ICMP echo (ping) requests and plotting their response time on a [subsecond offset heatmap](http://www.brendangregg.com/HeatMaps/subsecondoffset.html), I could visualize the latency and confirm that it spiked from ~5ms to ~300ms (or more) every half second. I eventually traced it to a specific VirtualBox VM. In the process, this tool was born.

## How it works

Most of the interesting bits happen in the `pinger` module, which continuously sends ICMP echo packets to a host with a frequency that you specify. When the responses arrive, the pinger records the timestamp and the delay (latency), and then that data is sent to the `dataPointPartitioner` to be divided into one-second buckets. Those buckets are then displayed as a heatmap.

Visually, it's important to understand that _both_ axes are displaying time. The x-axis shows a rolling window of passing time, which is common, but the y-axis shows the latency at various points _within each second_. If you're dealing with a problem that occurs very briefly and/or more frequently than once per second, as I was, then this extra resolution is critical.

## Installing

If you have the Go toolchain installed, `go get` should work:

```
go get -u github.com/acj/ping-heatmap
```

## Building from source

Clone this repository, and then:

```
go build
```

and optionally:

```
go test ./...
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
