PORT = 8080

GITHUB_SECRET = $(shell echo $(GITHUB_WEBHOOKSECRET))

USER = ukhomeoffice-bot-test
REPO = testgithubintegration
OFFENSES_X0 = 47797c0123bc0f5adfcae3d3467a2ed12e72b2cb
OFFENSES_X1 = f591c33a1b9500d0721b6664cfb6033d47a00793

FIXT_DIR = test/fixtures
RULES_FILE = $(FIXT_DIR)/rules/gitrob.json
DIFF_FILE = $(FIXT_DIR)/github_event_push.json
DIFF_FILE_OFFENSES = $(FIXT_DIR)/github_event_push_offenses.json
RULES_URL = https://raw.githubusercontent.com/michenriksen/gitrob/master/signatures.json

cli:
	@go install -race ./cmd/scanrepo

install: deps
	@go install -race --ldflags=\"-s\" .

deps: get-tools
	@trash

get-tools:
	@go get -u github.com/rancher/trash

lint:
	@golint
	@go vet

rules:
	@curl -s $(RULES_URL) > $(RULES_DIR)/gitrob.json

# curl -s https://api.github.com/repos/ukhomeoffice-bot-test/testgithubintegration/commits/f591c33a1b9500d0721b6664cfb6033d47a00793 -H "Accept: application/vnd.github.VERSION.diff"
diff-no-offenses:
	@curl -s \
		-H "Accept: application/vnd.github.VERSION.diff" \
		https://api.github.com/repos/$(USER)/$(REPO)/commits/$(OFFENSES_X0)

diff-offenses:
	@curl -s \
		-H "Accept: application/vnd.github.VERSION.diff" \
		https://api.github.com/repos/$(USER)/$(REPO)/commits/$(OFFENSES_X1)

struct:
	@gojson \
		-name githubResponseFull \
		-input test/fixtures/github_event_push.json

watch:
	@realize run

run:
	@go build -race . && ./repo-security-scanner

mac-diff-file:
	@cat $(DIFF_FILE) | openssl sha1 -hmac $(GITHUB_SECRET) | sed 's/^.* //'

mac-diff-file-offenses:
	@cat $(DIFF_FILE_OFFENSES) | openssl sha1 -hmac $(GITHUB_SECRET) | sed 's/^.* //'


test-run:
	@wget -O- \
		-X POST \
		--header="X-GitHub-Event: push" \
		--header="X-Hub-Signature: sha1=$(shell make mac-diff-file)" \
		--post-file "$(DIFF_FILE)" \
		http://localhost:$(PORT)/github

test-run-offenses:
	@wget -O- \
		-X POST \
		--header="X-GitHub-Event: push" \
		--header="X-Hub-Signature: sha1=$(shell make mac-diff-file-offenses)" \
		--post-file "$(DIFF_FILE_OFFENSES)" \
		http://localhost:$(PORT)/github

test-run-fail:
	@wget -O- \
		-X POST \
		--header="X-GitHub-Event: push" \
		--header="X-Hub-Signature: sha1=123456" \
		--post-file "$(DIFF_FILE)" \
		http://localhost:$(PORT)/github

test-run-dev:
	@curl \
		-X POST \
		-H "x-github-event: push" \
		--header="X-Hub-Signature: sha1=$(shell make mac-diff-file)" \
		-d @$(DIFF_FILE) \
		http://repo-security-scanner.notprod.homeoffice.gov.uk/github

test:
	@go test

.PHONY: test* run deps install lint rules struct diff
