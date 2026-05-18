package cnty

import (
	"fmt"
	"crypto/md5"
	FU "github.com/fbaube/fileutils"
	L "github.com/fbaube/mlog"
	"github.com/fbaube/m5db"
	SU "github.com/fbaube/stringutils"
	CA "github.com/fbaube/contentanalysis"
	CT "github.com/fbaube/ctoken"
)

// NewContentityFromString returns a Contentity Nord where the
// Nord (i.e. parent and ordered children) is empty. This func
// does not require file system access. If there is an error,
// and the returned value is non-nil, then the error is returned
// in the embedded Errer. However this func should only be called
// from a well-defined context, so error reporting can be curt. 
// .
func NewContentityFromString(inString, filext string) *Contentity {
	if len(inString) == 0 {
		return nil
	}
	if filext == "" || filext == "." {
		return nil
	}
	// =============================
	//  Note that the path is nil ! 
	//   Declare some useful vars.
	// =============================
	var pCty     *Contentity
	var pCR *m5db.ContentityRow
	var pFSO  *FU.FSObject
	var pCA   *CA.ContentAnalysis
	
	pCty = new(Contentity)
	// fmt.Printf("\t Nord seqID %d \n", p.SeqID())
	L.L.Okay(SU.Ybg("===> New Contentity: (no path)"))
	// ============================
	//  FSObject has its zero value. 
	// ============================
	// Code lifted from [fileutils/fsitem.go] LoadContents(),
	// which reads the file into the field [TypedRaw], takes
	// the hash, and quickly checks for XML and HTML5 declarations.
	// -----------------------------
	// Allocate this to prevent NPEs
	pCty.FSO.TypedRaw = new(CT.TypedRaw)
	if len(inString) > FU.MAX_FILE_SIZE {
		L.L.Warning("NewContentityFromString: " +
			"too large (%d)", len(inString))
		pCty.FSO.TypedRaw.Raw_type = SU.Raw_type_NIL
		pCty.SetError(fmt.Errorf("NewContentityFromString: " +
		       "content string too large: %d", len(inString)))
		return pCty 
	}
	if len(inString) < FU.MIN_FILE_SIZE { // Suspiciously tiny ?
		L.L.Warning("NewContentityFromString: tiny (%d)", len(inString))
		pCty.FSO.TypedRaw.Raw_type = SU.Raw_type_NIL
		pCty.SetError(fmt.Errorf("NewContentityFromString: " +
                       "content string too small: %d", len(inString)))
		return pCty 
	}
	pCty.FSO.Raw = CT.Raw(inString)
	// Take the hash and set the field.
	// pCty.Hash = *new([16]byte)
        pCty.FSO.Hash = md5.Sum([]byte(inString))

	// =================================
	//  Now we can load in the contents
	//  string that was passed in.
	// =================================
	var e error
	// =============================
	//  "Promote" to a ContentAnalysis
	// =============================
	// NewContentAnalysis return (nil,nil) for DIRLIKE 
	pCA, e = CA.NewContentAnalysis(pFSO)
	if e != nil {  
	   L.L.Error("NewContentityFromString: %s", e)
	   println("LINE 110")
	   pCty.SetError(fmt.Errorf("FSI.NewContentAnalysis.FromString: %w", e))
	   return pCty 
	}
	if pCA == nil { panic("WTF") }
	// =================================
	//  "Promote" to a ContentityRecord
	// =================================
	pCR, e = m5db.NewContentityRow(pFSO)
	if e != nil {
		L.L.Error("NewContentityFromString: %s", e)
	   	println("LINE 92")
		panic("newcntyrow.fromstr:newcntyrow")
		return pCty // ,errors.New("newcntyrow.fromstr:newcntyrow")
	}
	if pCR.RawType() == "" { // or SU.MU_type_UNK {
		panic("UNK MarkupType in NewContentity")
	}
	// NOW if we want to exit, we can
	// do the necessary assignments
	pCty.ContentityRow = *pCR
	if pFSO.IsDirlike() {
		L.L.Info(SU.Ybg(" Directory " + SU.Tildotted(pFSO.FPs.AbsFP)))
		pCty.ContentityRow.FSO = *pFSO
		return pCty 
	}
	L.L.Info(SU.Gbg(" " + pFSO.String() + " "))

	// ==================================
	//  Now fill in the ContentityRecord
	// ==================================
	pCty.GLinks = *new(GLinks)
	// println("D=> NewContentity:", p.String()) // p.MType, p.AbsFP())
	// fmt.Printf("D=> NewContentity: %s / %s \n", p.MType, p.AbsFP())
	return pCty 
}
