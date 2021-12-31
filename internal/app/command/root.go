package command

import (
	"context"
	"fmt"
)

func (c *Command) Root(ctx context.Context) {
	fmt.Println("It works")
}
