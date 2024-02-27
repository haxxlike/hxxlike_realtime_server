FROM heroiclabs/nakama-pluginbuilder:3.20.1-arm AS go-builder

ENV GO111MODULE on
ENV CGO_ENABLED 1

WORKDIR /backend
COPY ./go.mod .
COPY ./main.go .
COPY ./vendor/ vendor/