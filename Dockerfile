# Package Merge Driver Builder
FROM golang:1.14 as builder
RUN go build -o bin/package-merge-driver

# Node, Git
FROM node:10.21.0

RUN apt-get update \
    && apt-get install -y git \
    && apt-get clean \

USER node

ADD --from=builder /user/src/myapp/bin/package-merge-driver /home/node/package-merge-driver

RUN git config --global merge."packagemerge-driver".name "Automatically merge npm semantic version of package.json" \
    && git config --global merge."packagemerge-driver".driver "/home/node/package-merge-driver %O %A %B"

CMD ["node"]