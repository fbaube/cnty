package cnty

import (
	"io/fs"
)

// handleWalkerErrorArgument handles calls to the 
// [fs.WalkDirFunc] where the passed-in error is non-nil:
//  - The err argument reports an error related to path,
//    such that WalkDir will not walk into that directory.
//    The function can decide how to handle that error; 
//    returning the error causes WalkDir to stop walking
//    the entire tree.
//  - WalkDir calls the function with a non-nil err argument
//    in two cases:
//  - First, if the initial Stat on the root directory fails,
//    WalkDir calls the function with path set to root, d set
//    to nil, and err set to the error from fs.Stat.
//  - Second, if a directory's ReadDir method (see
//    https://pkg.go.dev/io/fs#ReadDirFile) fails, WalkDir
//    calls the function with path set to the directory's 
//    path, d set to an DirEntry describing the directory, 
//    and err set to the error from ReadDir. In this second 
//    case, the function is called twice with the path of 
//    the directory: the first call is before the directory 
//    read is attempted and has err set to nil, giving the 
//    function a chance to return SkipDir or SkipAll and 
//    avoid the ReadDir entirely. The second call is after 
//    a failed ReadDir and reports the error from ReadDir.
// . 
func (p *ContentityFS) handleWalkerErrorArgument(inPath string, inDE *fs.DirEntry, inErr error) error {

	// First call to walker func ? 
	if // p.mustInitRoot() {
	      len(p.asSlice) == 0 || len(p.asMapOfAbsFP) == 0 {
	   	// If we're not init'd yet, then the initial Stat
		// on the root directory failed, so WalkDir calls 
		// the function with path set to root, d set to 
		// nil, and err set to the error from fs.Stat .
		if inDE != nil {
		   return &fs.PathError { Path:inPath, Err:inErr,
		   	Op:"ctyfswalker bad root state" } 
		}
		return &fs.PathError { Path:inPath, Err:inErr,
		       	Op:"ctyfswalker root stat failed" }
	} else {
		// Else it's a dir, and ReadDir on it failed. So,
		// this func has been called with path set to the 
		// directory's path, d set to a DirEntry describing
		// the directory, and err set to ReadDir's error.
		return &fs.PathError { Path:inPath,
			Op:"cntyfswalker.readdir", Err:inErr }
	}	
}
