package cnty

import (
	S "strings"

	MU "github.com/fbaube/miscutils"
	XU "github.com/fbaube/xmlutils"
)

// DoValidation TODO If no DOCTYPE, make a guess based on Filext but it can't be fatal.
func (p *Contentity) DoValidation(pXCF *XU.XmlCatalogFile) (dtdS string, docS string, errS string) {
	errS = ""
	if !p.IsXML() {
		panic("DoValidation !IsXML")
	}
	if p.ParsedDoctype == nil {
		return "No_DcTp", "valid??", ""
	}
	var ppid = p.ParsedDoctype.PIDSIDcatalogFileRecord
	// print("\t" + ppid.PTDesc + " --> ")
	val := pXCF.GetByPublicID(ppid.String())
	if val == nil {
		// println("DTD NOT FOUND")
		return "DTD_Unk", "valid??", ""
	}
	// print("DTD OK; ")
	// NOTE At CLI, can use :: alias validate-lw-topic =
	// 'xmllint --noout --dtdvalid file:///LwDTD/lw-topic.dtd'
	// NOTE If we don't specify the DTD, then xmllint only checks for
	// well-formedness, but we've already done this ourselves when
	// building the GTree. So, we should specify the DTD when we invoke
	// xmllint, or else not even bother.
	// NOTE We have to pass the flag "--nowarning", or else
	// validation will fail if the SYSTEM ID can't be found.
	stdOut, stdErr, err := MU.RunCommand(
		"xmllint", "--noout", "--nowarning", "--nonet", "--dtdvalid",
		"file://"+string(val.AbsFilePath), string(p.AbsFilePath))
	// NOTE:1060 that the return value "err" is dumb:
	// it contains stuff like "exit status 3".
	if S.TrimSpace(stdErr) == "" {
		// print("Document is valid \n")
		return "isFound", "isValid", ""
	}
	// print("Validation failed: ")
	if S.TrimSpace(stdOut) != "" {
		errS += stdOut + "\n"
	}
	errS += stdErr
	errS += "==> End of validation errors <== \n"
	if err != nil {
		errS += "==> xmllint: " + err.Error() + "\n"
	}
	return "isFound", "Errors!", errS
}
