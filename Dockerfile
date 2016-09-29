FROM alpine:3.4
RUN apk add --update bash ca-certificates && rm -rf /var/cache/apk/*
ADD artifact /
ENTRYPOINT ["/artifact"]
EXPOSE 50050
