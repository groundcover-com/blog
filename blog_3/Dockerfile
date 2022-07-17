FROM ubuntu:focal

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y apt-transport-https ca-certificates gnupg curl unzip

ENV VERSION="1.18"
ENV ARCH="amd64"
ENV PATH=$PATH:/usr/local/go/bin

RUN curl --silent -O -L "https://golang.org/dl/go${VERSION}.linux-${ARCH}.tar.gz" && \
    tar -C /usr/local -xzf "go${VERSION}.linux-${ARCH}.tar.gz" && \
    rm "go${VERSION}.linux-${ARCH}.tar.gz"
    
RUN curl --silent "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && \
    unzip -qq awscliv2.zip && \
    ./aws/install && \
    rm awscliv2.zip

WORKDIR /builder
COPY pkg/ pkg/
COPY go.mod go.sum main.go ./
RUN go mod download
RUN go build -o /bin/restart-watcher ./main.go

ENTRYPOINT ["/bin/restart-watcher"]
