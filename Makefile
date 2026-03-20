tag:
	scripts/tag.sh

test-unit:
	go test ./... -coverprofile=unit_coverage.out 2>&1 | grep -E "^(ok|FAIL|---)"

test-integration:
	go test ./test/integration/... -coverpkg=./... -coverprofile=integration_coverage.out 2>&1 | grep -E "^(ok|FAIL|---)"

test: test-unit test-integration
	@echo "\nUnit test coverage:"
	@go tool cover -func=unit_coverage.out | grep total | awk '{print $3}' | sed 's/%//'
	@echo "\nIntegration test coverage:"
	@go tool cover -func=integration_coverage.out | grep total | awk '{print $3}' | sed 's/%//'