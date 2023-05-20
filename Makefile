export COVERAGE_PACKAGES=http service

coverage:
	echo "mode: count" > coverage-all.out
	$(foreach pkg, $(COVERAGE_PACKAGES),\
		go test -cover -coverprofile=coverage.out -covermode=count ./pkg/$(pkg)/...;\
		tail -n +2 coverage.out >> coverage-all.out;)
	cat coverage-all.out | grep -v "mock_" > coverage-all.out
	cat coverage-all.out | grep -v "router.go" > coverage-all.out
	go tool cover -func=coverage-all.out