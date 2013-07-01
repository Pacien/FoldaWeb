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
	"fmt"
	"github.com/Pacien/fcmd"
	"github.com/howeyc/fsnotify"
	"os"
	"path"
	"strings"
	"time"
)

func watch(dirPath string, watcher *fsnotify.Watcher) *fsnotify.Watcher {
	watcher.Watch(dirPath)
	dirs, _ := fcmd.Explore(dirPath)
	for _, dir := range dirs {
		if !strings.HasPrefix(dir, *settings.outputDir) {
			err := watcher.Watch(dir)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
	return watcher
}

func parseParents(dir, sourceDir string, exts []string) map[string][]byte {
	dirs := strings.Split(strings.TrimPrefix(dir, sourceDir), "/")
	elements := make(map[string][]byte)
	for _, dir := range dirs {
		elements = parse(path.Join(sourceDir, dir), elements, exts, false)
	}
	return elements
}

func interactive(sourceDir, outputDir string, exts []string, saveAs string) {

	// compile the whole site
	compiled(sourceDir, outputDir, exts, saveAs)

	// watch the source dir
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println(err)
	}
	defer watcher.Close()
	watcher = watch(sourceDir, watcher)

	for {
		select {
		case ev := <-watcher.Event:
			fmt.Println(ev)

			// ignore hidden files
			if fcmd.IsHidden(ev.Name) {
				break
			}

			// manage watchers
			if ev.IsDelete() || ev.IsRename() {
				err = watcher.RemoveWatch(ev.Name)
				if err != nil {
					fmt.Println(err)
					return
				}
			} else if ev.IsCreate() && fcmd.IsDir(ev.Name) {
				watcher = watch(ev.Name, watcher)
			}

			dir, _ := path.Split(ev.Name)

			// remove previously compiled files
			if ev.IsDelete() || ev.IsRename() || ev.IsModify() {
				var err error
				if fcmd.IsDir(ev.Name) || !isParsable(ev.Name, exts) {
					err = os.RemoveAll(path.Join(outputDir, strings.TrimPrefix(ev.Name, sourceDir)))
				} else {
					err = os.RemoveAll(path.Join(outputDir, strings.TrimPrefix(dir, sourceDir)))
				}
				if err != nil {
					fmt.Println(err)
					return
				}
			}

			// recompile changed files
			if ev.IsCreate() || ev.IsModify() {
				if fcmd.IsDir(ev.Name) {
					elements := parseParents(ev.Name, sourceDir, exts)
					dirPath := path.Join(sourceDir, strings.TrimPrefix(ev.Name, sourceDir))
					go compile(dirPath, elements, sourceDir, outputDir, saveAs, exts, true)
					go copyFiles(dirPath, sourceDir, outputDir, exts, true)
				} else {
					dirPath := path.Join(sourceDir, strings.TrimPrefix(dir, sourceDir))
					if isParsable(path.Ext(ev.Name), exts) {
						elements := parseParents(dir, sourceDir, exts)
						go compile(dirPath, elements, sourceDir, outputDir, saveAs, exts, true)
					}
					go copyFiles(dirPath, sourceDir, outputDir, exts, false)
				}
			}

			// sleep some milliseconds to prevent early exit
			time.Sleep(time.Millisecond * 100)

			// wait until all tasks are completed
			wait.Wait()

		case err := <-watcher.Error:
			fmt.Println(err)
		}
	}
}
