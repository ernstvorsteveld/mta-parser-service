package sta

import (
	"bufio"
	"container/list"
	"fmt"
	"github.com/ernstvorsteveld/mta-common/common"
	"log"
	"os"
	"strings"
)

type MT940Data struct {
	Sender                string
	MessageType           string
	Receiver              string
	TransactionNumber     string
	AccountIdentification string
	StatementNumber       string
	OpeningBalance        AccountBalance
	Transactions          []Transaction
}

type FileLines struct {
	lines   *list.List
	current *list.Element
	data    MT940Data
}

func Start(ch chan common.FilenameMessage) {
	go listener(ch)
}

func listener(ch chan common.FilenameMessage) {
	var c common.FilenameMessage
	i := 0
	for {
		select {
		case c = <-ch:
			log.Print("Message nr: ", i, " message: ", c)
			handle(c)
		}
		i++
	}
}

func handle(c common.FilenameMessage) {
	file, err := os.Open(c.Dst)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	handleFile(file)
}

func handleFile(file *os.File) {
	scanner := bufio.NewScanner(file)
	l := FileLines{
		lines: list.New(),
		data:  MT940Data{},
	}
	for scanner.Scan() {
		l.lines.PushBack(scanner.Text())
	}

	l.handlFileLines()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (f *FileLines) handlFileLines() {
	log.Println("Start handling MT940 lines, number of lines: ", f.lines.Len())

	f.senderLine()
	f.typeLine()
	f.receiverLine()
	f.current = f.current.Next()
	f.handleLines()
}

func (f *FileLines) senderLine() {
	f.current = f.lines.Front()
	f.data.Sender = fmt.Sprintf("%s", f.current.Value)
	log.Println("Sender:", f.data.Sender)
}

func (f *FileLines) typeLine() {
	f.current = f.current.Next()
	f.data.MessageType = fmt.Sprintf("%s", f.current.Value)
	log.Println("Message type:", f.data.MessageType)
}

func (f *FileLines) receiverLine() {
	f.current = f.current.Next()
	f.data.Receiver = fmt.Sprintf("%s", f.current.Value)
	log.Println("Receiver:", f.data.Receiver)
}

func (f *FileLines) handleLines() {
	var cont = "CONTINUE"
	for cont == "CONTINUE" {
		cont = f.handleLine()
		if cont == "CONTINUE" {
			f.current = f.current.Next()
		}
	}
}

func (f *FileLines) handleLine() string {
	if f.current == nil {
		return "EOF"
	}
	v := strings.TrimSpace(fmt.Sprintf("%s", f.current.Value))
	log.Println("Handling line:", v)

	t := getType(v)
	switch t {
	case "20":
		f.data.TransactionNumber = transactionNumber(v)
	case "25":
		f.data.AccountIdentification = acountIdentification(v)
	case "28":
		f.data.StatementNumber = statementNumber(v)
	case "60F":
		f.data.OpeningBalance = openingBalance(v)
	case "61":
		f.data.Transactions = append(f.data.Transactions, transaction(v))
	case "86":
		f.data.Transactions[len(f.data.Transactions)-1].InformationToAccountOwner = informationToAccountOwner(v)
	case "62F":
		f.data.Transactions[len(f.data.Transactions)-1].ClosingBalance = closingBalance(v)
	case "MessageTrailerSection":
		log.Println("MessageTrailerSection")
	case "-":
		return "EOF"
	default:
		return "EOF"
	}
	return "CONTINUE"
}

type AccountBalance struct {
	Mark     string
	Date     string
	Currency string
	Amount   string
}

type Transaction struct {
	Date1                     string
	Date2                     string
	Mark                      string
	Amount                    string
	Code                      string
	Reference                 string
	InformationToAccountOwner InformationToAccountOwner
	ClosingBalance            ClosingBalance
}

type InformationToAccountOwner struct {
	values map[string]string
}

type ClosingBalance struct {
	DebitCredit string
	Date        string
	Currency    string
	Amount      string
}

func transaction(v string) Transaction {
	//Optional and provided in format
	//6!n[4!n]2a[1!a]15d1!a3!c16x[//16x]
	//[34x] - The last element supplementary details must be on a new line.
	// date, mark, amount,
	// 1901310131C0,3NTRFNONREF
	// date1: 190131
	// date2 : 0131
	// mark: C
	// amount: 0,3
	// transaction code: NTRF
	// transaction reference: NONREF

	s := getValue("61:Transactions:", v)
	amount, pos := transactionAmount(s[11:])
	t := Transaction{
		Date1:     s[0:6],
		Date2:     s[6:10],
		Mark:      s[10:11],
		Amount:    amount,
		Code:      s[11+pos : 11+pos+4],
		Reference: s[11+pos+4:],
	}
	log.Println(t)
	return t
}

func informationToAccountOwner(v string) InformationToAccountOwner {
	// :86:/IBAN/NL65BUNQ2206724936/NAME/P.C. Wacki/REMI/
	itao := InformationToAccountOwner{
		values: make(map[string]string),
	}

	containsEnOf := strings.Index(v, "en/of") > 0
	parts := strings.Split(v, "/")
	log.Println(parts)
	parts = mergeEnOf(parts, containsEnOf)

	remiFound := false

	for i := 0; i < len(parts); i++ {
		// parts[0] = :86:
		if i == 0 || i%2 == 0 {
			continue
		}
		if remiFound {
			itao.values["REMI"] = itao.values["REMI"] + " " + parts[i]
			continue
		}
		if parts[i] == "REMI" {
			remiFound = true
		}
		itao.values[parts[i]] = parts[i+1]
	}

	return itao
}

func mergeEnOf(parts []string, containsEnOf bool) []string {
	if !containsEnOf {
		return parts
	}
	newparts := make([]string, 0)
	enFound := false
	for _, s := range parts {
		if enFound {
			newparts[len(newparts)-1] = newparts[len(newparts)-1] + "/" + s
			enFound = false
		} else {
			if strings.HasSuffix(s, "en") {
				enFound = true
			}
			newparts = append(newparts, s)
		}
	}
	return newparts
}

func transactionAmount(s string) (string, int) {
	i := strings.Index(s, "N")
	return s[0:i], i
}

func closingBalance(v string) ClosingBalance {
	s := getValue("62F:", v)
	t := ClosingBalance{
		DebitCredit: s[0:1],
		Date:        s[1:7],
		Currency:    s[7:10],
		Amount:      s[10:],
	}
	return t
}

func openingBalance(v string) AccountBalance {
	//Mandatory and of format 1!a6!n3!a15d (D/C Mark)(Date)(Currency)(Amount).
	//In D/C Mark, C for Credit Balance and D for Debit balance.
	s := getValue("60F:Opening Balance:", v)

	ab := AccountBalance{
		Mark:     s[0:1],
		Date:     s[1:7],
		Currency: s[7:10],
		Amount:   s[10:],
	}
	log.Println(ab)
	return ab
}

func statementNumber(v string) string {
	return getValue("Statement number:", v)
}

func acountIdentification(v string) string {
	return getValue("Account identification:", v)
}

func transactionNumber(v string) string {
	return getValue("Transactions number:", v)
}

func getValue(t string, v string) string {
	s := strings.Split(v, ":")
	log.Println(t, s[2])
	return s[2]
}

func getType(v string) string {
	s := strings.Split(v, ":")
	if len(s) > 1 {
		log.Println("Statement type:", s[1])
		return s[1]
	} else if strings.Index(v, "-") == 0 {
		return "-"
	} else {
		return "MessageTrailerSection"
	}
}
