package pkg

import (
	"fmt"
	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
)

type CMDFunc func(config *WildFireConfig, cmd *cobra.Command, args []string) (*WildFireConfig, bool, error)
type CobraCMDFunc func(cmd *cobra.Command, args []string) error

func ProjectFunc(c CMDFunc) CobraCMDFunc {
	return func(cmd *cobra.Command, args []string) error {
		fmt.Println()
		config, update, err := c(GetConfig(), cmd, args)
		fmt.Println()

		if err != nil {
			return err
		}

		if update {
			cobra.CheckErr(config.SaveConfig())
			emoji.Println(":cloud: Configuration has been updated.")

			return nil
		}

		emoji.Println(":cloud: Dousing WildFire.")
		return nil
	}
}
