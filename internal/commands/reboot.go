package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/buglloc/yubictld/internal/ykman"
)

var rebootCmd = &cobra.Command{
	Use:           "reboot",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "Reboot Yubikeys",
	RunE: func(_ *cobra.Command, _ []string) error {
		ykm := ykman.NewYkMan()
		if err := ykm.ReloadDevices(); err != nil {
			return fmt.Errorf("could not reload yubikeys: %w", err)
		}

		for _, dev := range ykm.Devices() {
			if err := dev.Reboot(); err != nil {
				fmt.Printf("could not reboot %s: %v\n", dev.String(), err)
				continue
			}

			fmt.Printf("%s was rebooted", dev.String())
		}

		return nil
	},
}
