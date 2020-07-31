FROM golang:1.14.2 AS go

# set working directory
RUN mkdir -p /go/src/github.com/codefresh-io/argocd-listener
WORKDIR /go/src/github.com/codefresh-io/argocd-listener

COPY . .

RUN go get -f -v ./agent/pkg

RUN CGO_ENABLED=0 go build -o /argocd-listener ./agent/pkg

#
# ------ Release ------
#
FROM alpine:3.6

RUN apk --no-cache upgrade && apk --no-cache add ca-certificates

COPY --from=go /argocd-listener /usr/local/bin/

CMD ["argocd-listener"]