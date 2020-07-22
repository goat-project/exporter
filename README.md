# Goat exporter

[![Build Status](https://travis-ci.org/goat-project/exporter.svg?branch=master)](https://travis-ci.org/goat-project/exporter)

Goat exporter exports *Storage*, *Public IP* and *Cloud Usage* records to [Prometheus](https://prometheus.io/).

![exporter](https://github.com/goat-project/exporter/blob/master/img/exporter.png)

The goat exporter is [a service](https://github.com/goat-project/exporter/tree/master/service) that waits for records using a Watcher. 
The [Watcher](https://github.com/goat-project/exporter/tree/master/watch) watches a root directory given by a configuration and its subdirectories. 
When a new directory is created, Watcher adds it to the list of watched directories. When a new record is written, 
Watcher adds it to the Event channel. The event channel is handled by Parser.

The [Parser](https://github.com/goat-project/exporter/tree/master/parse) takes the event, opens a file given by the event, distinguishes 
the file format, and parses it. The Parser recognizes 3 file formats:
- `"text/plain; charset=utf-8"` - for Cloud Usage record (The first line has to contain APEL message.)
- `"application/json"` - for Public IP Usage record
- `"text/xml; charset=utf-8"` - for Storage Usage record

The file data is processed by a respective parser - IP, Storage, VM. IP uses Go encoding library for JSON, 
Storage uses Go encoding library for XML and VM is parsed manually according to APEL format. Parsed records are put to 
the Record channel which is handled by Exporter.

The [Exporter](https://github.com/goat-project/exporter/tree/master/export) takes the record and exports it to the Prometheus 
according to its type. Export is provided by a respective gauge. The [Gauges](https://github.com/goat-project/exporter/tree/master/gauge) 
must be registered in Prometheus before exporting and satisfy the correct format with all registered labels.


## Requirements
* Go 1.12 or newer to compile
* Promethues, Grafana

## Installation
The recommended way to install this tool is using `go get`:
```
go get -u github.com/goat-project/exporter
```

## Configuration
Exporter is configured by a file or command line flags.
```
Flags:
  -d, --debug string                 debug
  -o, --dir-path string              Directory path [PATH] (required)
  -g, --goat-endpoint string         Goat endpoint [GOAT_ENDPOINT] (required)
  -h, --help                         help for exporter
      --log-path string              path to log file
  -p, --prometheus-endpoint string   Prometheus endpoint [PROMETHEUS_ENDPOINT] (required)
  -v, --version                      version for exporter
```
The default configuration file is in [`config/` folder](https://github.com/goat-project/exporter/tree/master/config). 
The exporter configuration, named `exporter.yml`, could be also placed in `/etc/exporter/` or `$HOME/.exporter/`.

## Usage example
```
go build
./exporter
```

## Contributing
1. [Fork exporter](https://github.com/goat-project/exporter/fork)
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request