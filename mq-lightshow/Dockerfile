FROM golang:1.17-alpine AS build
ARG BUILD_VERSION

RUN apk update && apk add --no-cache libc-dev gcc

WORKDIR $GOPATH/src/gitlab.local/hassio-addons/mq-lightshow/
COPY . .

#RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -ldflags "-X main.version=$BUILD_VERSION -linkmode external -extldflags '-static'" -o /mq-lightshow .
RUN CGO_ENABLED=1 go build -a -ldflags "-X main.version=$BUILD_VERSION -linkmode external -extldflags '-static'" -o /mq-lightshow .

FROM scratch
COPY --from=build /mq-lightshow /
COPY www/*.tpl /www/
COPY www/css/*.css* /www/css/
COPY www/js/*.js* /www/js/
COPY www/icons/*.svg /www/icons/
COPY www/images/*.png /www/images/

EXPOSE 8099/tcp
ENTRYPOINT ["/mq-lightshow"]