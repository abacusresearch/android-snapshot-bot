FROM golang
ADD . /go/src/github.com/abacusresearch/android-snapshot-bot
RUN go get github.com/abacusresearch/android-snapshot-bot/...
RUN go install github.com/abacusresearch/android-snapshot-bot
ENTRYPOINT /go/bin/android-snapshot-bot
