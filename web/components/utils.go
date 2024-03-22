package components

import (
	"fmt"

	"github.com/a-h/templ"
)

func mergeClass(class string, attrs templ.Attributes) string {
	class2Any, ok := attrs["class"]
	if !ok {
		class2Any = ""
	}
	class2, ok := class2Any.(string)
	if !ok {
		class2 = ""
	}
	return fmt.Sprintf("%s %s", class, class2)
}
