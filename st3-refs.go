package cnty

import (
	L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
)

// st3_Refs gathers all the various types of internal 
// and external references. Summary:
//   - 3a DoBlockList() // make list of blocks
//   - 3b DoGLinks()    // gather links
//   - 3c DoTableOfContents() // make ToC
//
// Some elms are not clearly either BLK or INL, so we
// do (or should) permit an indeterminate third state.
// .
func (p *Contentity) st3_Refs() *Contentity {
	if p.HasError() {
		return p
	}
	p.Lstage = "33"
	p.L(LDebug, "=== 33:Refs ===")
	ret := p.DoBlockList().DoGLinks().DoTableOfContents()
	if !p.HasError() { p.L(LOkay, "=== 33:Refs: Success ===") }
	return ret 
}

// DoBlockList makes a list of all the nodes that are
// blocks, so that they cn be traversed for rendering,
// and targeted for references.
// .
func (p *Contentity) DoBlockList() *Contentity {
	if p.HasError() {
		return p
	}
	switch p.RawType() {
	case SU.Raw_type_XML:
		L.L.Warning("TODO> 3a. DoBlockList XML")
	case SU.Raw_type_HTML:
		L.L.Warning("TODO> 3a. DoBlockList HTML")
	case SU.Raw_type_MKDN:
		L.L.Warning("TODO> 3a. DoBlockList MKDN")
	}
	return p
}

// DoGLinks gathers links.
// .
func (p *Contentity) DoGLinks() *Contentity {
	if p.HasError() {
		return p
	}
	switch p.RawType() {
	case SU.Raw_type_XML:
		L.L.Info("Calling GatherXmlGLinks...")
		p.GatherXmlGLinks()
		L.L.Info("Called! GatherXmlGLinks")
	case SU.Raw_type_HTML:
		L.L.Warning("TODO> 3b. DoGLinks HTML")
	case SU.Raw_type_MKDN:
		L.L.Warning("TODO> 3b. DoGLinks MKDN")
	}
	return p
}

// DoTableOfContents makes a ToC.
// .
func (p *Contentity) DoTableOfContents() *Contentity {
	if p.HasError() {
		return p
	}
	switch p.RawType() {
	case SU.Raw_type_XML:
		L.L.Warning("TODO> 3c. DoTableOfContents XML")
	case SU.Raw_type_HTML:
		L.L.Warning("TODO> 3c. DoTableOfContents HTML")
	case SU.Raw_type_MKDN:
		L.L.Warning("TODO> 3c. DoTableOfContents MKDN")
	}
	return p
}
