package enum

import (
	"fmt"
	"testing"
)

func TestEnum(t *testing.T) {
	t.Run("Just testing enums here", func(t *testing.T) {
		fmt.Println(IntFriday)
		fmt.Println(StrFriday)
		fmt.Println(DayOfWeekFriday.String())
	})
}
