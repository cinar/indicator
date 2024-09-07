// Copyright (c) 2021-2024 Onur Cinar.
// The source code is provided under GNU AGPLv3 License.
// https://github.com/cinar/indicator

package asset

import (
	"fmt"
)

const (
	// InMemoryRepositoryBuilderName is the name for the in memory repository builder.
	InMemoryRepositoryBuilderName = "memory"

	// FileSystemRepositoryBuilderName is the name for the file system repository builder.
	FileSystemRepositoryBuilderName = "filesystem"

	// TiingoRepositoryBuilderName is the name of the Tiingo repository builder.
	TiingoRepositoryBuilderName = "tiingo"
)

// RepositoryBuilderFunc defines a function to build a new repository using the given configuration parameter.
type RepositoryBuilderFunc func(config string) (Repository, error)

// repositoryBuilders provides mapping for the repository builders.
var repositoryBuilders = map[string]RepositoryBuilderFunc{
	InMemoryRepositoryBuilderName:   inMemoryRepositoryBuilder,
	FileSystemRepositoryBuilderName: fileSystemRepositoryBuilder,
	TiingoRepositoryBuilderName:     tiingoRepositoryBuilder,
}

// RegisterRepositoryBuilder registers the given builder.
func RegisterRepositoryBuilder(name string, builder RepositoryBuilderFunc) {
	repositoryBuilders[name] = builder
}

// NewRepository builds a new repository by the given name type and the configuration.
func NewRepository(name, config string) (Repository, error) {
	builder, ok := repositoryBuilders[name]
	if !ok {
		return nil, fmt.Errorf("unknown repository: %s", name)
	}

	return builder(config)
}

// inMemoryRepositoryBuilder builds a new in memory repository instance.
func inMemoryRepositoryBuilder(_ string) (Repository, error) {
	return NewInMemoryRepository(), nil
}

// fileSystemRepositoryBuilder builds a new file system repository instance.
func fileSystemRepositoryBuilder(config string) (Repository, error) {
	return NewFileSystemRepository(config), nil
}

// tiingoRepositoryBuilder builds a new Tiingo repository instance.
func tiingoRepositoryBuilder(config string) (Repository, error) {
	return NewTiingoRepository(config), nil
}
