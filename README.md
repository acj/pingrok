# ping-heatmap

A tool for creating subsecond offset heatmaps for ICMP echo (ping) replies

![subsecond offset heatmap example](https://user-images.githubusercontent.com/27923/65378522-7699d700-dc87-11e9-83fe-b0ebd0c3251d.png)


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
  -r int
        number of pings per second (default 10)
  -s string
        IP and port for web server (default "0.0.0.0:8086")
  -t int
        seconds of data to display (default 30)
```

Once the program is running, point your browser at http://0.0.0.0:8086 (or whatever you've configured via `-s`).
