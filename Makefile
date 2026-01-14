.EXPORT_ALL_VARIABLES:
NOTIFY_HOSTNAME=https://api.staging.notification.cdssandbox.xyz
NOTIFY_APIKEY=gcntfy-test-00000000-0000-4000-8000-000000000000-00000000-0000-4000-8000-0000000000000
NOTIFY_TEMPLATE_ID=00000000-0000-4000-8000-000000000000

SMTP_USE_TLS=true
SMTP_TLS_CERT_FILE=./example_certs/server.crt
SMTP_TLS_KEY_FILE=./example_certs/server.key
SMTP_HOSTNAME=localhost
SMTP_PORT=1025
SMTP_USERNAME=username
SMTP_PASSWORD=longpasswordgo

TEST_SENDER=author@localhost
TEST_RECIPIENT=max.neuvians+staging-notify-1@cds-snc.ca

.PHONY: dev generate-keys release release-test script-test-python test
.DEFAULT_GOAL := release

dev:
	@echo "Starting dev server..."
	@go run .

generate-keys:
	@cd example_certs && \
	rm -f server.key server.crt && \
	echo "Generating keys..." && \
	openssl genrsa -out server.key 2048 && \
	openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650 -subj /CN=localhost/O=smtp-proxy-for-notify/C=CA

release:
	@mkdir -p release/latest
	@docker build -t smtp-proxy-for-notify-build -f Dockerfile.build .
	@docker create -ti --name smtp-proxy-for-notify-build smtp-proxy-for-notify-build bash 
	@docker cp smtp-proxy-for-notify-build:/smtp-proxy-for-notify release/latest/smtp-proxy-for-notify
	@docker rm -f smtp-proxy-for-notify-build

release-test:
	@mkdir -p release/latest
	@docker build -t smtp-proxy-for-notify-build -f Dockerfile.build .
	@docker create -ti --name smtp-proxy-for-notify-build smtp-proxy-for-notify-build bash 
	@docker cp smtp-proxy-for-notify-build:/smtp-proxy-for-notify release/latest/smtp-proxy-for-notify-test
	@docker rm -f smtp-proxy-for-notify-build

script-test-python:
	@python3 ./bin/test_send.py

test:
	@go test -cover ./...