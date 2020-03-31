FROM golang:1.14.1-alpine

# install required packages
RUN apk add --update --no-cache \
	build-base \
	ca-certificates \
	curl \
	git \
	groff \
	jq \
	less \
	make \
	mysql-client \
	openssh-client \
	openssl \
	py-pip \
	python \
	python-dev \
	unzip \
	vim \
	zip

RUN go get github.com/golang/mock/mockgen

RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.24.0

RUN git config --global url."git@github.com:".insteadOf "https://github.com/"
