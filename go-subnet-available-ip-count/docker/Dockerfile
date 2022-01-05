#
# (C) Copyright 2022 HP Development Company, L.P.
# Confidential computer software. Valid license from HP required for possession, use or copying.
# Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
# Computer Software Documentation, and Technical Data for Commercial Items are licensed
# to the U.S. Government under vendor's standard commercial license.
#

FROM golang:1.16.12-alpine3.15

LABEL maintainer="Runtime Team (cwp-gdrs-runtime@external.groups.hp.com)"
LABEL version=1.0

WORKDIR /app

COPY ./code ./

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -o main .

CMD ["./main"]