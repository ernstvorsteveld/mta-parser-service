package mta

import (
	"fmt"
	"regexp"
	"testing"
)

func Test_61_regex(t *testing.T) {
	str := `:61:1901310131C2,NTRFNONREF`
	r, _ := regexp.Compile(`^:61:(.{6})(.{4})(C|D)(\d*,\d*)(.{4})(.*)`)
	//print(str, *r)

	str = `:86:/IBAN/NL65BUNQ2206724936/NAME/P.C. Wacki/REMI/`
	r, _ = regexp.Compile(`(^:86:)(/)([a-zA-Z0-9\.\s]+)`)
	print(str,*r)
}

func print(str string, r regexp.Regexp) {
	for index, match := range r.FindStringSubmatch(str) {
		fmt.Printf("[%d] %s\n", index, match)
	}
}