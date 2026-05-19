package cnty

import (
	// BF "github.com/fbaube/blackfriday" // Version 2 !!
	// BF "github.com/yuin/goldmark"
	AST "github.com/yuin/goldmark/ast"
)

var unlinkList []AST.Node

func NormalizeTextLeaves(rootNode AST.Node) {
	unlinkList = make([]AST.Node, 0, 64)
	// rootNode.Walk(func(node AST.Node, entering bool) AST.WalkStatus {
	AST.Walk(rootNode, func(node AST.Node, entering bool) (AST.WalkStatus, error) {
		/* code to use ?
		if string(node.Literal) == "" {
			return AST.WalkContinue // GoToNext
		}
		if node.Type != BF.Text {
			return AST.WalkContinue // GoToNext
		}
		*/
		var parent = node.Parent()
		/* code to use ?
		if parent.Type == BF.Paragraph {
			return AST.WalkContinue // GoToNext
		}
		if string(parent.Literal) != "" {
			// panic("Literal overload")
		}
		parent.Literal = node.Literal
		*/
		if node == parent.FirstChild() && node == parent.LastChild() {
			unlinkList = append(unlinkList, node)
			// node.Unlink()
			return AST.WalkContinue, nil // GoToNext
		}
		if node == parent.FirstChild() || node == parent.LastChild() {
			unlinkList = append(unlinkList, node)
			// node.Unlink()
			return AST.WalkContinue, nil // GoToNext
		}
		println("OOPS")

		return AST.WalkContinue, nil // GoToNext
	})
	for _, n := range unlinkList {
		// RemoveChild(self, child Node)
		// n.Unlink()
		n.Parent().RemoveChild(n.Parent(), n)
	}
}
