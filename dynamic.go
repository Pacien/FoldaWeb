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
	"github.com/drbawb/mustache"
	"net/http"
	"path"
	"strings"
)

func handle(w http.ResponseWriter, r *http.Request) {
	// serve static files
	if !(path.Ext(r.URL.Path) == "" || isParsable(path.Ext(r.URL.Path), settings.exts)) {
		http.ServeFile(w, r, path.Join(*settings.sourceDir, r.URL.Path))
		return
	}

	// get the list of dirs to parse
	request := strings.TrimSuffix(r.URL.Path, "/")
	dirs := strings.Split(request, "/")
	for i, dir := range dirs {
		if i != 0 {
			dirs[i] = path.Join(dirs[i-1], dir)
		}
	}

	// parse these dirs
	elements := make(map[string][]byte)
	for i := len(dirs) - 1; i >= 0; i-- /*_, dir := range reverse dirs*/ {
		parse(path.Join(*settings.sourceDir, dirs[i]), elements, settings.exts, false)
	}

	// render the page
	template := merge(elements)
	page := mustache.Render(string(template), makeContext(r.URL.Path, *settings.sourceDir, settings.exts))

	// serve the page
	_, err := w.Write([]byte(page))
	if err != nil {
		fmt.Println(err)
		return
	}
}

func dynamic(port string) {
	fmt.Println("Listening on: localhost:" + port)
	http.HandleFunc("/", handle)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err)
	}
}
