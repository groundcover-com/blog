FROM ubuntu:focal

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y apt-transport-https ca-certificates gnupg curl unzip 

RUN echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | \
    tee -a /etc/apt/sources.list.d/google-cloud-sdk.list && \
    curl --silent https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add -

RUN apt-get update && \
    apt-get install -y google-cloud-cli

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
