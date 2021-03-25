// +build amd64,linux amd64,darwin arm64,darwin

package machine

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/containers/common/pkg/completion"
	"github.com/containers/podman/v3/cmd/podman/registry"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/containers/podman/v3/pkg/machine"
	"github.com/containers/podman/v3/pkg/machine/qemu"
	"github.com/spf13/cobra"
)

var (
	removeCmd = &cobra.Command{
		Use:               "remove [options] NAME",
		Short:             "Remove an existing machine",
		Long:              "Remove an existing machine ",
		RunE:              remove,
		Args:              cobra.ExactArgs(1),
		Example:           `podman machine remove myvm`,
		ValidArgsFunction: completion.AutocompleteNone,
	}
)

var (
	destoryOptions machine.RemoveOptions
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Mode:    []entities.EngineMode{entities.ABIMode, entities.TunnelMode},
		Command: removeCmd,
		Parent:  machineCmd,
	})

	flags := removeCmd.Flags()
	formatFlagName := "force"
	flags.BoolVar(&destoryOptions.Force, formatFlagName, false, "Do not prompt before removeing")

	keysFlagName := "save-keys"
	flags.BoolVar(&destoryOptions.SaveKeys, keysFlagName, false, "Do not delete SSH keys")

	ignitionFlagName := "save-ignition"
	flags.BoolVar(&destoryOptions.SaveIgnition, ignitionFlagName, false, "Do not delete ignition file")

	imageFlagName := "save-image"
	flags.BoolVar(&destoryOptions.SaveImage, imageFlagName, false, "Do not delete the image file")
}

func remove(cmd *cobra.Command, args []string) error {
	var (
		err    error
		vm     machine.VM
		vmType string
	)
	switch vmType {
	default:
		vm, err = qemu.LoadVMByName(args[0])
	}
	if err != nil {
		return err
	}
	confirmationMessage, doIt, err := vm.Remove(args[0], machine.RemoveOptions{})
	if err != nil {
		return err
	}

	if !destoryOptions.Force {
		// Warn user
		fmt.Println(confirmationMessage)
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Are you sure you want to continue? [y/N] ")
		answer, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		if strings.ToLower(answer)[0] != 'y' {
			return nil
		}
	}
	return doIt()
}
