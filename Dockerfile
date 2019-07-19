FROM gcr.io/distroless/base
LABEL maintainer="filmil@gmail.com"
COPY tap2junit /
ENTRYPOINT ["/tap2junit"]
