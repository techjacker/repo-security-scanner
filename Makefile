PORT = 8080

USER = ukhomeoffice-bot-test
REPO = testgithubintegration
COMMIT_ID = f591c33a1b9500d0721b6664cfb6033d47a00793

FIXT_DIR = test/fixtures
RULES_FILE = $(FIXT_DIR)/rules/gitrob.json
DIFF_FILE = $(FIXT_DIR)/github_diff_response.json
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
	@curl -s \
		https://api.github.com/repos/$(USER)/$(REPO)/commits/$(COMMIT_ID) \
		-H "Accept: application/vnd.github.VERSION.diff"

struct:
	@gojson \
		-name githubResponseFull \
		-input test/fixtures/github_diff_response.json

watch:
	@realize run

run:
	@go build -race . && ./repo-security-scanner

test-run:
	@curl \
		-X POST \
		-d @$(DIFF_FILE) \
		http://localhost:$(PORT)/github

test:
	@go test ./...

test-cover:
	@go test -cover ./...

test-race:
	@go test -race ./...


.PHONY: test* run
