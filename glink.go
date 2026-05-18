package cnty

import (
	FU "github.com/fbaube/fileutils"
	"github.com/fbaube/gtree"
)

type Flags int

const (
	IsRef      Flags = 1 << iota // 1 << 0 i.e. 0000 0001
	IsExtl                       // 1 << 1 i.e. 0000 0010
	IsURI                        // 1 << 2 i.e. 0000 0100
	IsKey                        // 1 << 3 i.e. 0000 1000
	IsResolved                   // 1 << 4 i.e. 0001 0000
)

func (b Flags) Set(flag Flags) Flags   { return b | flag }
func (b Flags) Reset(flag Flags) Flags { return b & ^flag }
func (b Flags) IsSet(flag Flags) bool  { return b&flag != 0 }

// func Toggle(b, flag Flags) Flags { return b ^ flag }

func (f Flags) String() string {
	var deff = "Def"
	var intl = "Intl"
	var ltyp = "N/A"
	var resd = "unReslvd"
	if f.IsSet(IsRef) {
		deff = "Ref"
	}
	if f.IsSet(IsExtl) {
		intl = "Extl"
	}
	if f.IsSet(IsURI) {
		ltyp = "Uri"
	}
	if f.IsSet(IsKey) {
		ltyp = "Key"
	}
	if f.IsSet(IsResolved) {
		resd = "Resolved"
	}
	return deff + "," + intl + "," + ltyp + "," + resd
}

// GLinks is used for (1) intra-file ref resolution,
// (2) inter-file ptr resolution, (3) ToC generation.
type GLinks struct {
	// OwnerP points back to the owning struct, so that
	// `GLink`s can be processed easily as simple data structures.
	OwnerP interface{}
	// KeyRefncs are outgoing key-based links/references
	KeyRefncs []*GLink // (Extl|Intl)KeyReferences
	// KeyRefnts are unique key-based definitions that are possible
	// referents (resolution targets) of same or other files' [KeyRefncs]
	KeyRefnts []*GLink // (Extl|Intl)KeyDefs
	// UriRefncs are outgoing URI-based links/references
	UriRefncs []*GLink // (Extl|Intl)UriReferences
	// UriRefnts are unique URI-based definitions that are possible
	// referents(resolution targets) of same or other files' [UriRefncs]
	UriRefnts []*GLink // (Extl|Intl)UriDefs
}

// GLink summarizes a link (or key) (or reference) found in markup content.
// It is either URI-based (`href conref id`) or key-based (`key keyref`). It
// applies to all LwDITA formats, but not all fields apply to all LwDITA formats.
type GLink struct {
	// IsRefnc - else is Refnt (Referents are much more numerous)
	IsRefnc bool
	// IsExtl - else is Intl (which are more numerous)
	IsExtl bool
	// AddressMode is "http", "key", "idref", "uri"
	AddressMode string
	// Att is the XML attribute - id, idref, href, xref, keyref, etc.
	Att string
	// Tag is the tag that has this link-related attribute of interest
	Tag string
	// Link_raw as redd in during parsing
	Link_raw string
	// RelFP can be a URI or the resolution of a keyref.
	// "" if target is in same file; NOTE This is relative to the
	// sourcing file, NOT to the current working directory during parsing!
	RelFP string
	// AbsFP can be a URI or the resolution of a keyref.
	// "" if target is in same file
	AbsFP FU.AbsFilePath
	// TopicID iff present (but isn't it mandatory ?)
	TopicID string
	// FragID is peeled off from Raw
	FragID string
	// Resolved is used to narrow in on difficult cases
	Resolved bool
	// LinkedFrom is the GTag where the GLink is defined
	LinkedFrom *gtree.GTag
	// Original can be nil: it is the tag where the GLink is resolved to,
	// i.e. the REFERENT, and is quite possibly in another file, which we
	// hope we also have available in memory.
	Original *gtree.GTag
}
