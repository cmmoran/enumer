// Code generated by "enumer -type Pill -json"; DO NOT EDIT.

package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _PillName = "PlaceboAspirinIbuprofenParacetamol"

var _PillIndex = [...]uint8{0, 7, 14, 23, 34}

const _PillLowerName = "placeboaspirinibuprofenparacetamol"

func (i Pill) String() string {
	if i < 0 || i >= Pill(len(_PillIndex)-1) {
		return fmt.Sprintf("Pill(%d)", i)
	}
	return _PillName[_PillIndex[i]:_PillIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _PillNoOp() {
	var x [1]struct{}
	_ = x[Placebo-(0)]
	_ = x[Aspirin-(1)]
	_ = x[Ibuprofen-(2)]
	_ = x[Paracetamol-(3)]
}

var _PillValues = []Pill{Placebo, Aspirin, Ibuprofen, Paracetamol}

var _PillNameToValueMap = map[string]Pill{
	_PillName[0:7]:        Placebo,
	_PillLowerName[0:7]:   Placebo,
	_PillName[7:14]:       Aspirin,
	_PillLowerName[7:14]:  Aspirin,
	_PillName[14:23]:      Ibuprofen,
	_PillLowerName[14:23]: Ibuprofen,
	_PillName[23:34]:      Paracetamol,
	_PillLowerName[23:34]: Paracetamol,
}

var _PillNames = []string{
	_PillName[0:7],
	_PillName[7:14],
	_PillName[14:23],
	_PillName[23:34],
}

// PillString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func PillString(s string) (Pill, error) {
	if val, ok := _PillNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _PillNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Pill values", s)
}

// PillValues returns all values of the enum
func PillValues() []Pill {
	return _PillValues
}

// PillStrings returns a slice of all String values of the enum
func PillStrings() []string {
	strs := make([]string, len(_PillNames))
	copy(strs, _PillNames)
	return strs
}

// IsAPill returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Pill) IsAPill() bool {
	for _, v := range _PillValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for Pill
func (i Pill) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Pill
func (i *Pill) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Pill should be a string, got %s", data)
	}

	var err error
	*i, err = PillString(s)
	return err
}
