package got

import (
	"context"
	"fmt"
)

func ExampleContent() {
	status := NewStatusOk()
	ctx := ContextWithContent(context.Background(), status)
	if IsContextContent(ctx) {
		fmt.Printf("Status : %v\n", ContextContent(ctx).(*Status).String())
	}

	//Output:
	//Status : 0 Successful
}
