FROM golang:1.10

WORKDIR $GOPATH/src/github.com/codefresh-id/argocd-listener

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

RUN go build -o argocd-listener ./src

EXPOSE 8080

CMD ["argocd-listener"]