package command

import (
	"context"
	"fmt"
)

func (c *Command) Root(ctx context.Context) error {
	fmt.Println("It works")
	return nil
}
