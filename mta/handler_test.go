package mta

import (
	"container/list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_61(t *testing.T) {
	a61Line := a61Line()
	fl := FileLines{
		lines:   a61Line,
		current: a61Line.Front(),
		data:    MT940Data{},
	}

	fl.handleLines()
	assert.Equal(t, "190131", fl.data.Transactions[0].Date1, "The Date1 value is not 190131.")
	assert.Equal(t, "0131", fl.data.Transactions[0].Date2, "The Date2 value is not 0131.")
	assert.Equal(t, "C", fl.data.Transactions[0].Mark, "The Mark value is not C.")
	assert.Equal(t, "0,3", fl.data.Transactions[0].Amount, "The Amount value is not 0,3.")
	assert.Equal(t, "NTRF", fl.data.Transactions[0].Code, "The Code value is not NTRF.")
	assert.Equal(t, "NONREF", fl.data.Transactions[0].Reference, "The Reference value is not NONREF.")
}

func Test_86_simple(t *testing.T) {
	a86Lines := a86Line(":86:/IBAN/NL65BUNQ2206724936/NAME/P.C. Wacki/REMI/")
	fl := FileLines{
		lines:   a86Lines,
		current: a86Lines.Front(),
		data:    MT940Data{},
	}

	fl.handleLines()
	assert.Equal(t, "190131", fl.data.Transactions[0].Date1, "The Date1 value is not 190131.")
	assert.Equal(t, "NL65BUNQ2206724936", fl.data.Transactions[0].InformationToAccountOwner.values["IBAN"], "IBAN is incoorect")
	assert.Equal(t, "P.C. Wacki", fl.data.Transactions[0].InformationToAccountOwner.values["NAME"], "IBAN is incoorect")
}

func Test_86_complex1(t *testing.T) {
	a86Lines := a86Line(":86:/IBAN/NL54INGB0006214323/NAME/Hr J F M Wacki en/of Mw K Wacki-Blokland/REMI/")
	fl := FileLines{
		lines:   a86Lines,
		current: a86Lines.Front(),
		data:    MT940Data{},
	}

	fl.handleLines()
	assert.Equal(t, "190131", fl.data.Transactions[0].Date1, "The Date1 value is not 190131.")
	assert.Equal(t, "NL54INGB0006214323", fl.data.Transactions[0].InformationToAccountOwner.values["IBAN"], "IBAN is incoorect")
	assert.Equal(t, "Hr J F M Wacki en/of Mw K Wacki-Blokland", fl.data.Transactions[0].InformationToAccountOwner.values["NAME"], "NAME is incoorect")
}

func Test_86_complex2(t *testing.T) {
	a86Lines := a86Line(":86:/NAME/TEST *SERVICES/REMI/TEST *SERVICES g.co/helppay#, GB")
	fl := FileLines{
		lines:   a86Lines,
		current: a86Lines.Front(),
		data:    MT940Data{},
	}

	fl.handleLines()
	assert.Equal(t, "190131", fl.data.Transactions[0].Date1, "The Date1 value is not 190131.")
	assert.Equal(t, "TEST *SERVICES", fl.data.Transactions[0].InformationToAccountOwner.values["NAME"], "NAME is incoorect")
	assert.Equal(t, "TEST *SERVICES g.co helppay#, GB", fl.data.Transactions[0].InformationToAccountOwner.values["REMI"], "REMI")
}

func Test_62F(t *testing.T) {
	//:62F:C190531EUR5,66
	a62FLines := a62FLine(":62F:C190531EUR5,66")
	fl := FileLines{
		lines:   a62FLines,
		current: a62FLines.Front(),
		data:    MT940Data{},
	}

	fl.handleLines()
	assert.Equal(t, "190531", fl.data.Transactions[0].ClosingBalance.Date, "The Date1 value is not 190531.")
	assert.Equal(t, "C", fl.data.Transactions[0].ClosingBalance.DebitCredit, "The DebitCredit is not C.")
	assert.Equal(t, "EUR", fl.data.Transactions[0].ClosingBalance.Currency, "The Currency is not EUR.")
	assert.Equal(t, "5,66", fl.data.Transactions[0].ClosingBalance.Amount, "The Amount is not 5,66.")
}

func a62FLine(s string) *list.List  {
	l := list.New()
	l.PushBack(":61:1901310131C0,3NTRFNONREF")
	l.PushBack(s)
	return l

}

func a61Line() *list.List {
	l := list.New()
	l.PushBack(":61:1901310131C0,3NTRFNONREF")
	return l
}

func a86Line(s string) *list.List {
	l := list.New()
	l.PushBack(":61:1901310131C0,3NTRFNONREF")
	l.PushBack(s)
	return l
}
