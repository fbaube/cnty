package cnty

import (
	"errors"
	"fmt"
	"io/fs"
)

// ContentityError is
// Contentity +
// SrcLoc (in source code) +
// PathError struct { Op, Path string; Err error }
//
// # Maybe use the format pkg.filename.methodname.Lnn
// .
type ContentityError struct {
	PE fs.PathError
	*Contentity
}

func WrapAsContentityError(e error, op string, cty *Contentity) ContentityError {
	ce := ContentityError{}
	ce.PE.Err = e
	ce.PE.Op = op
	if cty == nil {
		ce.PE.Path = "(contentity path not found!)"
	} else {
		ce.PE.Path = cty.FSO.FPs.AbsFP
	}
	return ce
}

func NewContentityError(ermsg string, op string, cty *Contentity) ContentityError {
	ce := ContentityError{}
	ce.PE.Err = errors.New(ermsg)
	ce.PE.Op = op
	if cty == nil {
		ce.PE.Path = "(contentity path not found!)"
	} else {
		ce.PE.Path = cty.FSO.FPs.AbsFP
	}
	return ce
}

func (ce ContentityError) Error() string {
	return ce.String()
}

func (ce *ContentityError) String() string {
	var s string
	s = fmt.Sprintf("%s(%s): %s", ce.PE.Op, ce.PE.Path, ce.PE.Err.Error())
	return s
}
