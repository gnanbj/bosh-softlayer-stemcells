package cleanup_stemcells

import (
	"errors"
	"fmt"
	"time"

	softlayer "github.com/maximilien/softlayer-go/softlayer"

	common "github.com/maximilien/bosh-softlayer-stemcells/common"
)

type CleanupStemcellsCmd struct {
	LastValidDate string
	ShipItTag     string

	client  softlayer.Client
	options common.Options
}

func NewCleanupStemcellsCmd(options common.Options, client softlayer.Client) *CleanupStemcellsCmd {
	return &CleanupStemcellsCmd{
		options: options,
		client:  client,
	}
}

func (cmd *CleanupStemcellsCmd) Println(a ...interface{}) (int, error) {
	return fmt.Println(a)
}

func (cmd *CleanupStemcellsCmd) Printf(msg string, a ...interface{}) (int, error) {
	return fmt.Printf(msg, a)
}

func (cmd *CleanupStemcellsCmd) Options() common.Options {
	return cmd.options
}

func (cmd *CleanupStemcellsCmd) CheckOptions() error {
	if cmd.options.CommandFlag != "cleanup-stemcells" {
		return errors.New(fmt.Sprintf("CommandFlag (%s) is incorrect, is not cleanup-stemcells", cmd.options.CommandFlag))
	}

	if cmd.options.NamePatternFlag == "" {
		return errors.New("NamePatternFlag not passed")
	}

	cmd.defaultOptionalFlags()

	return nil
}

func (cmd *CleanupStemcellsCmd) Run() error {
	return nil
}

// Private

func (cmd *CleanupStemcellsCmd) defaultOptionalFlags() {
	if cmd.options.LastValidDateFlag == "" {
		cmd.LastValidDate = time.Now().Format(time.UnixDate)
	}

	if cmd.options.ShipItTagFlag == "" || cmd.options.ShipItTagFlag == "SHIPIT" {
		cmd.ShipItTag = "SHIPIT"
	}
}
