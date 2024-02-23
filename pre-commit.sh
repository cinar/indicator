#!/bin/bash

(command -v revive) || go install github.com/mgechev/revive@v1.3.4
(command -v staticcheck) || go install honnef.co/go/tools/cmd/staticcheck@latest
(command -v gomarkdoc) || go install github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest
(command -v gosec) || go install github.com/securego/gosec/v2/cmd/gosec@latest
(command -v markdownfmt) || go install github.com/shurcooL/markdownfmt@latest

PACKAGES=(asset helper momentum strategy trend volatility volume)
PACKAGE_FILES=("${PACKAGES[@]/#/./}")
PACKAGE_FILES=("${PACKAGE_FILES[@]/%//...}")

go fmt ./...
go fix ./...
go vet ./...
go test -cover "${PACKAGE_FILES[@]}"
gosec ./...

revive -config=revive.toml ./...
staticcheck ./...

for package in "${PACKAGES[@]}";
do
    echo Package "$package"
    gomarkdoc --repository.default-branch v2 --output "$package"/README.md ./"$package"
done

# markdownfmt -w README.md

