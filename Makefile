PORT = 8080

USER = ukhomeoffice-bot-test
REPO = testgithubintegration
OFFENSES_X0 = 47797c0123bc0f5adfcae3d3467a2ed12e72b2cb
OFFENSES_X1 = f591c33a1b9500d0721b6664cfb6033d47a00793

FIXT_DIR = test/fixtures
RULES_FILE = $(FIXT_DIR)/rules/gitrob.json
DIFF_FILE = $(FIXT_DIR)/github_event_push.json
RULES_URL = https://raw.githubusercontent.com/michenriksen/gitrob/master/signatures.json

install: deps
	@go install -race --ldflags=\"-s\" .

deps: get-tools
	@trash

get-tools:
	@go get -u \
			github.com/rancher/trash

# @interfacer $(go list ./... | grep -v /vendor/)
lint:
	@golint
	@go vet

rules:
	@curl -s $(RULES_URL) > $(RULES_DIR)/gitrob.json

# curl -s https://api.github.com/repos/ukhomeoffice-bot-test/testgithubintegration/commits/f591c33a1b9500d0721b6664cfb6033d47a00793 -H "Accept: application/vnd.github.VERSION.diff"
diff-no-offenses:
	@curl -s \
		https://api.github.com/repos/$(USER)/$(REPO)/commits/$(OFFENSES_X0) \
		-H "Accept: application/vnd.github.VERSION.diff"

diff-offenses:
	@curl -s \
		https://api.github.com/repos/$(USER)/$(REPO)/commits/$(OFFENSES_X1) \
		-H "Accept: application/vnd.github.VERSION.diff"

struct:
	@gojson \
		-name githubResponseFull \
		-input test/fixtures/github_event_push.json

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
	@go test

test-cover:
	@go test -cover

test-race:
	@go test -race


.PHONY: test* run deps install lint rules struct diff
