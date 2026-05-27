package cnty

import (
	"fmt"
	// L "github.com/fbaube/mlog"
)

/*
func (p *Contentity) SetError(s string) {
     	var e error
	e = fmt.Errorf("[F%02d:S%s] %s", p.logIdx, p.logStg, s)
	p.Errer.Err = e
	p.L(LError(e.Error())
}
*/

func (p *Contentity) WrapError(s string, e error) {
     	var e2 error
	e2 = fmt.Errorf("[F%02d:S%s] %s: %w", p.Lindex, p.Lstage, s, e)
	p.Errer.Err = e2
	p.L(LError, e2.Error())
}
