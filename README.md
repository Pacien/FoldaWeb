FoldaWeb
========

### Description

FoldaWeb is a "keep last legacy" website generator: the program generates a page from parts of pages that are presents inside a directory and directly inside its parents. If multiple parts have the same name, it uses the last one (the one located the deepest in the directories before the current included).

This behaviour makes particularly easy to create well-organized websites with many subpages of different types with associated layouts for each.

___

### Features

- Unique "keep last legacy" generation (no pun intended)
- Mustache templating: FoldaWeb uses [Mustache](http://mustache.github.io/mustache.5.html) as template engine and adds several handy contextual variables
- Markdown compatible: pages can be written using the [Markdown syntax](http://daringfireball.net/projects/markdown/syntax)

Moreover, because FoldaWeb generates static files, generated websites are:

- **Portable**: any host and web server software can serve flat files.
- **Fast**: no server-side scripting is required everytime someone loads a page
- **Secure**: no CMS security flaws

___

### Example

[Multiverse Inc. Global Website](http://multiverse.pacien.net) is an example of website generated using FoldaWeb.

Its sources are available on GitHub at [Pacien/FoldaWeb-example](https://github.com/Pacien/FoldaWeb-example)

___

### Usage

Simply put the binary inside a directory containing a `source` folder with the website's sources inside and run the program (simply open the executable). Another folder named `out` containing the generated website will be created instantly.

You can also pass custom settings via command line arguments:

    -sourceDir="./source": Path to the source directory.
    -outputDir="./out": Path to the output directory.
    -parsableExts="html, txt, md": Parsable file extensions separated by commas.
    -saveAs="index.html": Save compiled files as named.
    -startWith="index": Name without extension of the first file that will by parsed.
    -wordSeparator="-": Word separator used to replace spaces in URLs.
