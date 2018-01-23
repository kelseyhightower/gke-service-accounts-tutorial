FROM golang:1.9.1
WORKDIR /go/src/github.com/kelseyhightower/gke-service-accounts-tutorial/echo
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o echo \
    -tags netgo -installsuffix netgo .

FROM scratch
COPY --from=0 /go/src/github.com/kelseyhightower/gke-service-accounts-tutorial/echo/echo .
ENTRYPOINT ["/echo"]
