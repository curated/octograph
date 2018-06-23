FROM golang:1.10 AS build
WORKDIR $GOPATH/src/github.com/curated/octograph
ADD https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 /usr/bin/dep
RUN chmod +x /usr/bin/dep
COPY Gopkg.toml Gopkg.lock ./
RUN dep ensure --vendor-only
COPY . ./
RUN ln -s $(pwd) /octograph
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o /app .

FROM alpine
# ENV CONFIG=config/prod.config.json
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /octograph/config/dev.config.json config/
COPY --from=build /octograph/graph/issues_query.gql graph/
COPY --from=build /octograph/mapping/issue.json mapping/
COPY --from=build /app ./
CMD ["/app", "-logtostderr=true"]
