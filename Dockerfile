FROM golang:1.23-bullseye

ENV APP_DIR=/app
WORKDIR ${APP_DIR}

# Requirements are installed here.
RUN apt -y update \
    && apt -y install protobuf-compiler clang-format libc-dev libstdc++-10-dev libgcc-9-dev netbase libudev-dev ca-certificates default-mysql-client fish
RUN sh -c "$(curl -Ssf https://pkgx.sh)" \
    && pkgm install task \
    && pkgm install fullstory.com/grpcurl
