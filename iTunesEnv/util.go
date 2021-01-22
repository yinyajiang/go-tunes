package itunes

import (
	tools "github.com/yinyajiang/go-ytools/utils"
)

//Elem ...
type Elem struct {
	Recursion bool     `json:"recursion"`
	Srcs      []string `json:"srcs"`
	DstDir    string   `json:"dstDir"`
	Exts      []string `json:"exts"`
	Exclues   []string `json:"exclues"`
}

//CopyElem ...
func copyElem(elems []Elem) {
	if len(elems) == 0 {
		return
	}
	for _, elem := range elems {
		for _, src := range elem.Srcs {
			var files []string
			if !tools.IsDir(src) && tools.IsInFilter(src, elem.Exts) {
				files = append(files, src)
			} else {
				if elem.Recursion {
					files = tools.FilterDeepFile(src, elem.Exts)
				} else {
					files = tools.FilterFile(src, elem.Exts)
				}
			}
			for _, file := range files {
				if !tools.IsInFilter(file, elem.Exclues) {
					tools.CopyFile(file, tools.ThePath(elem.DstDir, tools.PathName(file)))
				}
			}
		}
	}
}
