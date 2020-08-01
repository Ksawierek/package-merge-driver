# Package Merge Driver Builder
FROM golang:1.14 as builder
WORKDIR /usr/src/app
COPY package-merge-driver.go .
RUN go mod init github.com/Ksawierek/package-merge-driver && go build -o bin/package-merge-driver

# Node, Git
FROM node:12.18.3

RUN wget --no-check-certificate -q https://raw.githubusercontent.com/petervanderdoes/gitflow-avh/develop/contrib/gitflow-installer.sh && bash gitflow-installer.sh install stable; rm gitflow-installer.sh

USER node

COPY --from=builder /usr/src/app/bin/package-merge-driver /home/node/package-merge-driver

RUN git config --global merge."packagemerge-driver".name "Automatically merge npm semantic version of package.json" \
    && git config --global merge."packagemerge-driver".driver "/home/node/package-merge-driver %O %A %B"

CMD ["node"]