package version

import (
	"fmt"
	"github.com/spf13/cobra"
)

// RootCmd 查看版本号命令
var RootCmd = &cobra.Command{
	Use:   "version",
	Short: "version number",
	Long:  `the version number of rednote sxt replyer`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("v1.0.0")
	},
}
