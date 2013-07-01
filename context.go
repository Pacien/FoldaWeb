/*

	This file is part of CompileTree (https://github.com/Pacien/CompileTree)

	CompileTree is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	CompileTree is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with CompileTree. If not, see <http://www.gnu.org/licenses/>.

*/

package main

import (
	"github.com/Pacien/fcmd"
	"path"
	"strings"
)

type page struct {
	Title string
	URL   string
}

type context struct {
	path      string
	IsCurrent func(params []string, data string) string
	IsParent  func(params []string, data string) string
}

// Methods accessible in templates

func (c context) URL() string {
	p := strings.TrimPrefix(c.path, *settings.sourceDir)
	return path.Clean("/" + p)
}

func (c context) Title() string {
	_, t := path.Split(strings.TrimRight(c.URL(), "/"))
	return t
}

func (c context) SubPages() (subPages []page) {
	dirs, _ := fcmd.Ls(c.path)
	for _, dir := range dirs {
		var page page
		page.Title = dir
		page.URL = path.Join(c.URL(), dir)
		subPages = append(subPages, page)
	}
	return
}

func (c context) IsRoot() bool {
	if c.URL() == "/" {
		return true
	}
	return false
}

func (c context) isCurrent(pageTitle string) bool {
	if c.Title() == pageTitle {
		return true
	}
	return false
}

func (c context) isParent(pageTitle string) bool {
	for _, parent := range strings.Split(c.URL(), "/") {
		if parent == pageTitle {
			return true
		}
	}
	return false
}

func makeContext(pagePath, sourceDir, outputDir string, exts []string) (c context) {
	c.path = pagePath
	c.IsCurrent = func(params []string, data string) string {
		if c.isCurrent(strings.Join(params, " ")) {
			return data
		}
		return ""
	}
	c.IsParent = func(params []string, data string) string {
		if c.isParent(strings.Join(params, " ")) {
			return data
		}
		return ""
	}
	return
}
