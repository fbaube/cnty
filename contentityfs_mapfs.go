package cnty

// Is this useful ?
/*
https://pkg.go.dev/testing/fstest@master#MapFS

A MapFS is a simple in-memory file system for use in tests,
represented as a map from path names (arguments to Open) to
information about the files or directories they represent.

The map need not include parent directories for files contained in
the map; those will be synthesized if needed. But a directory can
still be included by setting the [MapFile.Mode]'s fs.ModeDir bit;
this may be necessary for detailed control over the directory's
fs.FileInfo or to create an empty directory.

File system operations read directly from the map, so that the file
system can be changed by editing the map as needed. An implication is
that file system operations must not run concurrently with changes to
the map, which would be a race. Another implication is that opening or
reading a directory requires iterating over the entire map, so a MapFS
should typically be used with not more than a few hundred entries or
directory reads.

https://pkg.go.dev/testing/fstest@master#MapFile

type MapFile struct {
	Data    []byte      // file content
	Mode    fs.FileMode // fs.FileInfo.Mode
	ModTime time.Time   // fs.FileInfo.ModTime
	Sys     any         // fs.FileInfo.Sys
}

*/
