package cnty

import (
	"fmt"
	"io"
	"github.com/fbaube/gtoken"
	"github.com/fbaube/gparse"
	"github.com/fbaube/gtree"
	"github.com/fbaube/m5db"
	N "github.com/fbaube/nork"
	SU "github.com/fbaube/stringutils"
)

type ContentityStage func(*Contentity) *Contentity

// For the record, ignore the API of
// https://godoc.org/golang.org/x/net/html#Node

// Contentity embeds [Nork] - and embeds [ContentityRow] embeds
// ([FSObject] embeds [Errer] - but Contentity does NOT embed [FSONork].
// .
type Contentity struct { // has Raw

     	// Nork provides hierarchical structure, only.
	// It is embedded as an instance, not as a pointer.
	N.Nork
	// ====================================
        //  Substructures for: Adjacency Lists
        // ====================================
        // Every implementation of adjacency lists points
        // (1) to parent, (2) to peers in own list, and
        // (3) to first & last kids in linked list of all its kids.
        // ------------------------------------------------
        //  Substructure for Adjacency List #1 based on Go
        //  ptrs-to-structs (not indices). Note that these
        //  are provided IN PARALLEL with the ptr-to-struct
        //  slices in the embedded Nork. This gives us
        //  redundancy for error checking, and maybe also
        //  lets us run performance comparisons.
        // ------------------------------------------------
	// THESE MUST DIE, REPLACED BY NORK
	// ================================
//      parent             N.Norker // level up   // dupe of Nork.prnt
        prevPeer, nextPeer N.Norker // level same // like Nork.kids []*Nork
        firstKid, lastKid  N.Norker // level down // like Nork.kids []*Nork
	// =========================================
        //  Substructures for: Data persisted to DB
        // =========================================
	// ContentityRow includes all fields that get persisted 
	// to the DB. It contains the field Raw (deeply embedded),
	// and embeds [FSObject] embeds [Errer]. 
	m5db.ContentityRow
	
	// LogInfo is (the index of the Contentity in 
	// the larger slice) + (the processing stage ID)
	LogInfo
	// logIdx int
	// logStg string
	
	// ParserResults is parseutils.ParserResults_ffs
	// (ffs = file format -specific = "html" or "mkdn" but not
	// "xml" cos Go's XML parser does not produce a tree structure) 
	ParserResults interface{}

	GTokens      []*gtoken.GToken
	GTags        []*gtree.GTag
	*gtree.GTree // maybe not need GRootTag or RootOfASTptr
	GTknsWriter, GTreeWriter,
	GEchoWriter io.Writer
	GLinks

	// GEnts is "ENTITY"" directives (both with "%" and without).
	GEnts map[string]*gparse.GEnt
	// DElms is "ELEMENT" directives.
//	DElms map[string]*gtree.GTag

	TagTally StringTally
	AttTally StringTally

	// FU.OutputFiles // This was useful at one point

	// ======================
	//  STUFF FOR FUTURE USE
	// ======================
/*
        // ------------------------------------------
        //  Substructure for Adjacency List #2 based
        //  on discrete indices into "arena" slice
        // ------------------------------------------
        iParent              int // level up
        iPrevPeer, iNextPeer int // level same
        iFirstKid, iLastKid  int // level down
*/
/*
        // ------------------------------------------
        //  Substructure for Adjacency List of KIDS
        //  where they are stored on a single string
        //  field that holds all applicable indices.
        // ------------------------------------------
        // kidIdxs when empty is "," (or ""), else
        // e.g. ",1,4,56,". The kidIdxs should be in
        // the same order as the Kid nodes themselves.
        // Comma-bracketing simplifies search (",%d,").
        // ------------------------------------------
        kidIdxs string
*/
}

// LogInfo exists mainly to provide a grep'able string:
// for example "(01:4a)", where 01 is the index of the
// Contentity and 4a is the processing stage. This is
// obv a candidate for replacement by stdlib's slog.
//
// The [io.Writer] field W exists outside of the
// [github.com/fbaube/mlog] logging subsystem 
// and should only be used if `mlog` is not.
// .
type LogInfo struct {
	Lindex int
	Lstage string
	W io.Writer
	}

func (p *LogInfo) String() string {
     return fmt.Sprintf("(%02d:%s)", p.Lindex, p.Lstage) 
     }

/*
func (p *Contentity) IsDir() bool {
	return p.FSObject.IsDir()
}

func (p *Contentity) IsDirlike() bool {
	return p.FSObject.IsDirlike()
}

type norderCreationState struct {
	// nexSeqID int // reset to 0 when doing another tree ?
	rootPath string
}
*/

// String is developer output. Hafta dump:
// FU.InputFile, FU.OutputFiles, GTree,
// GRefs, *XmlFileMeta, *XmlItems, *DitaInfo
func (p Contentity) String() string {
	var sGTree string
	if p.GTree != nil {
		sGTree = p.GTree.String()
	}
	// s := fmt.Sprintf("[len:%d]", p.Size())
	s := fmt.Sprintf("||%s||GTree|%s||OutbKeyLinks|%+v|KeyLinkTgts|%+v|OutbUriLinks|%+v|UriLinkTgts|%+v||",
		SU.Tildotted(p.FSO.FPs.AbsFP), sGTree, p.KeyRefncs, 
		p.KeyRefnts, p.UriRefncs, p.UriRefnts)
	/* code to use ?
			if p.XmlFileMeta != nil {
				s += fmt.Sprintf("XmlFileMeta|%s||", p.XmlFileMeta.String())
			}
		* /
		if p.IDinfo != nil {
			s += fmt.Sprintf("xf.IDinfo|%s||", p.IDinfo.String())
		}
	if p.GEnts != nil {
		// FIXME s += fmt.Sprintf("GEnts|%s||", p.GEnts.String())
		 * 	}
		 * 	if p.DElms != nil {
		// FIXME s += fmt.Sprintf("DElms|%s||", p.DElms.String())
	}
	== */
	// if p.DitaInfo != nil {
	s += fmt.Sprintf("DitaInfo|Flav:%s|Cntp:%s|", p.DitaFlavor, p.DitaContype)
	// }
	return s
}
