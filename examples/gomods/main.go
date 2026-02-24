package main

import "fmt"

//go:generate go run github.com/cmmoran/enumer -type=Pill -json
type Pill int

const (
	Placebo       Pill = iota
	Aspirin            //enumer:alias=bayer
	Ibuprofen          //enumer:alias=advil
	Paracetamol        //enumer:alias=tylenol
	Acetaminophen = Paracetamol
)

func main() {
	fmt.Println(PillStrings())
	fmt.Println(Placebo.IsAPill())
	fmt.Println(Placebo)
	asprinAlias, _ := PillString("bayer")
	fmt.Println("bayer =", asprinAlias)
	ibuprofenAlias, _ := PillString("advil")
	fmt.Println("advil =", ibuprofenAlias)
	paracetamolAlias, _ := PillString("tylenol")
	fmt.Println("tylenol =", paracetamolAlias)
}
