package cnty

import (
	"fmt"

	LU "github.com/fbaube/logutils"
	L "github.com/fbaube/mlog"
)

func (p *Contentity) LogPrefix(mid string) string {
	return fmt.Sprintf("%02d%sst%s", p.Lindex, mid, p.Lstage)
}

func (p *Contentity) L(level LL, format string, a ...interface{}) {
	// p.L(LLog(level, format, a...)
	L.L.LogWithString(LU.Level(level), format,
		fmt.Sprintf("F%02d", p.Lindex)+"|stg"+p.Lstage, a...)
}

/*
func (p *Contentity) LogTextQuote(level LL, textquote string, format string, a ...interface{}) {
	// p.L(LLog(level, format, a...)
	p.L(LLogWithString(LU.Level(level), format,
		fmt.Sprintf("%02d", p.Lindex)+","+p.Lstage, a...)
	panic("FIXME")
	// p.L(LLogMultilineAsIs(SU.IndentWith("   |  ", textquote))
}
*/

type LL LU.Level

var LDebug, LInfo, LOkay, LWarning, LError, LPanic LL

func init() {
	LDebug = LL(LU.LevelDebug)
	// LProgress = LL(LU.LevelProgress)
	LInfo = LL(LU.LevelInfo)
	LOkay = LL(LU.LevelOkay)
	LWarning = LL(LU.LevelWarning)
	LError = LL(LU.LevelError)
	LPanic = LL(LU.LevelPanic)
}
