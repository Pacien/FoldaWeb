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
	"github.com/Pacien/fcmd"
	"github.com/drbawb/mustache"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"path"
	"strings"
	"sync"
)

var wait sync.WaitGroup

// Common templating

func isParsable(fileName string, exts []string) bool {
	for _, ext := range exts {
		if path.Ext(fileName) == ext {
			return true
		}
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

func parse(dirPath string, elements map[string][]byte, exts []string, overwrite bool) map[string][]byte {
	_, filesList := fcmd.Ls(dirPath)
	for _, fileName := range filesList {
		if isParsable(fileName, exts) && (overwrite || elements[fileName[:len(fileName)-len(path.Ext(fileName))]] == nil) {
			var err error
			elements[fileName[:len(fileName)-len(path.Ext(fileName))]], err = read(path.Join(dirPath, fileName))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return elements
}

func compile(dirPath string, elements map[string][]byte, sourceDir, outputDir, saveAs string, exts []string, recursive bool) {
	defer wait.Done()

	if strings.HasPrefix(dirPath, outputDir) {
		return
	}

	elements = parse(dirPath, elements, exts, true)

	if recursive {
		dirs, _ := fcmd.Ls(dirPath)
		for _, dir := range dirs {
			wait.Add(1)
			go compile(path.Join(dirPath, dir), elements, sourceDir, outputDir, saveAs, exts, recursive)
		}
	}

	pagePath := strings.TrimPrefix(dirPath, sourceDir)

	template := merge(elements)
	page := mustache.Render(string(template), makeContext(pagePath, sourceDir, exts))

	err := fcmd.WriteFile(path.Join(outputDir, pagePath, saveAs), []byte(page))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func copyFiles(dirPath, sourceDir, outputDir string, exts []string, recursive bool) {
	defer wait.Done()

	if strings.HasPrefix(dirPath, outputDir) {
		return
	}

	dirs, files := fcmd.Ls(dirPath)
	for _, file := range files {
		if !isParsable(file, exts) {
			err := fcmd.Cp(path.Join(dirPath, file), path.Join(outputDir, strings.TrimPrefix(dirPath, sourceDir), file))
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	if recursive {
		for _, dir := range dirs {
			wait.Add(1)
			go copyFiles(path.Join(dirPath, dir), sourceDir, outputDir, exts, recursive)
		}
	}
}
