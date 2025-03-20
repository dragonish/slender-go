FROM golang:alpine3.21 AS Builder
ARG VERSION
ARG COMMIT

RUN apk add git bash gcc musl-dev upx

WORKDIR /app

ENV GO111MODULE=on
ENV CGO_ENABLED=1
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN export BUILD_DATE=`date +%FT%T%z` && \
  echo "$BUILD_DATE / $COMMIT / $VERSION" && \
  go build -ldflags "-w -s -X 'slender/internal/version.Version=$VERSION' -X 'slender/internal/version.Commit=$COMMIT' -X 'slender/internal/version.BuildDate=$BUILD_DATE'" -o slender main.go
RUN upx -9 -o slender.minify slender && mv slender.minify slender

FROM alpine:3.21
COPY --from=Builder /app/slender /bin/slender
COPY assets /app/assets
COPY web /app/web

EXPOSE 8080
WORKDIR /app
ENTRYPOINT ["slender"]
