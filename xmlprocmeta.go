package cnty

import (
	CT "github.com/fbaube/ctoken"
	"github.com/fbaube/gtree"
	// SU "github.com/fbaube/stringutils"
)

// RefineDirectives scans to patch Directives with correct keyword.
func (p *Contentity) RefineDirectives() error {
	// pX := p.TheXml()

	var pTag *gtree.GTag
	for _, pTag = range p.GTags {
		if pTag.TDType != CT.TD_type_DRCTV {
			continue
		}
		// Here are the directives we are dealing with.
		// Promote the Directive subtype keyword at "keynoun" to "Type", and
		// then promote the first word of "keyargs" (the Name) to "keynoun"".
		// DOCTYPE  Name ExtID [SysID]
		// ENTITY   Name EntDef   // General entity
		// ENTITY % Name EntDef   // Parameter entity (DTD only)
		// ELEMENT  Name contentspec
		// ATTLIST  Name AttDef's
		// NOTATION Name ExtID
		// ilog.Printf("Dir.PRE |%s|%s|", RT.string1, RT.string2)

		pTag.TDType = CT.TDType(pTag.Text)
		if pTag.TDType == CT.TD_type_DRCTV {
			panic("YIKES, leftover Dir Tagtype")
		}
		panic("cnty/xmlprocmeta.go :L125")
		// !! pTag.TagOrPrcsrDrctv, pTag.Datastring = SU.SplitOffFirstWord(pTag.Datastring)

		println("    --> RefineDirectives:", pTag.TDType, ",", pTag.Text) // TagOrPrcsrDrctv)

		if pTag.TDType == CT.TD_type_Entity { // !! && pTag.TagOrPrcsrDrctv == "%" {
			pTag.EntityIsParameter = true
			panic("cnty/xmlprocmeta.go:L40")
			// !! pTag.TagOrPrcsrDrctv, pTag.Datastring = SU.SplitOffFirstWord(pTag.Datastring)
			// pXI.GotDtdDecls = true
		}
		// fmt.Printf("Normalized Directive: %+v \n", RT)
		// At this point we are OK:
		// Type:DOCTYPE  keynoun:HTML keyargs:  PUBLIC "-//W3C//DTD HTML 4.0 ..."
		// Type:ENTITY   keynoun:foo  keyargs: "FOO"  entityIsParameter:false
		// Type:ENTITY   keynoun:bar  keyargs: "BAR"  entityIsParameter:true
		// Type:ELEMENT  Name contentspec
		// Type:ATTLIST  Name AttDef's
	}
	return nil
}
