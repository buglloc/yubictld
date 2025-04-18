package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/buglloc/yubictld/internal/ykman"
)

var touchArgs struct {
	pin uint8
}

var touchCmd = &cobra.Command{
	Use:           "touch",
	SilenceUsage:  true,
	SilenceErrors: true,
	Short:         "Touch Yubikey",
	RunE: func(_ *cobra.Command, _ []string) error {
		if touchArgs.pin == 0 {
			return fmt.Errorf("must specify a pin to touch")
		}

		ykman.NewYkMan().ReloadDevices()
		//tDev, err := toucher.FirstDevice()
		//if err != nil {
		//	return fmt.Errorf("could not find a toucher device")
		//}
		//
		//return toucher.NewHwToucher(tDev).Touch(touchArgs.pin)
		return nil
	},
}

func init() {
	flags := touchCmd.PersistentFlags()
	flags.Uint8Var(&touchArgs.pin, "pin", 0, "Pin number to touch")
}
