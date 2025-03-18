FROM golang:alpine AS build-img
RUN mkdir /build_output/
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN go build -o /build_output/sound_server

FROM alpine AS final-image
RUN mkdir /sounds/
RUN apk add --no-cache alsa-utils
WORKDIR /app
COPY --from=build-img /build_output/sound_server /app/sound_server

EXPOSE 8500
CMD ["/app/sound_server", "--sounds-dir=/sounds", "--timeout-seconds=600"]