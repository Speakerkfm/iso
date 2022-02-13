package command

import (
	"context"
	"log"
)

func (c *Command) Root(ctx context.Context) error {
	log.Println("It works")
	return nil
}
