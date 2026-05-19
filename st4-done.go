package cnty

import (
	L "github.com/fbaube/mlog"
	SU "github.com/fbaube/stringutils"
)

// st4_Done does final cleanup and beautification.
// .
func (p *Contentity) st4_Done() *Contentity {
	if p.HasError() {
		return p
	}
	// p.L(LProgress, "Done")
	L.L.Debug("=== 44:Done ===")
	switch p.RawType() {
	case SU.Raw_type_XML:
		// L.L.Warning("TODO> 4. Done XML")
	case SU.Raw_type_HTML:
		// L.L.Warning("TODO> 4. Done HTML")
	case SU.Raw_type_MKDN:
		// L.L.Warning("TODO> 4. Done MKDN")
	}
	if !p.HasError() { p.L(LOkay, "=== 44:Done: Success ===") }
	return p // ret 
}
