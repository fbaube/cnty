package cnty

import (
	"io/fs"
	"os"
	"fmt"
	"errors"
	FP "path/filepath"
	S "strings"
	// "time"

	FU "github.com/fbaube/fileutils"
	SU "github.com/fbaube/stringutils"
	CT "github.com/fbaube/ctoken"
	L "github.com/fbaube/mlog"
)

// NewContentityFS is the entrypoint for processing an input
// directory tree of files and it proceeds as follows:
//  - initialize
//  - create an [os.Root] (see https://pkg.go.dev/os@go1.26.3#Root )
//  - create a [nork.NorkFactory] 
//  - walk the RootFS, creating Contentities
//    and appending them to a slice
//  - process the list to identify and make
//    parent/child links
//
// ContentityFS embeds [FSObject] embeds [Errer]. 
//
// About [os.Root] and RootFS:
//  - Root may be used to only access files within a single directory tree.
//  - Methods on Root can only access files and directories beneath a root
//    directory. If any component of a file name passed to a method of Root
//    references a location outside the root, the method returns an error.
//    File names may reference the directory itself (.).
//  - Methods on Root will follow symbolic links, but symbolic links may
//    not reference a location outside the root. Symbolic links must not
//    be absolute.
//  - Methods on Root do not prohibit traversal of filesystem boundaries,
//    Linux bind mounts, /proc special files, or access to Unix device files.
//  - Methods on Root are OK for use from multiple goroutines simultaneously.
//  - On most platforms, creating a Root opens a file descriptor or handle
//    referencing the directory. If the directory is moved, methods on Root
//    reference the original directory in its new location.
//
// The path argument should probably be an absolute FP, 
// because a relative FP might cause problems. Note that
// this is the opposite of the advice for lower-level items.
//
// Note that when we used [os.DirFS], it appeared 
// to make no difference whether path 
//  - is relative or absolute
//  - ends with a trailing slash or not
//  - is a directory or a symlink to a directory
// It probably needs serious testing, especially now with RootFS.
//
// The only error returns for this func are:
//  - a bad path, rejected by FU func [NewFilepaths] 
//  - the path is not a directory (altho it can be
//    a symlnk to a directory ?)
//  - TBD: WHat happens of [os.Root] barfs on something ? 
//
// TODO Add two flags ?
// Maybe it needs two boolean arguments:
//  - One to say whether to be strict about security 
//    (using [os.Root] (YES!!) and Valid/Local, and
//  - One to say whether to follow symlinks.
// 
// These two flags might have some interesting interactions.
// OTOH since this func can (and does?) use [os.Root], it can 
// easily (and should probably) also default to higher security 
// using funcs [io/fs.ValidPath] and [path/filepath.IsLocal].
//
// Accumulated NewContentity errors are counted
// in the field CotentityFS.nErrors 
// .
func NewContentityFS(aPath string, okayFilexts []string) (*ContentityFS, error){
     	var pCntyFS = new(ContentityFS)
	var e error 
	var pRootCnty    *Contentity 
	var pCntyFactory *ContentityFactory 
	var pPE *os.PathError
	pPE = new(os.PathError{Path:aPath})

	// -----------------------
	//  Prepare filepaths and
	//   check for directory 
	// -----------------------
	L.L.Info("Making NewContentityFS: " + aPath)
	pRootCnty = NewContentity(aPath)
	pCntyFS.rootCnty = pRootCnty
	if pRootCnty.HasError() {
	     	L.L.Error("NewCntyFS: bad path: %s", aPath)
		pPE.Op = "newcntyfs:newfso"
		pPE.Err = fmt.Errorf("bad root path: %w",
			  pRootCnty.GetError())
		return pCntyFS, pPE 
	}
	pathToUse := pRootCnty.FSO.FPs.RelFP // FU.EnsureTrailingPathSep
	pPE.Path = pathToUse
	// [os.DirFS] does not check or report problems with
	// the path argument, while OTOH [os.Root] is groovy.
	if pRootCnty.FSO.FPs.DoesNotExist || !pRootCnty.FSO.FPs.IsDir {
		L.L.Error("NewCntyFS: Not a directory: %s", aPath)
		pPE.Op = "newcntyfs.root"
		pPE.Err = errors.New("not a valid directory")
		return pCntyFS, pPE
	}
	// -----------------
	// Create the RootFS
	// -----------------
	// 2025.01 Change from RelFP to AbsFP (Altho 
	// there was probly a good reason to use RelFP)
	L.L.Info("Path for new os.RootFS: " + aPath)
	var osRoot *os.Root 
	osRoot, e = os.OpenRoot(aPath)
	pCntyFS.FS = osRoot.FS()
	// ----------------------
	//  Create root Nork and
	//   FSOTreeNorkFactory
	// ----------------------
//	pCntyFactory, pRootCnty, e = N.NewFSOTreeNorkFactory(pathToUse)
	pCntyFactory = new(ContentityFactory)
	pCntyFactory.rootPath = pathToUse
/*
	if e != nil {
	   L.L.Error("contentityfs_new: new factory failed: %w", e)
	   return nil, e
	   }
*/
	// Initialize slice & map. Their length 
	// 0 will be detected by func [mustInit]
	pCntyFS.asSlice = make([]*Contentity, 0)
	pCntyFS.asSlice = append(pCntyFS.asSlice, pRootCnty)
	pCntyFS.asMapOfAbsFP = make(map[string]*Contentity)
	pCntyFS.asMapOfAbsFP[pCntyFS.rootCnty.FSO.FPs.AbsFP] = pRootCnty

	var initHasBeenDone bool 

	// ==================
	//    FIRST PASS
	//  Load slice & map
	// ==================
	// NOTE that rel.path "." seems to be necessary 
	// here or else really weird errors occur.
// ======================================================================
//	e = fs.WalkDir(pCntyFS.FS, ".", wfnBuildContentityTree)
//
// wfnBuildContentityTree is a [fs.WalkDirFunc]
// 
// Altho it returns [*fs.PathError], this has to be declared as error
// because of the problems of an interface that is both nil and not nil.
//
// The basic procedure of the func is is:
//  - handle an error passed in 
//  - check validity of path argument (and reject if it is a file)
//  - handle first time thru (i.e. root node) 
//  - filter out unwanted values (and if unwanted dir, return [fs.SkipDir])
//  - add to slice - and also map - whether dir or file
//  - use materialised paths in slice to form links to build a tree 
//
// Note that [os.Root] should handle symlinks securely, 
//
// This func filters out several file types:
//  - hidden (esp'ly .git directory)
//  - leading underbars ("_")
//  - emacs backup ("myfile~")
//  - this app's info  files: "*gtk,*gtr"
//  - this app's debug files: "*_(echo,tkns,tree)"
//  - filenames without a dot (indicating no file extension)
//  - NOTE that zero-length files (no content to analyse)
//    should NOT be filtered out 
//
// Note that as path separator, "/" is usually assumed, not [os.PathSep]. 
// .
// wfnBuildContentityTree(inPath string, inDE fs.DirEntry, inErr error) error {
// ======================================================================
e = fs.WalkDir(pCntyFS.FS, ".",
func(inPath string, inDE fs.DirEntry, inErr error) error { // fs.WalkDirFunc
	// --------------------------
	//  Were we passed an error?
	// --------------------------
	if inErr != nil {
	   	 return pCntyFS.handleWalkerErrorArgument(
			inPath, &inDE, inErr)
	}
	// --------------------
	//  Set some variables 
	// --------------------
//	var isFirst = // pCntyFS.mustInitRoot() // first call ?
//	    len(pCntyFS.asSlice) == 0 && len(pCntyFS.asMapOfAbsFP) == 0
	var inName  = inDE.Name()
	var inDEisDir = inDE.IsDir()
	// If it's a directory, make sure it has a trailing slash.
	if inDEisDir {
	   inPath = FU.EnsureTrailingPathSep(inPath)
	   inName = FU.EnsureTrailingPathSep(inName)
	}
	// func [filepath.Abs] fails here cos it needs more than just 
	// a Base file name cos it does only lexical processing. 
	// absfp,_ = FP.Abs(path)

	// var p *Contentity
	var e error 
	// ==================
	//  HANDLE ROOT NODE 
	// (without filtering)
	// ==================
//	if isFirst {
	if !initHasBeenDone {
		L.L.Info("pCntyFSWalker: inPath: " + inPath)
	   	if !inDEisDir { return &fs.PathError { Path:inPath,
		   	Op:"cntyfswalker.root", Err:errors.New("not a dir") } }
		L.L.Debug("cntyfswalker.root: path: %s / %s", inName, inPath)
		L.L.Debug("cntyfswalker.root: dirEntry: %+v", inDE)
	   	e = pCntyFS.doInitRoot()
		initHasBeenDone = true
		if e == nil { return nil }
		return &fs.PathError { Err:e, Path:inPath,
		       Op:"newrootcnty.doinitroot" }
	}
	// ---------------------------
	//  Filter out unwanted stuff 
	// ---------------------------
	bad, rsn := excludeFilenamepath(inPath)
	if bad {
		L.L.Debug("Rejecting (%s): %s", inPath, rsn)
		if inDEisDir { return fs.SkipDir } 
		return nil
	}
	// -----------------------------------------------
	//  Now at this point, even if it's a directory,
	//  it's OK ! So let's go ahead and form the path
	//  of the file-or-dir and make the Contentity
	// -----------------------------------------------
	absPathToUse := FU.EnsureTrailingPathSep(
		        FP.Join(pCntyFS.RootAbsPath(), inPath))
	var pCty *Contentity
	pCty = NewContentity(absPathToUse)
	if pCty.HasError() { 
		L.L.Warning("Rejecting (new Contentity(%s) failed): %T %+v",
			absPathToUse, e, e)
		pCntyFS.nErrors++
		return nil
	}
	// ================================
	//  Now we do things based on just 
	//  what exactly the input DirEntry
	//  (inDE) is - file, dir, wotevs.
	// ================================
	// This is where bugs appeared when it's a directory,
	// cos calling code was assuming a valid Contentity.
	// TODO: Not sure what happens with symlinks
	if inDEisDir {
	   	if pCty.FSO.TypedRaw == nil {
		   pCty.FSO.TypedRaw = new(CT.TypedRaw)
		   } 
	        pCty.FSO.TypedRaw.Raw_type = SU.Raw_type_DIRLIKE
		pCntyFS.nDirs++ // just a simple counter
		// println("================ DIR ========")
		// These next two stmts should barf, cos
		// they should not be allocated for a dir !
		// p.MimeType = "dir"
		// p.MType = SU.MU_type_DIRLIKE
		L.L.Okay("Item (DIR/) OK; no MmeType or MType") 
	} else if inDE.Type() == 0 { // regular file
		pCntyFS.nFiles++ // just a simple counter 
		L.L.Okay("Item (FILE) OK: MType<%s> MarkupType<%s>",
			pCty.MType, pCty.RawType())
	} else if (inDE.Type() & fs.ModeSymlink) != 0 { // Symlink
	       	if pCty.FSO.TypedRaw == nil {
		   pCty.FSO.TypedRaw = new(CT.TypedRaw)
		}
	        pCty.FSO.TypedRaw.Raw_type = SU.Raw_type_DIRLIKE // OOPS
		pCntyFS.nMiscs++ // just a simple counter 
		L.L.Okay("Item (SYML) OK: what to do ?!")
	} else { // Some weirdness in the Mode bits 
	        pCty.FSO.TypedRaw.Raw_type = SU.Raw_type_DIRLIKE
             // pCntyFS.nMiscs++ // just a simple counter
		pCntyFS.nErrors++
                L.L.Error("Item (WTF) BAD: what to do ?!")
	}
	// -------------------------
	//   Also add it to the
	//  arena-slice and the map
	// -------------------------
	pCntyFS.asSlice = append(pCntyFS.asSlice, pCty)
	pCntyFS.asMapOfAbsFP[absPathToUse] = pCty
	// L.L.Info("ADDED TO MAP L241: " + pathToUse)
	return nil
})
// ======================================================================
	if e != nil {
		// L.L.Panic("NewpCntyFS.WalkDir: " + e.Error())
		return nil, &fs.PathError { Op:"newcntyfs.walkdir",
		       Err:e, Path:aPath } 
	}
	L.L.Okay("NewpCntyFS: walked OK %d nords from path %s",
		 len(pCntyFS.asSlice), pathToUse)

	// Debuggery 
	for ii, cc := range pCntyFS.asSlice {
	    if cc == nil {
	       L.L.Error ("OOPS, pCntyFS.asSlice[%02d] is NIL", ii)
	       continue
	    }
	    /* if cc.FSObject == nil || cc.FSObject.FileMeta == nil {
	       L.L.Error("WTF, man!")
	       continue } */
	    if cc.FSO.IsDirlike() {
	        L.L.Debug("[%02d] isDIRLIKE: AbsFP: %s",
			ii, cc.FSO.FPs.AbsFP)
	    } else {
	     // fmt.Fprintf(os.Stderr, "ii is <%d> cc is <%p> \n", ii, cc)
	        fmt.Fprintf(os.Stderr, "path: %s\n", cc.CreatPath())
	     // fmt.Fprintf(os.Stderr, "cc.RawType is <%p> \n", cc.RawType())
	//	time.Sleep(200 * time.Millisecond)
	//	L.L.Debug("[%02d] MarkupType: %s", ii, cc.RawType())
	    }
	}

	// ==============================================
	//      SECOND PASS
	//  Range over the slice, using the materialised
	//  paths in asMapToAbsFS to identify parent/kid 
	//  Nord relationships and link together
	// ==============================================
	// TODO This needs to be in some generalized 
	// form, such as TreeFromMaterializedPaths
	// =========================================
	var i int
	var pC *Contentity
	for i, pC = range pCntyFS.asSlice {
		if i == 0 { // skip over root 
			continue
		}
		// ---------------------------
		//  Shortcut if child of root
		// ---------------------------
		if !S.Contains(pC.FSO.FPs.RelFP, FU.PathSep) {
			pCntyFS.rootCnty.AddKid(&pC.Nork)
			continue
		}
		// --------------------------
		//   Get dir portion of path
		// --------------------------
		itsDir := FP.Dir(pC.FSO.FPs.RelFP)
		itsDir = FU.EnsureTrailingPathSep(itsDir)
		// println(n.Path, "|cnex2|", itsDir)
		// L.L.Warning("itsDir: " + itsDir)
		// L.L.Warning("theMap: %+v", pCntyFS.asMap)
		var pPar *Contentity
		var ok bool
		// PROBLEMS HERE ?
		// The parent directory should be in the map.
		// If it's not, then possibly we have messed
		// up with trailing separators. 
		if pPar, ok = pCntyFS.asMapOfAbsFP[itsDir]; !ok {
			L.L.Error("findParentInMap: failed for: " +
				itsDir + " of " + pC.FSO.FPs.AbsFP)
			println(fmt.Sprintf("%+v", pCntyFS.asMapOfAbsFP))
			panic(pC.FSO.FPs.AbsFP)
		}
		/*
		if itsDir != par.AbsFP() { // or, Rel? 
			panic(itsDir + " != " + par.AbsFP())
		}
		*/
		pPar.AddKid(&pC.Nork)
	}
	// TODO Look for entries that do not have a parent assigned !
	
	/* more debugging
	println("DUMP LIST")
	for _, n := range pFTFS.asSlice {
		println(n.LinePrefixString(), n.LineSummaryString())
	}
	println("DUMP MAP")
	for k, v := range pFTFS.asMap {
		fmt.Printf("%s\t:: %s %s \n", k, v.LinePrefixString(), v.LineSummaryString())
	}
	*/
	// println(SU.Gbg("=== TREE ==="))
	// pCntyFS.rootNord.PrintAll(os.Stdout)
	return pCntyFS, nil
}
