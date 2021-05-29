# Build image
FROM alpine:latest AS build

# Build requirements
RUN apk add --no-cache ca-certificates

# Copy binary
COPY vcreport /vcreport

# ---

# Runtime image
FROM scratch
LABEL maintainer="Matthias Blaser <matthias.blaser@weekend4two.com>"

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /vcreport /vcreport

ENTRYPOINT ["/vcreport"]
CMD []
