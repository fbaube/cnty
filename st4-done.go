package cnty

import (
	// L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
)

// st4_Done does final cleanup and beautification.
// .
func (p *Contentity) st4_Done() *Contentity {
	if p.HasError() {
		return p
	}
	// p.L(LProgress, "Done")
	p.L(LDebug, "=== 44:Done ===")
	switch p.RawType() {
	case SU.Raw_type_XML:
		// p.L(LWarning("TODO> 4. Done XML")
	case SU.Raw_type_HTML:
		// p.L(LWarning("TODO> 4. Done HTML")
	case SU.Raw_type_MKDN:
		// p.L(LWarning("TODO> 4. Done MKDN")
	}
	if !p.HasError() { p.L(LOkay, "=== 44:Done: Success ===") }
	return p // ret 
}
