package cnty

import (
	"errors"
	"os"
	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	N "github.com/fbaube/nork"
	"github.com/fbaube/m5db"
	SU "github.com/fbaube/stringutils"
	CA "github.com/fbaube/contentanalysis"
)

// Arg s is a filepath. The returned [Contentity] embeds [Nork]
// and [FSO] (embeds [Errer]) but does NOT embed [FSONork].
// If (Errer.HasError], it is a [*os.PathError].
func newContentity(s string) *Contentity {
        var pC = new(Contentity)
        if s == "" {
                pC.SetError(errors.New("newcnty: missing path"))
                return pC
        }
//	pC.FSONork = FU.NewFSOLoneNork(s)
        pC.Nork = *N.NewNork(s)
        pC.FSO = *FU.NewFSObject(s)
        var pPE = new(os.PathError{Path:s})
        if pC.Nork.HasError() {
           pPE.Op = "newcnty:newnork"
           pPE.Err = pC.Nork.GetError()
           pC.SetError(pPE)
        }
        if pC.FSO.HasError() {
           pPE.Op = "newcnty:newfsonork"
           pPE.Err = pC.FSO.GetError()
           pC.SetError(pPE)
        }
        return pC
 }

// NewContentity returns a Contentity -nork (i.e. a [Nork] node 
// with content and an embedded [FSObject] ) that can NOT be the
// root of a Contentity tree. For error, see embedded [Errer].
//
// FIXME: Should we have Tree and Lone versions and Factory ?
//
// It should accept either an absolute or a relative FP, altho
// relative is preferred for various reasons, mainly because of
// security-related preferences of the path and filepath stdlibs.
//
// TODO: Maybe it needs two boolean arguments:
//  - One to say whether to be strict about security 
//    (using [os.Root] and Valid/Local, and
//  - One to say whether to follow symlinks.
// 
// These two flags might have some interesting interactions.
// Since this func could (but does not) use [os.Root], these can be
// left as calling options, rather than implementing higher security	
// using funcs [io/fs.ValidPath] and [path/filepath.IsLocal].
// 
// We want everything to be in a nice tree of Norks, and it means that
// we have to create Contenties for directories too (where `Raw_type
// == SU.Raw_type_DIRLIKE`), so we have to handle that case too. 
// .
func NewContentity(aPath string) *Contentity {
	var e error
	// Use the unexported utility func defined above. 
	// FIXME We are not using a ContentityFactory yet,
	// so we do not try to get a Factory root path 
	// If we were passed an Abs.FP, it's okay.
	// It was something like 
/*	if FP.IsAbs(aPath) {
		pFSONork.FSO = FU.NewFSObject(aPath)
	 } else {
	// else create an Abs.FP 
		pFSO = FU.NewFSObject(FP.Join(CntyEng.rootPath, aPath))
	}
*/
	var pNewCnty = newContentity(aPath)
	// If error, quick return.
        if pNewCnty.FSO.HasError() {
                return pNewCnty
        }
	// Have a pre-filled error ready 
	pPE := new(os.PathError{Path:aPath})
	L.L.Debug("NewContentity.FSO: %s", pNewCnty.FSO.Infos())
	// =====================================
	//  pNewCnty.(FSO,Nork) are OK, so from
	//    here on we can return pNewCnty 
	// =====================================
/*	var pFSONork = N.NewFSOLoneNork(aPath)
	if  pFSONork.HasError() {
	     	pPE.Op = "newcnty:newrootfsotreenork"
		pPE.Err = pNewCnty.FSO.GetError()
		pNewCnty.SetError(pPE)
		return pNewCnty
	}
	pNewCnty.Nork = *N.NewNork(aPath)
*/
	L.L.Okay(SU.Ybg("===> New Contentity: %s"), SU.Tildotted(aPath))

	// ======================================
	//  If it's a directory (or similar,
	//  such as symlink) handle it here cos
	//  we don't need to do ContentAnalysis.
	//  FIXME If it's a symlnk, we should probly
	//  read the target and store it somewhere. 
	// ======================================
	if pNewCnty.FSO.FPs.IsDirlike {
	   	// This should fail only if the item does not exist.
		pNewCntyRow, e := m5db.NewContentityRow(&pNewCnty.FSO)
		if e != nil {
			L.L.Error("NewContentity(Dirlike)<%s>: %s", aPath, e)
			println("LINE 122")
			pPE.Op = "newcnty:newcntyrow:dirlike"
			pPE.Err = e 
			pNewCnty.FSO.SetError(pPE) 
			return pNewCnty
		}
		L.L.Info(SU.Ybg(" Dir " + SU.Tildotted(pNewCnty.FSO.FPs.AbsFP)))
                pNewCntyRow.FSO = pNewCnty.FSO
		pNewCnty.ContentityRow = *pNewCntyRow
		return pNewCnty 
        }
	// =============================================
	//  Here onward, we should be dealing only with
	//   regular files, which can contain content.
	//   Start by forcing a fetch of the contents.
	// =============================================
	if !pNewCnty.FSO.IsFile() {
	     panic("LINE 132 it's not a file")
	}
	_, e = pNewCnty.FSO.Contents()
	// L.L.Warning("LENGTH %d", len(pNewCnty.FSO.TypedRaw.Raw))
	if e != nil {
   	   println("LINE 137")
	   pPE.Op = "newcnty.contents"
	   pPE.Err = e
	   pNewCnty.FSO.SetError(pPE)
	   return pNewCnty
	}
	// ===================================
	//  Now it gets interesting - working 
	//  with fields persisted to the DB.
	//  Declare some useful vars.
	// ===================================
	var pNewCntyRow *m5db.ContentityRow
	var pNewCntyAnlys *CA.ContentAnalysis
	// =======================
	//  "Promote" FSObject to
	//   a ContentAnalysis
	// =======================
	// NewContentAnalysis return (nil,nil) for DIRLIKE 
	pNewCntyAnlys, e = CA.NewContentAnalysis(&pNewCnty.FSO)
	if e != nil { 
	   L.L.Error("NewContentity(PP=>PA)<%s>: %s", aPath, e)
	   println("LINE 158")
	   pPE.Op = "newcnty:newcntanls"
	   pPE.Err = e
	   pNewCnty.FSO.SetError(pPE)
	   return pNewCnty
	}
	if pNewCntyAnlys == nil { panic("WTF") }
	// ===========================
	// "Promote" ContentAnalysis
	//        to ContentityRecord
	// ===========================
	pNewCntyRow, e = m5db.NewContentityRow(&pNewCnty.FSO)
	if e != nil { // pNewCntyRow.HasError() {
		L.L.Error("NewContentity(PA=>CR)<%s>: %s", aPath, e)
	   	println("LINE 172")
		pPE.Op = "newcnty:newcntyrow"
		pPE.Err = e
		pNewCnty.FSO.SetError(pPE)
		return pNewCnty
	}
	if pNewCntyRow.RawType() == "" { // or SU.MU_type_UNK {
		panic("UNK MarkupType in NewContentity")
	}
	// NOW if we want to exit, we can
	// do the necessary assignments
	pNewCnty.ContentityRow = *pNewCntyRow
	if pNewCnty.FSO.IsDirlike() {
	   	// Whoops, not sposta be here
		panic("Late Dirlike")
		L.L.Info(SU.Ybg(" Directory " + SU.Tildotted(pNewCnty.FSO.FPs.AbsFP)))
		pNewCnty.ContentityRow.FSO = pNewCnty.FSO
		return pNewCnty 
	}
	L.L.Info(SU.Gbg(" " + pNewCnty.FSO.String() + " "))

	// ==================================
	//  Now fill in the ContentityRecord
	// ==================================
	pNewCnty.GLinks = *new(GLinks)
	// println("D=> NewContentity:", p.String()) // p.MType, p.AbsFP())
	// fmt.Printf("D=> NewContentity: %s / %s \n", p.MType, p.AbsFP())
	return pNewCnty 
}
