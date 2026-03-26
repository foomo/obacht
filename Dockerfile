FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN adduser -D -u 1001 -g 1001 obacht

COPY obacht /usr/bin/

USER obacht
WORKDIR /home/obacht

ENTRYPOINT ["obacht"]
