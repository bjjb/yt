test:
	@go test -coverprofile=coverage.out | grep 'coverage: 100.0% of statements'
