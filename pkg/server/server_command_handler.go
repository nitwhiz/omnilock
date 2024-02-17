package server

import "github.com/nitwhiz/omnilock/v2/pkg/client"

type CommandHandler func(c *client.Client, argv ...string) (result bool, err error)
