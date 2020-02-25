// +build tools

// Place any runtime dependencies as imports in this file.
// Go modules will be forced to download and install them.
package tools

import (
	_ "github.com/golang/protobuf/proto"
	_ "k8s.io/code-generator"
)
