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

// FoldaWeb, a "keep last legacy" website generator
package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/Pacien/fcmd"
	"github.com/drbawb/mustache"
	"github.com/frankbille/sanitize"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"path"
	"strings"
	"sync"
)

type generator struct {
	// parameters
	sourceDir, outputDir string
	startWith, saveAs    string
	wordSeparator        string
	skipPrefix           string
	parsableExts         []string

	// go routine sync
	tasks sync.WaitGroup
}

type page struct {
	// properties accessible in templates
	Title               string
	AbsPath, Path       string
	IsRoot              bool
	IsCurrent, IsParent func(params []string, data string) string

	// properties used for page generation
	dirPath string
	parts   parts
	body    []byte
}

type parts map[string][]byte

func (p parts) clone() parts {
	c := parts{}
	for k, v := range p {
		c[k] = v
	}
	return c
}

// Creates an initiated generator
func newGenerator() (g generator) {

	// Read the command line arguments
	flag.StringVar(&g.sourceDir, "sourceDir", "./source", "Path to the source directory.")
	flag.StringVar(&g.outputDir, "outputDir", "./out", "Path to the output directory.")
	flag.StringVar(&g.startWith, "startWith", "index", "Name without extension of the first file that will by parsed.")
	flag.StringVar(&g.saveAs, "saveAs", "index.html", "Save compiled files as named.")
	flag.StringVar(&g.wordSeparator, "wordSeparator", "-", "Word separator used to replace spaces in URLs.")
	flag.StringVar(&g.skipPrefix, "skipPrefix", "_", "Folders with this prefix will be hidden in the output.")
	var parsableExts string
	flag.StringVar(&parsableExts, "parsableExts", "html, txt, md", "Parsable file extensions separated by commas.")

	flag.Parse()

	g.sourceDir = path.Clean(g.sourceDir)
	g.outputDir = path.Clean(g.outputDir)
	for _, ext := range strings.Split(parsableExts, ",") {
		g.parsableExts = append(g.parsableExts, "."+strings.Trim(ext, ". "))
	}

	return

}

func (g *generator) sanitizePath(filePath string) string {
	sanitizedFilePath := strings.Replace(filePath, " ", g.wordSeparator, -1)
	return sanitize.Path(sanitizedFilePath)
}

func (g *generator) sourcePath(filePath string) string {
	return path.Join(g.sourceDir, filePath)
}

func (g *generator) outputPath(filePath string) string {
	pathElements := strings.Split(filePath, "/")
	var finalFilePath string
	for _, element := range pathElements {
		if !strings.HasPrefix(element, g.skipPrefix) {
			finalFilePath = path.Join(finalFilePath, element)
		}
	}
	return path.Join(g.outputDir, g.sanitizePath(finalFilePath))
}

func (g *generator) isFileParsable(fileName string) bool {
	for _, ext := range g.parsableExts {
		if path.Ext(fileName) == ext {
			return true
		}
	}
	return false
}

func (g *generator) copyFile(filePath string) {
	defer g.tasks.Done()
	fmt.Println("Copying: " + filePath)
	err := fcmd.Cp(g.sourcePath(filePath), g.outputPath(filePath))
	if err != nil {
		fmt.Println(err)
	}
}

func (g *generator) parseFile(filePath string) []byte {
	fileBody, err := ioutil.ReadFile(g.sourcePath(filePath))
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if path.Ext(filePath) == ".md" {
		fileBody = blackfriday.MarkdownCommon(fileBody)
	}
	return fileBody
}

func (g *generator) mergeParts(parts parts) []byte {
	merged := parts[g.startWith]
	for pass := 0; bytes.Contains(merged, []byte("{{> ")) && pass < 1000; pass++ {
		for partName, partBody := range parts {
			merged = bytes.Replace(merged, []byte("{{> "+partName+"}}"), partBody, -1)
		}
	}
	return merged
}

func (g *generator) contextualize(page page) page {
	_, page.Title = path.Split(page.dirPath)
	if page.dirPath == "" {
		page.IsRoot = true
		page.AbsPath, page.Path = "/", "/"
	} else {
		page.AbsPath = g.sanitizePath("/" + page.dirPath)
		_, page.Path = path.Split(page.AbsPath)
	}

	page.IsCurrent = func(params []string, data string) string {
		if page.Path == path.Clean(params[0]) {
			return data
		}
		return ""
	}

	page.IsParent = func(params []string, data string) string {
		if strings.Contains(page.AbsPath, path.Clean(params[0])) {
			return data
		}
		return ""
	}

	return page
}

func (g *generator) generate(page page) {
	defer g.tasks.Done()

	dirs, files := fcmd.Ls(g.sourcePath(page.dirPath))

	// Parse or copy files in the current directory
	containsParsableFiles := false
	for _, file := range files {
		filePath := path.Join(page.dirPath, file)
		if g.isFileParsable(file) {
			containsParsableFiles = true
			page.parts[file[:len(file)-len(path.Ext(file))]] = g.parseFile(filePath)
		} else {
			g.tasks.Add(1)
			go g.copyFile(filePath)
		}
	}

	// Generate subpages in surdirectories
	for _, dir := range dirs {
		subPage := page
		subPage.dirPath = path.Join(page.dirPath, dir)
		subPage.parts = page.parts.clone()
		g.tasks.Add(1)
		go g.generate(subPage)
	}

	// Generate the page at the current directory
	if _, currentDir := path.Split(page.dirPath); containsParsableFiles && !strings.HasPrefix(currentDir, g.skipPrefix) {
		fmt.Println("Rendering: " + page.dirPath)
		page.body = []byte(mustache.Render(string(g.mergeParts(page.parts)), g.contextualize(page)))
		err := fcmd.WriteFile(g.outputPath(path.Join(page.dirPath, g.saveAs)), page.body)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func main() {
	fmt.Println("FoldaWeb <https://github.com/Pacien/FoldaWeb>")

	g := newGenerator()

	// Remove previously generated site
	err := fcmd.Rm(g.outputDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Generate everything
	page := page{}
	page.parts = make(parts)
	g.tasks.Add(1)
	go g.generate(page)

	// Wait until all tasks are completed
	g.tasks.Wait()
	fmt.Println("Done.")
}
