.PHONY: test-backend build-backend run-backend smoke-backend deploy-local

test-backend:
	cd backend && go test ./...

build-backend:
	cd backend && go build ./...

run-backend:
	cd backend && go run ./cmd/api

smoke-backend:
	bash deploy/scripts/smoke_backend.sh

deploy-local:
	bash deploy/scripts/deploy_local.sh
