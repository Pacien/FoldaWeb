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
	"os"
)

func compiled(sourceDir, outputDir string, exts []string, saveAs string) {
	// remove previously compiled site
	err := os.RemoveAll(outputDir)
	if err != nil {
		fmt.Println(err)
		return
	}

	// compile everything
	wait.Add(2)
	go compile(sourceDir, make(map[string][]byte), sourceDir, outputDir, saveAs, exts, true)
	go copyFiles(sourceDir, sourceDir, outputDir, exts, true)

	// wait until all tasks are completed
	wait.Wait()
	fmt.Println("Compilation done.")
}
