FROM golang:1.23-bullseye

ENV APP_DIR /app
WORKDIR ${APP_DIR}

# Requirements are installed here.
RUN apt -y update && apt -y install protobuf-compiler
RUN sh -c "$(curl -Ssf https://pkgx.sh)" \
    && pkgx install task \
    && pkgx install mysql.com \
    && pkgx install fishshell.com
