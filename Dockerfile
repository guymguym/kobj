FROM alpine:3

RUN apk add --update ca-certificates curl bash \
 && curl -L https://storage.googleapis.com/kubernetes-release/release/`curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt`/bin/linux/amd64/kubectl -o /usr/local/bin/kubectl \
 && chmod +x /usr/local/bin/kubectl \
 && rm /var/cache/apk/*

WORKDIR /
COPY ./kobj-linux ./kobj

ENTRYPOINT [ "/kobj" ]
CMD [ "server" ]
