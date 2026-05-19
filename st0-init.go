package cnty

import (
           "errors"
	SU "github.com/fbaube/stringutils"
	 L "github.com/fbaube/mlog"
)

// st0_Init does pre-processing prep and checks.
//
// Input: A `Contentity` that has had its contents analyzed.
//
// For `st0_Init()` to work,
//   - the `Contentity` *must* refer to a readable file, and
//   - the field `Contentity.MType` *must* already be set
//
// Output: A `Contentity` that is in-memory and analyzed
// (shallowly) as `XML` or `MKDN` (Markdown) or `HTML`.
//
// - SetTypeSpecific()
// - SanityCheck()
// .
func (p *Contentity) st0_Init() *Contentity {
	if p.HasError() {
		return p
	}
	if p.MType == "" {
		// panic("st0_Init: nil MType")
		L.L.Error("st0_Init: nil MType")
	}
	p.Lstage = "00"
	p.L(LDebug, "=== 00:Init ===")
	// panic("TEST PANIC in st0")
	return p.st0a_SanityCheck()
}

// st0a_SanityCheck checks that `mcfile.MarkupType()` &
// `mcfile.IsXML()` are OK and that `MCFile.MType[]` is set.
func (p *Contentity) st0a_SanityCheck() *Contentity {
	p.Lstage = "0a"
	var errStr string 
	if p.FSO.TypedRaw.Raw_type != SU.Raw_type_DIRLIKE && p.MType == "" {
		errStr = "MType is empty"
	}
	if len(p.RawType()) < 3 || len(p.RawType()) > 7 {
		panic("BAD RawType in st0a: " + string(p.RawType()))
	}
	switch p.RawType() {
	case SU.Raw_type_XML:
		if !p.IsXML() {
			errStr = "is XML but: !XML?!"
		}
	case SU.Raw_type_MKDN:
		if p.IsXML() {
			errStr = "is Mkdn but: XML?!"
		}
	case SU.Raw_type_HTML:
		if !p.IsXML() {
			errStr = "is HTML but: !XML?!"
		}
	case SU.Raw_type_BIN:
		if p.IsXML() {
			// panic("Init error: is BIN but: XML?!")
			errStr = "is BIN but: XML?!"
		}
	case SU.Raw_type_SQL, SU.Raw_type_DIRLIKE:
	     // No problem!
	default:
		errStr = "bad/missing markup type: " +
			string(p.RawType())
	}
	if errStr != "" {
	// panic("errStr: " + errStr)
	   p.FSO.SetError(errors.New(errStr))
	}
	return p
}
