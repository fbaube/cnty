package cnty

import (
	"github.com/nbio/xml"
	"fmt"
	"sort"
	S "strings"

	CT "github.com/fbaube/ctoken"
	"github.com/fbaube/gtoken"
	// XU "github.com/fbaube/xmlutils"
)

type StringTally map[string]int

// func (tt TagTally) String() string {
//	return "tt"
// }

var GlobalTagTally StringTally
var GlobalAttTally StringTally
var GlobalTagCount int
var GlobalAttCount int

func init() {
	GlobalTagTally = make(map[string]int)
	GlobalAttTally = make(map[string]int)
}

func (p *Contentity) TallyTags() {
	p.TagTally = make(map[string]int)
	p.AttTally = make(map[string]int)
	for _, pGT := range p.GTokens {
		if pGT.CName.Local == "" {
			continue
		}
		AddInXName(p.TagTally, p.AttTally, pGT)
		AddInXName(GlobalTagTally, GlobalAttTally, pGT)
		GlobalTagCount++
		GlobalAttCount += len(pGT.CAtts)
	}
}

func AddInXName(ElmT StringTally, AttT StringTally, gT *gtoken.GToken) {
	// Is there an entry for the tag yet ?
	// if val, ok := dict["foo"]; ok {
	var ok bool
	var n int
	if n, ok = ElmT[gT.CName.Local]; ok {
		ElmT[gT.CName.Local] = n + 1
	} else {
		ElmT[gT.CName.Local] = 1
	}
	// Now process the attributes
	for _, A := range gT.CAtts { // A is a *XAtt i.e. *xml.Attr
		var gat CT.CAtt = A // *A
		var xat = xml.Attr(gat)
		var sat = xat.Name.Local
		if n, ok = AttT[sat]; ok {
			AttT[sat] = n + 1
		} else {
			AttT[sat] = 1
		}
	}
}

func (st StringTally) StringSortedValues() string {
	n := map[int][]string{}
	var a []int
	for k, v := range st {
		n[v] = append(n[v], k)
	}
	for k := range n {
		a = append(a, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(a)))
	var sb S.Builder
	for _, k := range a {
		for _, s := range n[k] {
			sb.WriteString(fmt.Sprintf("%s:%d,", s, k))
		}
	}
	return sb.String()
}
