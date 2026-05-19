package cnty

import (
       // "log"
	"runtime/debug"
	"fmt"
	"os"
	// "github.com/fbaube/must"
	SU "github.com/fbaube/stringutils"
)

// ExecuteStages processes a Contentity to completion in an isolated
// thread, and can eaily be converted to run as a goroutine. Summary:
//   - st0_Init()
//   - st1_Read()
//   - st2_Tree()
//   - st3_Refs()
//   - st4_Done() (not currently called, but
//     will work on all input files at once !)
//
// An interesting question is, how can we indicate an error and
// terminate a thread prematurely ? The method currently chosen
// is to use interface [github.com/fbaube/miscutils/Errer]. This
// has to be checked for at the start of a func. But then we can
// chain functions by writing them left-to-right. Winning!
//
// (If functions accept and return a ptr+error pair then they
// chain right-to-left, which is a big fail for readability.) 
//
// We could also pass in a `Context` and use its cancellation
// capability. Yet another way might be simply to `panic`,
// and so this function already has code to catch panics.
// .
func (p *Contentity) ExecuteStages() *Contentity {

     // Start by filtering out pathological
     // and uninteresting cases. 
	if p.RawType() == "" { // or SU.MU_type_UNK {
		panic("UNK MarkupType in ExecuteStages")
	}
	if p.FSO.TypedRaw == nil ||
	  (p.FSO.Raw == "" && p.RawType() != SU.Raw_type_DIRLIKE) {
		p.L(LWarning, "No content: skipping")
		return p
	}
	if p.HasError() {
		p.L(LInfo, "Has error: skipping")
		return p
	}
	if p.RawType() == SU.Raw_type_BIN {
		p.L(LWarning, "Binary file: skipping")
		return p
	}
	// p.L(LInfo, "LENGTH %d SIZE %d", len(p.FSObject.Raw), p.Size())
	if p.FSO.FPs.IsDirlike {
		p.L(LInfo, "Is dir (or similar): skipping")
		return p
	}
	p.Lstage = "--"
	defer func() {
		if r := recover(); r != nil {
		     	fmt.Fprintf(os.Stderr, "RECOVER: %v \n", r) 
			p.L(LPanic, "= = = = = = = = = = = = = = = = = = = =")
			p.L(LPanic, " ** PANIC caught in ExecuteStages ** ")
			var sRecovered string
			var eRecovered error
			p.L(LInfo, "recover() got a: %T", r)
			switch r.(type) {
			case string:
				sRecovered = r.(string)
			case error:
				eRecovered = r.(error)
			}
			if sRecovered == "" {
				sRecovered = eRecovered.Error()
			}
			p.L(LError, "defer'd-recover()-string: "+sRecovered)
			debug.PrintStack()
			/* Barfed!
			bb := debug.Stack()
			ss := string(bb)
			p.L(LError, "stdlib:runtime/debug.Stack(): "+
				ss) // string(debug.Stack()))
			*/
			p.L(LError, "= = = = = = = = = = = = = = = = = = = =")
		}
	}()
	// Execute stages/steps
	return p.
		st0_Init().
		st1_Read().
		st2_Tree().
		st3_Refs()
	// Stage/Step 4 works on all input files at once !
	// st4_Done()
}
