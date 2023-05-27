FROM golang:1.19.0-alpine3.16

RUN apk update \
    && apk add curl logrotate

COPY . .

ENV GOPATH=${pwd}
# RUN go mod tidy
RUN cd ./ && unset GOPATH && go get

CMD go run *.go
