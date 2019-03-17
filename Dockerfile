FROM golang:1.12

WORKDIR $GOPATH/src/github.com/tsoporan/music_findr_api

COPY . .

# Deps
RUN go get -d -v ./...

RUN go install -v ./...

EXPOSE 4000

CMD ["music_findr_api"]
