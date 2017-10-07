// The UXToolkit Project
// Copyright (c) Wirecog, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license, which can be found in the LICENSE file.

package cog

import (
	"os"
)

var DefaultTemplatesDirectoryName string
var DefaultGoSourcePath string
var TemplateFileExtension = ".tmpl"
var ReactivityEnabled = true
var VDOMEnabled = true

type Cog interface {
	Render() error
}

func init() {

	DefaultTemplatesDirectoryName = "templates"
	DefaultGoSourcePath = os.Getenv("GOPATH") + "/src"

}
