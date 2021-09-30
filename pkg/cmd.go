package pkg

import (
	"github.com/kyokomi/emoji/v2"
	"github.com/spf13/cobra"
)

type CMDFunc func(config *WildFireConfig, cmd *cobra.Command, args []string) (*WildFireConfig, bool)
type CobraCMDFunc func(cmd *cobra.Command, args []string)

func ProjectFunc(c CMDFunc) CobraCMDFunc {
	return func(cmd *cobra.Command, args []string) {
		emoji.Println()
		config, update := c(GetConfig(), cmd, args)
		emoji.Println()

		if update {
			cobra.CheckErr(config.SaveConfig())
			emoji.Println(":cloud: Configuration has been updated.")

			return
		}

		emoji.Println(":cloud: Dousing WildFire.")
	}
}
