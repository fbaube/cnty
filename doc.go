// Package cnty defines a per-file structure [Contentity] that
// holds all relevant per-file information. This includes:
//
//   - file path info
//   - file content (UTF-8, tietysti)
//   - file type information (MIME and more)
//   - the results of markup-specific file analysis (in the most analysable
//     case, i.e. XML, this comprises tokens, gtokens, gelms, gtree)
//
// For a discussion of tree walk functions, see `doc_wfn.go`
//
// Note that if we do not get an explicit XML DOCTYPE
// declaration, there is some educated guesswork required.
//
// The first workflow was based on XML, and comprises:
// `text => XML tokens => GTokens => GTags => GTree`
//
// First, package `gparse` gets as far as the `GToken`s,
// which can only be in a list: they have no tree structure.
// Then package `gtree` handles the rest.
//
// XML analysis starts off with tokenization (by the stdlib),
// so it makes sense to then have separate steps for making
// `GToken's, GTag's, GTree`. <br/>
// MKDN and HTML analyses use higher-level libraries that
// deliver CSTs (Concrete Syntax Tree, i.e. parse tree).
// We choose to do this processing in `package gparse`
// rather than in `package gtree`.
//
// MKDN gets a tree of `yuin/goldmark/ast/Node`, and HTML
// gets a tree of stdlib `golang.org/x/net/html/Node`.
// Since a CST is delivered fully-formed, it makes sense
// to have Step 1 that attaches to each node its `GToken´
// and `GTag`, and then Step 2 that builds a `GTree`.
//
// There are three major types of `Contentity`,
// corresponding to how we process the file content:
// - "XML"
// - - (§1) Use stdlib `encoding/xml` to get `[]XU.XToken`
// - - (§1) Convert `[]XU.XToken` to `[]gparse.GToken`
// - - (§2) Build `GTree`
// - "MKDN"
// - - (§1) Use `yuin/goldmark` to get tree of `yuin/goldmark/ast/Node`
// - - (§1) From each Node make a `MkdnToken` (in a list?) incl. `GToken` and `GTag`
// - - (§2) Build `GTree`
// - "HTML"
// - - (§1) Use `golang.org/x/net/html` to get a tree of `html.Node`
// - - (§1) From each Node make a `HtmlToken` (in a list?) incl. `GToken` and `GTag`
// - - (§2) Build `GTree`
//
// In general, all go files in this protocol stack should be organised as: <br/>
// - struct definition()
// - constructors (named `New*`)
// - printf stuff (Raw(), Echo(), String())
//
// Some characteristic methods:
// - Raw() returns the original string passed from the golang XML parser
// (with whitespace trimmed)
// - Echo() returns a string of the item in normalised form, altho be
// aware that the presence of terminating newlines is not treated uniformly
// - String() returns a string suitable for runtime nonitoring and debugging
//
// NOTE The use of shorthand in variable names: Doc, Elm, Att.
//
// NOTE We use `godoc2md`, so we can use Markdown in these code comments.
package cnty
