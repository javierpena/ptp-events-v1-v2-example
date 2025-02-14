FROM ubi9/ubi:latest AS builder
WORKDIR /tmp/simpleapp
COPY . .

FROM ubi9/ubi:latest as bin

COPY --from=builder /tmp/simpleapp/bin/simpleapp /
ENTRYPOINT ["/simpleapp"]
