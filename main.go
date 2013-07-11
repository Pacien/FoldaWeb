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
	"flag"
	"fmt"
	"github.com/Pacien/fcmd"
	"strings"
)

var settings struct {
	mode      *string // compiled, interactive or dynamic
	sourceDir *string
	outputDir *string // for compiled site
	port      *string // for the integrated web server (dynamic mode only)
	exts      []string
	saveAs    *string
}

func init() {
	fcmd.DefaultPerm = 0755 // -rwxr-xr-x

	// read settings
	settings.mode = flag.String("mode", "compiled", "compiled|interactive|dynamic")
	settings.sourceDir = flag.String("source", "./source", "Path to sources directory.")
	settings.outputDir = flag.String("output", "./out", "[compiled mode] Path to output directory.")
	settings.port = flag.String("port", "8080", "[dynamic mode] Port to listen.")
	exts := flag.String("exts", "html, txt, md", "List parsable file extensions. Separated by commas.")
	settings.saveAs = flag.String("saveAs", "index.html", "[compiled and interactive modes] Save compiled files as named.")
	flag.Parse()
	settings.exts = strings.Split(*exts, ",")
	for i, ext := range settings.exts {
		settings.exts[i] = "." + strings.Trim(ext, ". ")
	}
	*settings.sourceDir = strings.TrimPrefix(*settings.sourceDir, "./")
	*settings.outputDir = strings.TrimPrefix(*settings.outputDir, "./")
}

func main() {
	fmt.Println("FoldaWeb <https://github.com/Pacien/FoldaWeb>")
	fmt.Println("Mode: " + *settings.mode)
	fmt.Println("Source: " + *settings.sourceDir)
	fmt.Println("Output: " + *settings.outputDir)
	fmt.Println("====================")

	switch *settings.mode {
	case "compiled":
		compiled(*settings.sourceDir, *settings.outputDir, settings.exts, *settings.saveAs)
	case "interactive":
		interactive(*settings.sourceDir, *settings.outputDir, settings.exts, *settings.saveAs)
	case "dynamic":
		dynamic(*settings.port)
	default:
		fmt.Println("Invalid mode.")
	}
}
