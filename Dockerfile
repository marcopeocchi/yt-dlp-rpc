FROM golang:alpine3.16

RUN mkdir -p /usr/src/yt-dlp-rpc/download
VOLUME /usr/src/yt-dlp-rpc/downloads
WORKDIR /usr/src/yt-dlp-rpc

RUN apk add curl psmisc python3 ffmpeg

COPY . .
RUN chmod +x ./fetch-yt-dlp.sh
RUN ./fetch-yt-dlp.sh

RUN go build -o yt-dlp-rpc *.go

EXPOSE 4444
CMD [ "./yt-dlp-rpc", "--driver", "./yt-dlp", "--port", "4444" ]