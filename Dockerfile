FROM alpine:3.4
MAINTAINER "Ikem Okonkwo <ikem.okonkwo@andela.com>"

COPY artifact /artifact

CMD ["/artifact"]

EXPOSE 50050