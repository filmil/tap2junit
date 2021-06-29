
FROM golang:1.16.5 as build
WORKDIR /go/src/github.com/filmil/tap2junit
COPY . .
RUN go build ./cmd/tap2junit
RUN chmod 755 tap2junit

FROM gcr.io/distroless/base
LABEL maintainer="filmil@gmail.com"
COPY --from=build /go/src/github.com/filmil/tap2junit/tap2junit /
ENTRYPOINT ["/tap2junit"]
