package cnty

// ContentityFactory tracks the state of a ContentityFS
// tree being assembled, for example when a directory
// is specified for recursive analysis.
// 
// FIXME: ID assignment should be offloaded to the DB ?
// .
type ContentityFactory struct {
	// nexSeqID should be reset to 0 when starting another tree ?
	// No, because every single entity (dir/file) gets one,
	// even if it is listed on the CLI as an individual file.
	// But maybe these should be handed out by SQLite.
	// nexSeqID      int
	rootPath      string
}



