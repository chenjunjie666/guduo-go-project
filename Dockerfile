FROM centos:8

RUN yum -y update

RUN yum -y install wget

#WORKDIR /go/src/guduo
RUN wget https://studygolang.com/dl/golang/go1.16.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf go1.16.linux-amd64.tar.gz

RUN export PATH=$PATH:/usr/local/go/bin && \
# set proxy for faster downloading
    go env -w GOPROXY=https://goproxy.cn,direct

#RUN mkdir /app
#COPY . ./app
#RUN cd /app && go mod download
#RUN cd /app/app/crawler/clean && go build main.go