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
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Filesystem utils

func isDir(dirPath string) bool {
	stat, err := os.Stat(dirPath)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

func isHidden(fileName string) bool {
	return strings.HasPrefix(fileName, ".")
}

func ls(path string) (dirs []string, files []string) {
	content, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, element := range content {
		if isHidden(element.Name()) {
			continue
		}
		if element.IsDir() {
			dirs = append(dirs, element.Name())
		} else {
			files = append(files, element.Name())
		}
	}
	return
}

func explore(dirPath string) (paths []string) {
	dirs, _ := ls(dirPath)
	for _, dir := range dirs {
		sourceDir := path.Join(dirPath, dir)
		paths = append(paths, sourceDir)
		subDirs := explore(sourceDir)
		for _, subDir := range subDirs {
			paths = append(paths, subDir)
		}
	}
	return
}

func cp(source, target string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	dir, _ := path.Split(target)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}

	targetFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, sourceFile)
	return err
}

func writeFile(target string, body []byte) error {
	dir, _ := path.Split(target)
	err := os.MkdirAll(dir, 0777)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(target, body, 0777)
	return err
}
