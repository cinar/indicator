#!/bin/bash

(command -v revive) || go install github.com/mgechev/revive@v1.3.4
(command -v staticcheck) || go install honnef.co/go/tools/cmd/staticcheck@latest
(command -v gomarkdoc) || go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest
(command -v gosec) || go install github.com/securego/gosec/v2/cmd/gosec@latest
(command -v markdownfmt) || go install github.com/shurcooL/markdownfmt@latest

PACKAGES=(
	asset
	helper
	momentum
	strategy
	strategy/trend
	strategy/volatility
	strategy/momentum
	strategy/compound
	trend
	volatility
	volume
)

PACKAGE_FILES=("${PACKAGES[@]/#/./}")
PACKAGE_FILES=("${PACKAGE_FILES[@]/%//...}")

go fmt ./...
go fix ./...
go vet ./...
go test -cover "${PACKAGE_FILES[@]}"
gosec ./...

revive -config=revive.toml ./...
staticcheck ./...
gomarkdoc ./...

markdownfmt -w README.md

