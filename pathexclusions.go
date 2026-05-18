package cnty

import (
	S "strings"
)

var excludePrefixes = []string { ".", "_" }

var excludeContains = []string { ".." }

// We don't necessarily want to exclude JSON 
// files, but we have to for now because we
// are auto-generating them as outputs. 
var excludeSuffixes = []string {
    "~", ".env", ".sh", ".rc", ".bashrc", "gtk", "gtr",
    "_echo", "_tkns", "_tree", ".tmp.json" }

// excludeFilenamepath returns true (plus a reason) for a file base
// name that matches a blacklist for prefix or "midfix" or suffix.
//  - Excluded prefixes must follow a path separator; this rule should
//    allow "." and "./" 8only( to pass thru unexcluded / unmolested.
//  - Excluded suffixes apply to all names, but will not apply to 
//    a directory name that has a path separator appended. 
func excludeFilenamepath(s string) (bool, string) {
     var reason string  
     for _, pfx := range excludePrefixes {
     	 if S.HasPrefix(s, pfx) {
	    reason += "prefix<" + pfx + "> " 
	    }
	 }
     for _, sfx := range excludeSuffixes {
     	 if S.HasSuffix(s, sfx) {
	    reason += "suffix<" + sfx + "> " 
	    }
	 }
     for _, fpc := range excludeContains {
     	 if S.Contains(s, fpc) {
	    reason += "contains<" + fpc + "> " 
	    }
	 }
     for _, pfx := range excludePrefixes {
     	 if S.Contains(s, "/" + pfx) {
	    reason += "/+prefix<" + pfx + "> " 
	    }
	 }
     return (reason != ""), reason 
}

