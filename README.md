# SMTP Proxy for Notify

This is a simple SMTP proxy for Notify that bridges the gap between the SMTP protocol and the Notify API. The proxy listens for SMTP connections and converts the SMTP message into a Notify API call. The Notify API then sends the message to the recipient. Attachments are supported.

## Usage

### Requirements

- A Notify account
- A Notify API key
- A Notify template ID

Your template must have a `subject` and `body` field. For example, the template should look something like this:

![An image showing a Notify template](https://github.com/cds-snc/smtp-proxy-for-notify/assets/867334/a868d28b-f4fb-4069-95ed-6fb11bbf5aae)

### Environment variables

| Variable | Description | Required | Default |
| --- | --- | --- | --- |
| NOTIFY_APIKEY | Your Notify API key | Yes | |
| NOTIFY_HOSTNAME | The hostname of the Notify API | No | https://api.notification.canada.ca |
| NOTIFY_TEMPLATE_ID | Your Notify template ID | Yes |  |
| SMTP_TLS_CERT_FILE | Path to your TLS certificate file | No | |
| SMTP_TLS_KEY_FILE | Path to your TLS key file | No | |
| SMTP_USE_TLS | Whether to use TLS or not | No | false |
| SMTP_HOSTNAME | The hostname to listen on | No | localhost |
| SMTP_PORT | The port to listen on | No | 1025 |
| SMTP_USERNAME | The username to use for authentication | Yes |
| SMTP_PASSWORD | The password to use for authentication | Yes |

### Running

#### Locally

You can run the proxy locally using the following command as long as you have all the environment variables set:

```bash
./release/latest/smtp-proxy-for-notify
```

#### Docker

The proxy can also be run using Docker. You can build the image using the Dockerfile in this repository. However, you should provide your own TLS certificate and key files. You can do this by building the image with the `SMTP_TLS_CERT_FILE` and `SMTP_TLS_KEY_FILE` build arguments. For example:

```bash
docker build --build-arg SMTP_TLS_CERT_FILE="example_certs/server.crt" --build-arg=SMTP_TLS_KEY_FILE="example_certs/server.key" -t smtp .
```

The alternative is to mount the certificate and key files into the container at runtime.

To run the container, you can use the following command assuming you built it with the Dockerfile in this repository:

```bash
docker run \
    -e NOTIFY_APIKEY=gcntfy-test-00000000-0000-4000-8000-000000000000-00000000-0000-4000-8000-0000000000000 \
    -e NOTIFY_TEMPLATE_ID=00000000-0000-4000-8000-0000000000008 \
    -e SMTP_TLS_CERT_FILE=/etc/ssl/certs/server.crt \
    -e SMTP_TLS_KEY_FILE=/etc/ssl/private/server.key \
    -e SMTP_USE_TLS=true \
    -e SMTP_HOSTNAME=0.0.0.0 \
    -e SMTP_USERNAME=username \
    -e SMTP_PASSWORD=longpasswordgo \
    -p 1025:1025 \
    smtp
```

# License

MIT License