package account

import (
	"fmt"
	"io"
)

func updateSetOneField(q io.StringWriter, fieldName string, value interface{}, n *int, args *[]interface{}) error {
	if *n != 1 {
		if _, err := q.WriteString(`, `); err != nil {
			return err
		}
	}

	if _, err := q.WriteString(fmt.Sprintf(`%s = $%d`, fieldName, *n)); err != nil {
		return err
	}
	*n++

	*args = append(*args, value)

	return nil
}
