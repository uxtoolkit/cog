// The UXToolkit Project
// Copyright (c) Wirecog, LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license, which can be found in the LICENSE file.

package cog

import (
	"errors"
	"reflect"
	"strings"

	"github.com/isomorphicgo/isokit"
	"github.com/uxtoolkit/reconcile"
	_ "golang.org/x/net/html"
	"honnef.co/go/js/dom"
)

type UXCog struct {
	Cog
	cogType         reflect.Type
	cogPrefixName   string
	cogPackagePath  string
	cogTemplatePath string
	templateSet     *isokit.TemplateSet
	Props           map[string]interface{}
	element         *dom.Element
	id              string
	hasBeenRendered bool
	parseTree       *reconcile.ParseTree
	cleanupFunc     func()
}

func (u *UXCog) getCogPrefixName() string {

	if u.cogType != nil {
		result := strings.Split(u.cogType.PkgPath(), `/`)
		return "cog:" + result[len(result)-1]
	} else {
		return ""
	}
}

func (u *UXCog) ID() string {
	return u.id
}

func (u *UXCog) SetID(id string) {
	u.id = id
}

func (u *UXCog) SetCleanupFunc(cleanupFunc func()) {
	u.cleanupFunc = cleanupFunc
}

func (u *UXCog) SetElement(element *dom.Element) {
	u.element = element
}

func (u *UXCog) Element() *dom.Element {
	return u.element
}

func (u *UXCog) CogInit(ts *isokit.TemplateSet) {
	u.hasBeenRendered = false
	u.Props = make(map[string]interface{})
	if ts != nil {
		u.templateSet = ts
	}
	u.cogTemplatePath = DefaultGoSourcePath + "/" + u.cogType.PkgPath() + "/" + DefaultTemplatesDirectoryName
	u.cogPrefixName = u.getCogPrefixName()

	if isokit.OperatingEnvironment() == isokit.ServerEnvironment {
		u.RegisterCogTemplates()
	}
}

func (u *UXCog) TemplateSet() *isokit.TemplateSet {
	return u.templateSet
}

func (u *UXCog) SetTemplateSet(ts *isokit.TemplateSet) {
	u.templateSet = ts
}

func (u *UXCog) CogType() reflect.Type {
	return u.cogType
}

func (u *UXCog) SetCogType(cogType reflect.Type) {
	u.cogType = cogType
}

func (u *UXCog) CogTemplatePath() string {
	return u.cogTemplatePath
}

func (u *UXCog) SetCogTemplatePath(path string) {
	u.cogTemplatePath = path
}

func (u *UXCog) RegisterCogTemplates() {
	u.templateSet.GatherCogTemplates(u.cogTemplatePath, u.cogPrefixName, ".tmpl")
}

func (u *UXCog) GetProps() map[string]interface{} {
	return u.Props

}

func (u *UXCog) SetProp(key string, value interface{}) {

	u.Props[key] = value
	if ReactivityEnabled == true && u.hasBeenRendered == true {
		u.Render()
	}

}

func (u *UXCog) BatchPropUpdate(props map[string]interface{}) {

	for k, v := range props {
		u.Props[k] = v
	}
	if ReactivityEnabled == true && u.hasBeenRendered == true {
		u.Render()
	}

}

func (u *UXCog) RenderCogTemplate() {

	var populateRenderedContent bool
	if u.hasBeenRendered == false {
		populateRenderedContent = true
	} else {
		populateRenderedContent = false
	}

	rp := isokit.RenderParams{Data: u.Props, Disposition: isokit.PlacementReplaceInnerContents, Element: *u.element, ShouldPopulateRenderedContent: populateRenderedContent}

	u.templateSet.Render(u.getCogPrefixName()+"/"+strings.Split(u.getCogPrefixName(), ":")[1], &rp)

	if u.hasBeenRendered == false {
		u.hasBeenRendered = true

		D := dom.GetWindow().Document()
		cogRoot := D.GetElementByID(u.id).FirstChild().(*dom.HTMLDivElement)
		contents := cogRoot.InnerHTML()
		parseTree, err := reconcile.NewParseTree([]byte(contents))
		if err != nil {
			println("Encountered an error: ", err)
		} else {
			u.parseTree = parseTree
		}

	}
}

func (u *UXCog) Render() error {
	document := dom.GetWindow().Document()
	e := document.GetElementByID(u.ID())

	if u.hasBeenRendered == true && e == nil {
		if u.cleanupFunc != nil {
			u.cleanupFunc()
			return nil
		}
	}

	if strings.ToLower(e.GetAttribute("data-component")) != "cog" {
		return errors.New("The cog container div must have a \"data-component\" attribute with a value specified as \"cog\".")
	}

	if u.hasBeenRendered == false {
		// Initial Render
		u.SetElement(&e)
		u.RenderCogTemplate()

		return nil
	} else if u.element != nil {
		// Re-render
		if VDOMEnabled == true {

			rp := isokit.RenderParams{Data: u.Props, Disposition: isokit.PlacementReplaceInnerContents, Element: *u.element, ShouldPopulateRenderedContent: true, ShouldSkipFinalRenderStep: true}
			u.templateSet.Render(u.getCogPrefixName()+"/"+strings.Split(u.getCogPrefixName(), ":")[1], &rp)

			D := dom.GetWindow().Document()
			cogRoot := D.GetElementByID(u.id).FirstChild().(*dom.HTMLDivElement)
			//contents := cogRoot.InnerHTML()
			newTree, err := reconcile.NewParseTree([]byte(rp.RenderedContent))

			if err != nil {
				println("Encountered an error: ", err)
			}

			changes, err := u.parseTree.Compare(newTree)
			if err != nil {
				println("Encountered an error: ", err)
			}
			if len(changes) > 0 {
				changes.ApplyChanges(cogRoot)
				u.parseTree = newTree
			}

		} else {
			u.RenderCogTemplate()

		}
		return nil
	}

	return nil

}
