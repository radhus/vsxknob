# vsxknob

Publish information from Pioneer VSX receivers on Prometheus and MQTT.

## Usage

```
$ go get -u github.com/radhus/vsxknob
$ vsxknob receiver:8102 mqtt-host:1883
```

For Prometheus metrics, check http://localhost:8080/metrics

For MQTT, subscribe to topic `vsx/state`.

## Docker

Published manually to Docker hub as `radhus/vsxknob:latest`.
