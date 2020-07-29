FROM golang:1.14.2 AS go

WORKDIR $GOPATH/src/github.com/codefresh-id/argocd-listener

COPY . .

RUN go get -f -v ./src

RUN go build -o /argocd-listener ./src

#
# ------ Release ------
#
FROM alpine:3.6

RUN apk --no-cache upgrade && apk --no-cache add ca-certificates

COPY --from=go /argocd-listener /usr/local/bin/

ENTRYPOINT ["argocd-listener"]