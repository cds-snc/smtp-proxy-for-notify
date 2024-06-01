FROM alpine:3.6@sha256:66790a2b79e1ea3e1dabac43990c54aca5d1ddf268d9a5a0285e4167c8b24475 as alpine

RUN apk add -U --no-cache ca-certificates

FROM scratch

ARG SMTP_TLS_CERT_FILE
ARG SMTP_TLS_KEY_FILE
ARG SMTP_TLS_CERT_FILE_DESTINATION_DIR=/etc/ssl/certs/
ARG SMTP_TLS_KEY_FILE_DESTINATION_DIR=/etc/ssl/private/

# Copy the CA certs from the alpine image
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy the TLS cert and key files if they are defined. 
# Copying the LICENSE just allows for an optional copy of the certs if they are not defined.
COPY LICENSE ${SMTP_TLS_CERT_FILE}* ${SMTP_TLS_CERT_FILE_DESTINATION_DIR}
COPY LICENSE ${SMTP_TLS_KEY_FILE}* ${SMTP_TLS_KEY_FILE_DESTINATION_DIR}

# Copy the binary
COPY  ./release/latest/smtp-proxy-for-notify /smtp-proxy-for-notify

# Run the binary
ENTRYPOINT ["/smtp-proxy-for-notify"]