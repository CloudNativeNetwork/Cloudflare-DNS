FROM golang:1.18 as build

WORKDIR /build

ADD go.mod go.mod

ADD go.sum go.sum

RUN go mod download -x

ADD . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -trimpath -installsuffix cgo -o cloudflare-dns main.go

FROM scratch

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=build /build/cloudflare-dns /cloudflare-dns
