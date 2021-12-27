package root

import (
	"os"

	"github.com/spf13/cobra"
	liked_playist "github.com/spotify-helper/sptfh/pkg/cmd/liked-playlist"
	"github.com/spotify-helper/sptfh/pkg/cmd/top"
)

// rootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "spth",
	Short: "A spotify helper tool",
	Long:  `Spotify refuses to add basic things to their application, so here is a helper tool`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.spth.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// Add subcommands
	RootCmd.AddCommand(liked_playist.NewCmdLikedPlaylist())
	RootCmd.AddCommand(top.NewCmdTopArtists())
	RootCmd.AddCommand(top.NewCmdTopTracks())
}

func Run() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
