NAME := apex2java
SRCS := $(shell find . -type d -name vendor -prune -o -type f -name "*.go" -print)
VERSION := 0.1.0
LDFLAGS := -ldflags="-s -w -X \"main.Version=$(VERSION)\""
DIST_DIRS := find * -type d -exec
ANTLR := java -Xmx500M -cp "/usr/local/lib/antlr-4.7.1-complete.jar:$(CLASSPATH)" org.antlr.v4.Tool

.PHONY: run
run: format
	@go run . run -a "Foo#action" -f example.cls

.PHONY: run/format
run/format:
	@go run . format -f example.cls

.PHONY: test
test: format
	@go test ./...

.PHONY: build
build: format
	@go build

.PHONY: format
format: import
	@gofmt -w .

.PHONY: import
import:
ifneq ($(shell command -v goimports 2> /dev/null),)
	@goimports -w ./main.go
endif

.PHONY: cross-build
cross-build: deps
	-@goimports -w $(SRCS)
	@gofmt -w $(SRCS)
	@for os in darwin linux windows; do \
	    for arch in amd64; do \
	        CC=clang GOOS=$$os GOARCH=$$arch CGO_ENABLED=1 go build -a -tags netgo \
	        -installsuffix netgo $(LDFLAGS) -o dist/$$os-$$arch/$(NAME); \
	    done; \
	done

.PHONY: dist
dist:
	@cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) tar zcf $(NAME)-$(VERSION)-{}.tar.gz {} \;

.PHONY: deploy
deploy:
	git push heroku master -f

.PHONY: tag
tag:
	git tag v$(VERSION) -f
	git push origin v$(VERSION) -f

