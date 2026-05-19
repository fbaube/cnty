// This file discusses walk functions, which are obviously important
// when dealing with tree structures.
//
// In file names and function names, `wfn` is used as shorthand for
// "walk function".
//
// It is important to make this code as generic as possible, so that
// additional content formats can easily be added. Therefore in some
// cases, function prototypes (i.e. `type myfunc func etc...`) will
// be used, even if they cannot be specific about the types involved,
// and can only refer to a generalized tree node type or a generalized
// array of tree node types.
//
// In fact, `interface gparse.MarkupStringer` serves this purpose.)
//
// In Go programming, the most common model is the one used by `package
// path/filepath`. The type signatures are as follows:
//
//   func Walk(root string, walkFn WalkFunc) error
//
//   type WalkFunc func(path string, info os.FileInfo, err error) error
//
// In the stdlib it is about walking a file tree, so errors have to
// bubble up, and sentinel value `SkipDir` says stop processing the
// current directory, i.e. skip the subtree. 
//
// In our generalized case, we can skip the complexity of SkipDir, but 
// we do need to bubble up errors. So in accordance with the discussion 
// above, we fill define types that describe generalized functions.
//
// In the stdlib it looks like this:
// https://golang.org/src/path/filepath/path.go
//
//  Walk walks the file tree rooted at root, calling walkFn
//  for each file or directory in the tree, including root.
//  func Walk(root string, walkFn WalkFunc) error {
// 	 info, err := os.Lstat(root)
//	 if err != nil {
//		err = walkFn(root, nil, err)
//	 } else {
//		err = walk(root, info, walkFn)
//	 }
//	 if err == SkipDir {
//		return nil
//	 }
//	 return err
//  }
//
// If we combine these signatures with the actual source code, we get a
// guide to the best practice. (We say this because we're always sposta
// look to the standard library for guidance.)
//
// File `doc.go` outlines the processing for types "XML", "MKDN", and "HTML".
// In each case, processing depends on the libraries being used. Some libraries
// do not perform explicit tokenization: instead they return a CST (Concrete
// Syntax Tree, i.e. parse tree), which we typically convert to an AST,
// abstract syntax tree.
//
// This distinction (tokenization versus CST) determines how we perform
// MCFile Stage 1, `Read(..)`. Therefore let us try to generalize the
// calling signatures used in Stage 1.
//
// (Note that in MCFile Stage 2 `Tree(..)`, both API styles build
// a `GTree` from `GTokens`s, but when the CST already exists, its
// hierarchy information can be used to help build the `GTree`.)
//
// - "XML"
// - - TextTokensFromString  s => `encoding/xml` => `XU.XToken`
// - - GTokensFromTextTokens `[]XU.XToken` => `[]gparse.GToken`
// - "MKDN"
// - - CSTfromString         s => `yuin/goldmark` => `(gm)/ast/Node`-tree
// - - TreeTokensFromCST     Node-tree => []`MkdnToken`
// - - GTokensFromTreeTokens []`MkdnToken` => `GToken`s & `GTag`s
// - "HTML"
// - - CSTfromString         s => `golang.org/x/net/html` => `html.Node`-tree
// - - TreeTokensFromCST     Node-tree => []`HtmlToken`
// - - GTokensFromTreeTokens []`HtmlToken` => `GToken`s & `GTag`s
//
// In this initial version, all these data structures are embedded in `MCFile`.
// In the future maybe they can be connected more directly, using some sort
// of typed pipes.
//
// Suggested signatures (GM = `yuin/goldmark`; not incl. `error` return values):
//  (Path when there is an explicit tokenization step)
//   func TokensFromString_xml(s string) ([]XU.XToken)
//        string => []XU.XToken
//        scalar => array
//  (Path when an CST is created directly from the input string)
//   func CSTfromString_(mkdn|html)(s string) ([]theType)
//        string => []GM/ast/Node, []golang.org/x/net/html/Node
//        scalar => tree
//   func TreeTokensFromCST_notXml(interface{}) ([]*MarkupStringer)
//        GM/ast/Node, golang.org/x/net/html/Node => []MkdnToken, []HtmlToken
//        tree => array
//  (Path that is common) (unless having TreeTokens demands a separate func)
//   func GTokensFromBaseTokens(interface{}) ([]*GToken)
//        XU.XToken, [][]MkdnToken, []HtmlToken => []GToken
//        array => array
//   func GTreeFromGTokens([]GToken) GTree
//        []GToken => GTree
//        array => tree
//
// Summary of dimensionality:
// - With explicit tokenization:
// - - (st1) scalar => array  => array   => (st2) tree
// - - (st1) string => texTox => GTokens => (st2) GTree
// - With CST directly:
// - - (st1) scalar => tree => array  => array   => (st2) tree
// - - (st1) string => CST  => triTox => GTokens => (st2) tree

package cnty
