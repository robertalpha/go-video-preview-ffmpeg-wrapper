FROM jrottenberg/ffmpeg:4-alpine

RUN apk add --no-cache go

WORKDIR testdir
COPY . .

ENTRYPOINT ["go", "test", "-v", "./...", "-coverprofile", "cover.out"]
