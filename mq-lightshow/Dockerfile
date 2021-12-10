ARG BUILD_VERSION
ARG ARCH
ARG GOOS
ARG GOARCH
FROM $ARCH/golang:1.17-alpine AS build

RUN apk update && apk add --no-cache libc-dev gcc

WORKDIR /go/src/github.com/lovesway/hassio-addons/mq-lightshow
COPY . .

RUN CGO_ENABLED=1 GOOS=$GOOS GOARCH=$GOARCH go build -a -ldflags "-X main.version=$BUILD_VERSION -linkmode external -extldflags '-static'" -o /mq-lightshow .

FROM scratch
COPY --from=build /mq-lightshow /
COPY www/*.tpl /www/
COPY www/css/*.css* /www/css/
COPY www/js/*.js* /www/js/
COPY www/icons/*.svg /www/icons/
COPY www/images/*.png /www/images/

EXPOSE 8099/tcp
ENTRYPOINT ["/mq-lightshow"]