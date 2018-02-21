FROM golang:1.10 as builder

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

WORKDIR /go/src/github.com/radhus/vsxknob
COPY . ./
RUN dep ensure
RUN CGO_ENABLED=0 go install \
      -a -tags netgo \
      -ldflags '-extldflags "-static"'


FROM scratch
COPY --from=builder /go/bin/vsxknob /vsxknob
ENTRYPOINT ["/vsxknob"]
