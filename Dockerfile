ARG ALPINE_VERSION=3.15.0
ARG GOLANG_VERSION=1.17.7-alpine3.15
ARG VERSION
ARG GIT_COMMIT

FROM golang:${GOLANG_VERSION} as builder

WORKDIR /go/src/metrp/
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ENV CGO_ENABLED=0

RUN go install \
    -installsuffix "static" \
    ./...

FROM alpine:${ALPINE_VERSION} as runtime

RUN set -x \
  && apk add --update --no-cache ca-certificates tzdata \
  && echo 'Etc/UTC' > /etc/timezone \
  && update-ca-certificates

ENV TZ=/etc/localtime                  \
    LANG=en_US.utf8                    \
    LC_ALL=en_US.UTF-8

COPY --from=builder /go/bin/metrp /
RUN chmod +x /metrp

RUN adduser -S appuser -u 1000 -G root
USER 1000

ENTRYPOINT ["/metrp"]
