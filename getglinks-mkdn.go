package cnty

// GatherLinks is:
// @conref to reuse block-level content,
// @keyref to reuse phrase-level content.
// TODO Each type of link (i.e. elm/att where it occurs) has to be categorised.
// TODO Each format of link target has to be categorised.
// Cross ref : <xref> : <a href> : [link](/URI "title")
// Key def : <keydef> : <div data-class="keydef"> : <div data- class="keydef"> in HDITA syntax
// Map : <map> : <nav> : See Example of an MDITA map (20)
// Topic ref : <topicref> : <a href> inside a <li> : [link](/URI "title") inside a list item
// TODO Stuff to get:
// XDITA map
// - topicref @href (w @format)
// - task @id
// HDITA
// - article @id
// - span @data-keyref
// - p @data-conref
// MDITA
// - has YAML "id"
// - uses <p @data-conref>
// - uses <span @data-keyref>
// - uses MD [link_text](link_target.dita)
// - uses ![The remote](../images/remote-control-callouts.png "The remote")
// XDITA
// - topic @id
// - ph @keyref
// - image @href
// - p @id
// - video/source @value
// - section @id
// - p @conref
func (p *Contentity) GatherLinks() error {
	println("    --> MF.GatherLinks TODO")
	/* code to use ?
	// if !pGF.IsXML { return nil }
	if pGF.Micodo == nil || pGF.Micodo[0] != "lwdita" {
		return nil
	}
	*/
	return nil
}
