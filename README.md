# vsxknob

Expose Prometheus compatible metrics from some Pioneer VSX receivers.

## Usage

```
$ go get -u github.com/radhus/vsxknob
$ vsxknob receiver:8102
```

... and check http://localhost:8080/metrics


Also published manually to Docker hub as `radhus/vsxknob:latest`.
