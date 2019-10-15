FROM golang:1.13 as builder

WORKDIR /go/src/github.com/radhus/vsxknob
COPY . ./

RUN CGO_ENABLED=0 go install

FROM gcr.io/distroless/base
COPY --from=builder /go/bin/vsxknob /vsxknob
ENTRYPOINT ["/vsxknob"]
