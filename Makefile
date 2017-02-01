USER = ukhomeoffice-bot-test
REPO = testgithubintegration
COMMIT_ID = f591c33a1b9500d0721b6664cfb6033d47a00793

RULES_DIR = rules/gitrob.json
RULES_URL = https://raw.githubusercontent.com/michenriksen/gitrob/master/signatures.json

install:
	@go install -race .

lint:
	@golint  -set_exit_status ./...
	@go vet ./...
	@interfacer $(go list ./... | grep -v /vendor/)


rules:
	@curl -s $(RULES_URL) > $(RULES_DIR)/gitrob.json

# curl -s https://api.github.com/repos/ukhomeoffice-bot-test/testgithubintegration/commits/f591c33a1b9500d0721b6664cfb6033d47a00793 -H "Accept: application/vnd.github.VERSION.diff"
diff:
	@curl -s https://api.github.com/repos/$(USER)/$(REPO)/commits/$(COMMIT_ID) \
		-H "Accept: application/vnd.github.VERSION.diff"

run:
	@go build -race . && ./repo-security-scanner

test:
	@go test ./...

test-cover:
	@go test -cover ./...

test-race:
	@go test -race ./...


.PHONY: test run
