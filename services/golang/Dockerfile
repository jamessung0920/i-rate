FROM golang

ARG APP_ENV
ENV GO111MODULE=on

WORKDIR /go/src/app
COPY www/golang .

# RUN go get github.com/pilu/fresh
# CMD ["fresh"]

RUN go get ./ && \
    go build

CMD if [ ${APP_ENV} = production ]; \
    then \
    app; \
    else \
    go get github.com/pilu/fresh && \
    fresh; \
    fi

EXPOSE 8080