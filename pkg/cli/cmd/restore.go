package cmd

import (
	"os"
	"path"

	"github.com/aserto-dev/topaz/pkg/cli/cc"
	"github.com/aserto-dev/topaz/pkg/cli/clients"
	"github.com/aserto-dev/topaz/pkg/cli/dockerx"
	"github.com/fatih/color"
	"github.com/google/uuid"
)

type RestoreCmd struct {
	File string `arg:""  default:"backup.tar.gz" help:"absolute file path to local backup tarball"`
	clients.Config
}

func (cmd *RestoreCmd) Run(c *cc.CommonCtx) error {
	if running, err := dockerx.IsRunning(dockerx.Topaz); !running || err != nil {
		if err != nil {
			return err
		}
		color.Yellow("!!! topaz is not running")
		return nil
	}

	cmd.Config.SessionID = uuid.NewString()

	dirClient, err := clients.NewDirectoryClient(c, &cmd.Config)
	if err != nil {
		return err
	}

	if cmd.File == "backup.tar.gz" {
		currentDir, err := os.Getwd()
		if err != nil {
			return err
		}
		cmd.File = path.Join(currentDir, "backup.tar.gz")
	}

	color.Green(">>> restore from %s", cmd.File)
	return dirClient.Restore(c.Context, cmd.File)
}
