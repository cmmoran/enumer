package main

// Arguments to format are:
//
//	[1]: type name
const valueMethod = `func (i %[1]s) Value() (driver.Value, error) {
	return i.String(), nil
}
`

const valueMethodInt = `func (i %[1]s) Value() (driver.Value, error) {
	return int64(i), nil
}
`

const scanMethod = `func (i *%[1]s) Scan(value any) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case int64:
		*i = %[1]s(v)
		return nil
	case int32:
		*i = %[1]s(v)
		return nil
	case int:
		*i = %[1]s(v)
		return nil
	case float64:
		*i = %[1]s(int64(v))
		return nil
	case float32:
		*i = %[1]s(int64(v))
		return nil
	}

	var str string
	switch v := value.(type) {
	case []byte:
		str = string(v)
	case string:
		str = v
	case fmt.Stringer:
		str = v.String()
	default:
		return fmt.Errorf("invalid value of %[1]s: %%[1]T(%%[1]v)", value)
	}

	if n, err := strconv.ParseInt(strings.TrimSpace(str), 10, 64); err == nil {
		*i = %[1]s(n)
		return nil
	}

	val, err := %[1]sString(str)
	if err != nil {
		return err
	}
	*i = val
	return nil
}`

func (g *Generator) addValueAndScanMethod(typeName string, returnInt bool) {
	g.Printf("\n")
	var vMethod string
	if returnInt {
		vMethod = valueMethodInt
	} else {
		vMethod = valueMethod
	}
	g.Printf(vMethod, typeName)
	g.Printf("\n\n")
	g.Printf(scanMethod, typeName)
}
