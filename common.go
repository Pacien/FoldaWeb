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
	"bytes"
	"fmt"
	"github.com/hoisie/mustache"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"path"
	"strings"
	"sync"
)

var wait sync.WaitGroup

// Common templating

func isParsable(fileName string) bool {
	switch path.Ext(fileName) {
	case ".md", ".html", ".txt":
		return true
	}
	return false
}

func read(fileName string) ([]byte, error) {
	fileBody, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	if path.Ext(fileName) == ".md" {
		fileBody = blackfriday.MarkdownCommon(fileBody)
	}
	return fileBody, nil
}

func merge(files map[string][]byte) (merged []byte) {
	merged = files["index"]
	for pass := 0; bytes.Contains(merged, []byte("{{> ")) && pass < 4000; pass++ {
		for fileName, fileBody := range files {
			merged = bytes.Replace(merged, []byte("{{> "+fileName+"}}"), fileBody, -1)
		}
	}
	return
}

// COMPILED and INTERACTIVE modes

// render and write everything inside

func parse(dirPath string, elements map[string][]byte, overwrite bool) map[string][]byte {
	_, filesList := ls(dirPath)
	for _, fileName := range filesList {
		if isParsable(fileName) && (overwrite || elements[fileName[:len(fileName)-len(path.Ext(fileName))]] == nil) {
			var err error
			elements[fileName[:len(fileName)-len(path.Ext(fileName))]], err = read(path.Join(dirPath, fileName))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return elements
}

func compile(dirPath string, elements map[string][]byte, sourceDir, outputDir string, recursive bool) {
	wait.Add(1)
	defer wait.Done()

	if strings.HasPrefix(dirPath, outputDir) {
		return
	}

	elements = parse(dirPath, elements, true)

	if recursive {
		dirs, _ := ls(dirPath)
		for _, dir := range dirs {
			go compile(path.Join(dirPath, dir), elements, sourceDir, outputDir, recursive)
		}
	}

	template := merge(elements)
	page := mustache.Render(string(template), nil /* TODO: generate contextual variables */)

	err := writeFile(path.Join(outputDir, strings.TrimPrefix(dirPath, sourceDir), "index.html"), []byte(page))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func copyFiles(dirPath, sourceDir, outputDir string, recursive bool) {
	wait.Add(1)
	defer wait.Done()

	if strings.HasPrefix(dirPath, outputDir) {
		return
	}

	dirs, files := ls(dirPath)
	for _, file := range files {
		if !isParsable(file) {
			err := cp(path.Join(dirPath, file), path.Join(outputDir, strings.TrimPrefix(dirPath, sourceDir), file))
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	if recursive {
		for _, dir := range dirs {
			go copyFiles(path.Join(dirPath, dir), sourceDir, outputDir, recursive)
		}
	}
}
