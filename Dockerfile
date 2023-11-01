FROM registry.access.redhat.com/ubi8/go-toolset:1.19.13-2.1697656138 AS builder

COPY . .

RUN CGO_ENABLED=0 go build -buildvcs=false -o ./sigterm-watcher .

FROM registry.access.redhat.com/ubi9/ubi:9.2-755.1697625012

COPY --from=builder --chown=1000:0 --chmod=775 /opt/app-root/src/sigterm-watcher /usr/local/bin/sigterm-watcher

CMD ["/usr/local/bin/sigterm-watcher"]
