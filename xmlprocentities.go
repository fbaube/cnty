package cnty

import (
	"fmt"
	S "strings"

	CT "github.com/fbaube/ctoken"
	"github.com/fbaube/gparse"
	SU "github.com/fbaube/stringutils"
)

// RULES
// https://www.xml.com/axml/axml.html
//
// "&"  "general"  entity references may appear in many places,
//       altho not at the top level of a DTD (i.e. inside defs only).
// "%" "parameter" entity references may appear in DTDs ONLY.
// They exist in different namespaces: "A parameter entity and a
//   general entity with the same name are two distinct entities. "
// %-refs are top-level only: "In the internal DTD subset, %-refs
//   i.e. parameter-entity references can occur only where markup
//   declarations can occur, not _within_ markup declarations."
// Complicating factor: "(This does not apply to references that
//   occur in external parameter entities or to the external subset.)""
//
// GenEntDef	::=	EntityValue | (ExternalID NDataDecl?)
// NDataDecl	::=	'NDATA' Name
// ParsedEntDef	::=	EntityValue | ExternalID
// EntityValue	::=	'"' ([^%&"] | PEReference | Reference)* '"'
//    "'" ([^%&'] | PEReference | Reference)* "'"
// ExternalID	::=	'SYSTEM' SystemLiteral
//              | 'PUBLIC' PubidLiteral SystemLiteral
// The literals are, basically,
// - double quotes around anything except double quotes, or
// - single quotes around anything except single quotes.
// "If 'NDATA' is present, this is a general unparsed entity;
//  else it is a parsed entity.
// "A [relative] URI might be relative to the document entity,
//  to the entity containing the external DTD subset, or to
//  some other external parameter entity."
//
// "Well-formedness in entities means that XML document's
//  logical and physical structures are properly nested;
//  no start-tag, end-tag, empty-element tag, element,
//  comment, PI, character reference, or entity reference
//  can begin in one entity and end in another."
// "Parameter-entity replacement text must be properly
//  nested with markup declarations. That is to say, if
//  either the first character or the last character of a
//  markup declaration (markupdecl above) is contained in
//  the replacement text for a parameter-entity reference,
//  both must be contained in the same replacement text.""
//
// Document Type Definition
// doctypedecl	::=	'<!DOCTYPE' Name (ExternalID)?
//      ('[' (markupdecl | PEReference)* ']')? '>'
// markupdecl	::=	ELEMENT... | ATTLIST... | ENTITY...
//      | NOTATION... | PI | Comment

// NewEntitiesList collects all entity definitions. -n
// Note that each Token is normalized. -n-
// rtType:ENTITY  string1:foo  string2:"FOO"  entityIsParameter:false -n-
// rtType:ENTITY  string1:bar  string2:"BAR"  entityIsParameter:true
//
// CALLED BY ProcessEntities only//
func (p *Contentity) NewEntitiesList() (gEnts map[string]*gparse.GEnt, err error) {
	// pX := p.TheXml()
	// pXI = new(XmlItems)
	// pXI.GEnts = make(map[string]*gtoken.GEnt)
	gEnts = make(map[string]*gparse.GEnt)

	for _, E := range p.GTags {

		if E.TDType != CT.TD_type_Entity {
			continue
		}
		// fmt.Printf("Collecting ENTITY directive: |%+v| \n", T)
		newEnt := new(gparse.GEnt)

		if E.EntityIsParameter {
			newEnt.TypeIsParm = true
			newEnt.RefChar = "%"
		} else {
			newEnt.RefChar = "&"
		}
		panic("FIXME cnty/xmlprocentities.go:L86")
		theSymbol := E.Text
		theRest := SU.NormalizeWhitespace(E.Text)

		newEnt.NameOnly = theSymbol
		newEnt.TheRest = theRest
		newEnt.NameAsRef = newEnt.RefChar + theSymbol + ";"

		// fmt.Printf("Got an ENTITY: |%s|%s| \n", newEnt.NameAsRef, newEnt)
		// FIXME This must be added AFTER the struct is fully loaded ?
		gEnts[newEnt.NameAsRef] = newEnt

		extIDtype, extIDtext := SU.SplitOffFirstWord(theRest)

		// Check for external entity
		newEnt.IsPublicID = (extIDtype == "PUBLIC")
		newEnt.IsSystemID = (extIDtype == "SYSTEM")
		if !(newEnt.IsPublicID || newEnt.IsSystemID) {
			continue
		}

		/* old code ?
		// e.g. "foo"
		NameOnly string
		// including "%|&" and ";" i.e. "&foo;" or "%foo;"
		NameAsRef string
		// true if parameter entity, false if general entity
		TypeIsParm bool
		// "%" if parameter entity, "&" if general entity
		RefChar    string
		IsSystemID bool
		IsPublicID bool
		// External entities only (PUBLIC, SYSTEM)
		ID  string
		URI string
		// Data model TBD
		Value interface{}
		// This appears to be used to hold "the rest".
		Output string
		*/

		// Got an ExternalID.
		// Two forms:
		// <!ENTITY name SYSTEM "URI">
		// <!ENTITY name PUBLIC "public_ID" "URI">
		// At this point here, "extIDtext" is the last one or two tokens.

		// ======================
		// public_ID: The ID should be in quotes - either single
		// or double - so if we don't find them, it's an error.
		// It would be nice to use Golang's CSV parswer, per
		// https://stackoverflow.com/questions/47489745/
		// splitting-a-string-at-space-except-inside-quotation-marks-go
		// NOTE Unfortunately it only works with double quotes,
		// NOT single quotes !
		// r := csv.NewReader(S.NewReader(extIDtext))
		// r.Comma = ' ' // space
		// fields, err := r.Read()
		// ======================
		var s1, s2 string
		var e error
		if newEnt.IsPublicID {
			s1, s2, e = SU.SplitOffQuotedToken(extIDtext)
			if e != nil {
				/* elog.Printf("PUBLIC external ID not properly quoted: (%s) |%s| \n",
				newEnt.NameAsRef, extIDtext) */
				return nil, fmt.Errorf("Bad quotes on ID for %s ID", extIDtext)
			}
			newEnt.ID = S.TrimSpace(s1)
			extIDtext = S.TrimSpace(s2)
			// NOTE fmt.Printf("PUBLIC:PID> %s \n", newEnt.ID)
		}
		// The URL should be in quotes - either single or
		// double - so if we don't find them, it's an error.
		if !SU.IsXmlQuoted(extIDtext) {
			/* elog.Printf("External ID's URL not properly quoted: (%s) |%s| \n",
			newEnt.NameAsRef, extIDtext) */
			return nil, fmt.Errorf("Bad quotes on URL %s for external ID", extIDtext)
		}
		newEnt.URI = SU.MustXmlUnquote(extIDtext)
		// NOTE fmt.Printf("SYSTEM:URI> %s \n", newEnt.URI)

		// fmt.Printf("==> Pass 4: added ENTITY: %s \n", newEnt)
		// ilog.Printf("Added new ENTITY: \n    %s", newEnt)

		// NOTE fmt.Printf("<<RESOLVE>> %s \n", newEnt.URI)
		// TODO Process search paths
	}
	/* more debugging
	for _,ent := range pRTx.DEnts {
		ilog.Printf("ENTs-DEF'd: %v \n", ent)
	}
	*/
	return gEnts, nil
}

// ====

// DoEntitiesList collects all entity definitions. -n
// Note that each Token has been normalized. -n-
// rtType:ENTITY  string1:foo  string2:"FOO"  entityIsParameter:false -n-
// rtType:ENTITY  string1:bar  string2:"BAR"  entityIsParameter:true
func (p *Contentity) DoEntitiesList() error {
	// pX := p.TheXml()
	println("    ==> DoEntitiesList TODO")
	/* code to use ?
	        if pGF.XmlItems == nil {
			pGF.XmlItems = new(XmlItems)
		}
		var pXI = pGF.XmlItems
	*/
	p.GEnts = make(map[string]*gparse.GEnt)

	for _, E := range p.GTags {

		if E.TDType != CT.TD_type_Entity {
			continue
		}
		fmt.Printf("    --> DoEntitiesList: Collecting directive: |%+v| \n", E)
		newEnt := new(gparse.GEnt)

		if E.EntityIsParameter {
			newEnt.TypeIsParm = true
			newEnt.RefChar = "%"
		} else {
			newEnt.RefChar = "&"
		}
		var theSymbol, theRest string
		panic("cnty/xmlprocentities.go:L214")
		// !! theSymbol = E.TagOrPrcsrDrctv
		// !! theRest = SU.NormalizeWhitespace(E.Datastring)

		newEnt.NameOnly = theSymbol
		newEnt.TheRest = theRest
		newEnt.NameAsRef = newEnt.RefChar + theSymbol + ";"

		// fmt.Printf("Got an ENTITY: |%s|%s| \n", newEnt.NameAsRef, newEnt)
		// FIXME This must be added AFTER the struct is fully loaded ?
		p.GEnts[newEnt.NameAsRef] = newEnt

		extIDtype, extIDtext := SU.SplitOffFirstWord(theRest)

		// Check for external entity
		newEnt.IsPublicID = (extIDtype == "PUBLIC")
		newEnt.IsSystemID = (extIDtype == "SYSTEM")
		if !(newEnt.IsPublicID || newEnt.IsSystemID) {
			continue
		}

		/* old code ?
		// e.g. "foo"
		NameOnly string
		// including "%|&" and ";" i.e. "&foo;" or "%foo;"
		NameAsRef string
		// true if parameter entity, false if general entity
		TypeIsParm bool
		// "%" if parameter entity, "&" if general entity
		RefChar    string
		IsSystemID bool
		IsPublicID bool
		// External entities only (PUBLIC, SYSTEM)
		ID  string
		URI string
		// Data model TBD
		Value interface{}
		// This appears to be used to hold "the rest".
		Output string
		*/

		// Got an ExternalID.
		// Two forms:
		// <!ENTITY name SYSTEM "URI">
		// <!ENTITY name PUBLIC "public_ID" "URI">
		// At this point here, "extIDtext" is the last one or two tokens.

		// ======================
		// public_ID: The ID should be in quotes - either single
		// or double - so if we don't find them, it's an error.
		// It would be nice to use Golang's CSV parswer, per
		// https://stackoverflow.com/questions/47489745/
		// splitting-a-string-at-space-except-inside-quotation-marks-go
		// NOTE Unfortunately it only works with double quotes,
		// NOT single quotes !
		// r := csv.NewReader(S.NewReader(extIDtext))
		// r.Comma = ' ' // space
		// fields, err := r.Read()
		// ======================
		var s1, s2 string
		var e error
		if newEnt.IsPublicID {
			s1, s2, e = SU.SplitOffQuotedToken(extIDtext)
			if e != nil {
				/* elog.Printf("PUBLIC external ID not properly quoted: (%s) |%s| \n",
				newEnt.NameAsRef, extIDtext) */
				return fmt.Errorf("Bad quotes on ID for %s ID", extIDtext)
			}
			newEnt.ID = S.TrimSpace(s1)
			extIDtext = S.TrimSpace(s2)
			// NOTE fmt.Printf("PUBLIC:PID> %s \n", newEnt.ID)
		}
		// The URL should be in quotes - either single or
		// double - so if we don't find them, it's an error.
		if !SU.IsXmlQuoted(extIDtext) {
			/* elog.Printf("External ID's URL not properly quoted: (%s) |%s| \n",
			newEnt.NameAsRef, extIDtext) */
			return fmt.Errorf("Bad quotes on URL %s for external ID", extIDtext)
		}
		newEnt.URI = SU.MustXmlUnquote(extIDtext)
		// NOTE fmt.Printf("SYSTEM:URI> %s \n", newEnt.URI)

		// fmt.Printf("==> Pass 4: added ENTITY: %s \n", newEnt)
		// ilog.Printf("Added new ENTITY: \n    %s", newEnt)

		// NOTE fmt.Printf("<<RESOLVE>> %s \n", newEnt.URI)
		// TODO Process search paths
	}
	/* more debugging
	for _,ent := range pRTx.DEnts {
		ilog.Printf("ENTs-DEF'd: %v \n", ent)
	}
	*/
	return nil
}

// ====

var s2check = ""

// SubstituteEntities does replacement in Entities for simple
// (single-token) entity references, i.e. that begin with "%" or "&".
func (p *Contentity) SubstituteEntities() error {
	// pX := p.TheXml()
	println("    ==> SubstituteEntities TODO")

	var chgs = true
	// var pXI = pGF.XmlItems

	for {
		// fmt.Printf("LOOP \n")
		if !chgs {
			break
		}
		chgs = false

		// TODO

		// First determine the longest sub string, and pick an
		// arbitrary multiple of it (such as 20x) as an absolute
		// upper limit for subs done here. Do this is order to
		// prevent a DOS attack via entity explosion.

		// Then set up two loops that go thru all DEnts.

		// Process all entity definitions and attribute
		// definitions [and also (most importantly?) transclusions?]
		for sEnt, E := range p.GEnts {

			// ilog.Printf("SubEntRksvly: chkg: [%s]%s \n", sEnt, E)
			if sEnt != E.NameAsRef {
				panic("Bad DEnt map keys in SubEntitiesRecursively")
			}
			// if (S.Index(sEnt,"%") == -1 && S.Index(sEnt, "&") == -1) { continue }

			// First let's identify everything that looks like an entity reference.
			for _, E := range p.GTags {
				/* code to use ?
				if RT.rtType == "CD" {
					if RT.string1 != "" { fmt.Printf("CDataDebugNONNIL|%+v| \n", RT) }
				}
				*/
				switch E.TDType {
				case CT.TD_type_ELMNT:
					continue
				case CT.TD_type_ENDLM:
					continue
				case CT.TD_type_PINST:
					continue
				case CT.TD_type_COMNT:
					continue
				case CT.TD_type_DRCTV:
					continue // panic("WTF")
				case CT.TD_type_Doctype:
					continue
				default:
					// CD, ELEMENT, ATTLIST, ENTITY, NOTATION
					if E.TDType == CT.TD_type_CDATA {
						s2check = E.Text
						// if s2check != "" { fmt.Printf("SubEnts got CDATA|%v| \n", RT) }
					} else {
						s2check = E.Text // !! E.Datastring
					}
					// Check for all the legal entity reference characters.
					if -1 == S.IndexAny(s2check, "&%;") {
						continue
					}
					// OK, let's brute force it
					if -1 != S.Index(s2check, sEnt) {
						// FIXME:30 FIXME fmt.Printf("(DD) Got a hit? |%s|%s| \n", sEnt, s2check)
					}
				}
			}

			// Now let's enumerate all the places where entity substitution is kosher.

			/* code to use ?
			chgs = true
			var pid *DPID
			var ok bool
			if pid, ok = pRTx.DPIDs[token]; !ok {
				panic("Entity not defined (yet?): " + token)
			}
			fnam := pid.URI
			ilog.Printf("ent2sub: |%s| <=> <%s> (%v) \n", token, fnam, pid)

			frdr, e := os.Open(fnam)
			defer frdr.Close()
			if e != nil {
				panic("Can't open file: " + e.Error())
			}
			var r io.Reader = frdr
			var transclusion *RichTokenization
			transclusion, e = TokenizeRichly(r)
			// ilog.Println(transcludedTokens)

			// The next token is at iToken, but now we want
			// to insert the slice of tokens we just redd.
			var transcludedTokens []RichToken
			transcludedTokens = transclusion.tokens
			ilog.Printf("iToken<%d:%s> outa <%d> \n",
				iToken, token, len(pRTx.tokens))
			ilog.Printf("iToken<%d> next<%s> all<%d>; inserting<%d>\n",
				iToken, pRTx.tokens[iToken], len(pRTx.tokens), len(transcludedTokens))
			var reddd = pRTx.tokens[:iToken]
			var unred = pRTx.tokens[iToken:]
			pRTx.tokens = append(reddd, transcludedTokens...)
			pRTx.tokens = append(pRTx.tokens, unred...)
			ilog.Printf("iToken<%d> next<%s> all<%d> \n",
				iToken, pRTx.tokens[iToken], len(pRTx.tokens))
			*/
		}
	}
	return nil
}
