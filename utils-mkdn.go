package cnty

import (
	"fmt"

	AST "github.com/yuin/goldmark/ast"
)

// =================
// ==== BF node ====
// == (BF library) =
// =================

func KidsAsSlice(p AST.Node) []AST.Node {
	var pp []AST.Node
	c := p.FirstChild()
	for c != nil {
		pp = append(pp, c)
		c = c.NextSibling()
	}
	return pp
}

func ListKids(p AST.Node) string {
	var pp []AST.Node
	pp = KidsAsSlice(p)
	if len(pp) == 0 {
		return "<0-kids>"
	}
	var s string
	for i, c := range pp {
		s += fmt.Sprintf("[%d:%s]", i, c.Type)
	}
	return s
}

// ================
// ===== GElm =====
// ================

func DumpGElm(p AST.Node) string {
	var s string
	/* code to use ?
	s = fmt.Sprintf("BFnode<%s>: ", gparse.MDnodeType[p.Type])
	s += "<|" + SU.NormalizeWhitespace(string(p.Literal)) + "|> "
	s += gparse.DumpHdg(p.HeadingData)
	s += gparse.DumpList(p.ListData)
	s += gparse.DumpCdBlk(p.CodeBlockData)
	s += gparse.DumpLink(p.LinkData)
	s += gparse.DumpTableCell(p.TableCellData)
	*/
	return s
}

// type myGElmVisitor func(node *Node, entering bool) WalkStatus
func myGElmVisitor(N AST.Node, entering bool) AST.WalkStatus {
	if !entering {
		return AST.WalkContinue // GoToNext
	}
	fmt.Printf("%s \n", DumpGElm(N)) // SU.GetIndent(lvl)
	return AST.WalkContinue          // GoToNext
}
