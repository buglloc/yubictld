package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/buglloc/yubictld/internal/ykman"
)

var listCmd = &cobra.Command{
	Use:           "list",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "List Yubikeys",
	RunE: func(_ *cobra.Command, _ []string) error {
		ykm := ykman.NewYkMan()
		if err := ykm.ReloadDevices(); err != nil {
			return fmt.Errorf("could not reload yubikeys: %w", err)
		}

		for _, dev := range ykm.Devices() {
			fmt.Printf("- %s:\n", dev.Path())
			fmt.Printf("\tserial: %d\n", dev.Serial())
			fmt.Printf("\tlocation: %v\n", dev.Location())
			fmt.Printf("\tfree: %v\n", dev.IsFree())
		}

		return nil
	},
}
