#
# (C) Copyright 2022 HP Development Company, L.P.
# Confidential computer software. Valid license from HP required for possession, use or copying.
# Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
# Computer Software Documentation, and Technical Data for Commercial Items are licensed
# to the U.S. Government under vendor's standard commercial license.
#

FROM golang:1.16.12-alpine3.15

LABEL maintainer="Runtime Team (cwp-gdrs-runtime@external.groups.hp.com)"
LABEL version=2.0

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the code
COPY ./code ./

# Download all dependencies.
RUN go get -d -v ./...
RUN go install -v ./...

# Build the application
RUN go build -o main .

# Expose port 9000
EXPOSE 9000

# Command to run the executable
CMD ["./main"]