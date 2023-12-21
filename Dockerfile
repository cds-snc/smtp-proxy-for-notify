FROM alpine:3.19@sha256:51b67269f354137895d43f3b3d810bfacd3945438e94dc5ac55fdac340352f48 as alpine

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