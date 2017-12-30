// The UXToolkit Project
// Copyright (c) Wirecog, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license, which can be found in the LICENSE file.

package cog

import (
	"go/build"
	"os"
)

var DefaultTemplatesDirectoryName string
var DefaultGoSourcePath string
var TemplateFileExtension = ".tmpl"
var ReactivityEnabled = true
var VDOMEnabled = true

type Cog interface {
	Render() error
	Start() error
}

func init() {

	var gopath string
	gp := os.Getenv("GOPATH")
	if gp != "" {
		gopath = gp
	} else {
		gopath = build.Default.GOPATH
	}
	DefaultTemplatesDirectoryName = "templates"
	DefaultGoSourcePath = gopath + "/src"

}
