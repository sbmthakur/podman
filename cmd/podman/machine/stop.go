package machine

import (
	"github.com/containers/common/pkg/completion"
	"github.com/containers/podman/v3/cmd/podman/registry"
	"github.com/containers/podman/v3/pkg/domain/entities"
	"github.com/containers/podman/v3/pkg/machine"
	"github.com/containers/podman/v3/pkg/machine/qemu"
	"github.com/spf13/cobra"
)

var (
	stopCmd = &cobra.Command{
		Use:               "stop NAME",
		Short:             "Stop an existing machine",
		Long:              "Stop an existing machine ",
		RunE:              stop,
		Args:              cobra.ExactArgs(1),
		Example:           `podman machine stop myvm`,
		ValidArgsFunction: completion.AutocompleteNone,
	}
)

func init() {
	registry.Commands = append(registry.Commands, registry.CliCommand{
		Mode:    []entities.EngineMode{entities.ABIMode, entities.TunnelMode},
		Command: stopCmd,
		Parent:  machineCmd,
	})
}

// TODO  Name shouldnt be required, need to create a default vm
func stop(cmd *cobra.Command, args []string) error {
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
	return vm.Stop(args[0], machine.StopOptions{})
}
