FROM golang:alpine AS build-env
ADD . /src
RUN cd /src && go build cmd/inbox.go

FROM alpine
WORKDIR /app
COPY --from=build-env /src/inbox /app/
COPY --from=build-env /src/public /app/public
CMD ./inbox
