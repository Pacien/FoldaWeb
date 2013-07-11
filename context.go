/*

	This file is part of FoldaWeb <https://github.com/Pacien/FoldaWeb>

	FoldaWeb is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	FoldaWeb is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with FoldaWeb. If not, see <http://www.gnu.org/licenses/>.

*/

package main

import (
	"github.com/Pacien/fcmd"
	"path"
	"strings"
)

type page struct {
	Title string
	Path  string
}

type context struct {
	filePath  string
	Path      string
	IsCurrent func(params []string, data string) string
	IsParent  func(params []string, data string) string
}

// Methods accessible in templates

func (c context) Title() string {
	_, t := path.Split(strings.TrimRight(c.Path, "/"))
	return t
}

func (c context) SubPages() (subPages []page) {
	dirs, _ := fcmd.Ls(c.filePath)
	for _, dir := range dirs {
		var page page
		page.Title = dir
		page.Path = path.Join(c.Path, dir)
		subPages = append(subPages, page)
	}
	return
}

func (c context) IsRoot() bool {
	if c.Path == "/" {
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
	for _, parent := range strings.Split(c.Path, "/") {
		if parent == pageTitle {
			return true
		}
	}
	return false
}

func makeContext(pagePath, sourceDir string, exts []string) (c context) {
	c.Path = path.Clean("/" + pagePath)
	c.filePath = path.Join(sourceDir, c.Path)
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
