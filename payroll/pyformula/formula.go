package pyformula

type Formula struct {
	TenantId       int64
	FormulaId      string
	Country        string
	CustomCodeText string
}

func (formula *Formula) Exec() {

}
