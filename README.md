# RTSPtoWeb share you ip camera to world!

RTSPtoWeb converts your RTSP streams to formats consumable in a web browser
like MSE (Media Source Extensions), WebRTC, or HLS. It's fully native Golang
without the use of FFmpeg or GStreamer!

## Table of Contents

- [Installation](#installation)
- [Configuration](#configuration)
- [Command-line](#command-line)
- [Web UI](#web-ui)
- [API documentation](#api-documentation)
- [Limitations](#Limitations)
- [Performance](#Performance)
- [Authors](#Authors)
- [License](#license)

## Installation

### Installation from source

1. Download source
   ```bash
   $ git clone https://github.com/deepch/RTSPtoWeb
   ```
3. CD to Directory
   ```bash
    $ cd RTSPtoWeb/
   ```
4. Test Run
   ```bash
    $ GO111MODULE=on go run *.go
   ```
5. Open Browser
    ```bash
    open web browser http://127.0.0.1:8083 work chrome, safari, firefox
    ```

## Configuration

### Server settings

```text
debug         - enable debug output
log_level     - log level (trace, debug, info, warning, error, fatal, or panic)

http_demo     - serve static files
http_debug    - debug http api server
http_login    - http auth login
http_password - http auth password
http_port     - http server port
http_dir      - path to serve static files from

https
https_auto_tls
https_auto_tls_name
https_cert
https_key
https_port

rtsp_port     - rtsp server port
```

### Stream settings

```text
name          - stream name
```

### Channel settings

```text
name          - channel name
url           - channel rtsp url
on_demand     - stream mode static (run any time) or ondemand (run only has viewers)
debug         - enable debug output (RTSP client)
status        - default stream status
```

#### RTSP pull modes

  * **on demand** (on_demand=true) - only pull video from the source when there's a viewer
  * **static** (on_demand=false) - pull video from the source constantly

### Example config.json

```json
{
  "server": {
    "debug": true,
    "log_level": "info",
    "http_demo": true,
    "http_debug": false,
    "http_login": "demo",
    "http_password": "demo",
    "http_port": ":8083",
    "rtsp_port": ":5541"
  },
  "streams": {
    "demo1": {
      "name": "test video stream 1",
      "channels": {
        "0": {
          "name": "ch1",
          "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
          "on_demand": true,
          "debug": false,
          "status": 0
        },
        "1": {
          "name": "ch2",
          "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
          "on_demand": true,
          "debug": false,
          "status": 0
        }
      }
    },
    "demo2": {
      "name": "test video stream 2",
      "channels": {
        "0": {
          "name": "ch1",
          "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
          "on_demand": true,
          "debug": false,
          "status": 0
        },
        "1": {
          "name": "ch2",
          "url": "rtsp://admin:admin@YOU_CAMERA_IP/uri",
          "on_demand": true,
          "debug": false,
          "status": 0
        }
      }
    }
  }
}
```

## Command-line

### Use help to show available args

```bash
./RTSPtoWeb --help
```

#### Response

```bash
Usage of ./RTSPtoWeb:
  -config string
        config patch (/etc/server/config.json or config.json) (default "config.json")
  -debug
        set debug mode (default true)
```

## Web UI

The web UI is available at http://127.0.0.1:8083/.

### Parameters for full-page multiview

#### Grid size

```
http://localhost:8083/pages/multiview/full?grid=5
```

#### Show additional options

  * Video protocol (MSE, WebRTC, or HLS)
  * Grid size
  * Background image

```
http://localhost:8083/pages/multiview/full?controls
```

## API documentation

See the [API docs](/docs/api.md)

## Limitations

Video Codecs Supported: H264 all profiles

Audio Codecs Supported: no

## Performance

```bash
CPU usage ≈0.2%-1% one (thread) core cpu intel core i7 per stream
```

## Authors

* **Andrey Semochkin** - *Initial work video* - [deepch](https://github.com/deepch)
* **Dmitriy Vladykin** - *Initial work web UI* - [vdalex25](https://github.com/vdalex25)

See also the list of [contributors](https://github.com/deepch/RTSPtoWeb/contributors) who participated in this project.

## License

This project licensed. License - see the [LICENSE.md](LICENSE.md) file for details

[webrtc](https://github.com/pion/webrtc) follows license MIT [license](https://raw.githubusercontent.com/pion/webrtc/master/LICENSE).

[joy4](https://github.com/nareix/joy4) follows license MIT [license](https://raw.githubusercontent.com/nareix/joy4/master/LICENSE).

## Other Example

Examples of working with video on golang

- [RTSPtoWeb](https://github.com/deepch/RTSPtoWeb)
- [RTSPtoWebRTC](https://github.com/deepch/RTSPtoWebRTC)
- [RTSPtoWSMP4f](https://github.com/deepch/RTSPtoWSMP4f)
- [RTSPtoImage](https://github.com/deepch/RTSPtoImage)
- [RTSPtoHLS](https://github.com/deepch/RTSPtoHLS)
- [RTSPtoHLSLL](https://github.com/deepch/RTSPtoHLSLL)

[![paypal.me/AndreySemochkin](https://ionicabizau.github.io/badges/paypal.svg)](https://www.paypal.me/AndreySemochkin) - You can make one-time donations via PayPal. I'll probably buy a ~~coffee~~ tea. :tea:
