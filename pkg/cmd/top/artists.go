package top

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	sptfh "github.com/spotify-helper/sptfh/pkg/spotify"
	"github.com/zmb3/spotify/v2"
)

var (
	count int
	term  string
)

func NewCmdTopArtists() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "top-artists",
		Short: "Get the top x artists",
		Run: func(cmd *cobra.Command, args []string) {

			if os.Getenv("SPOTIFY_ID") == "" {
				log.Fatal("SPOTIFY_ID not set")
			}

			if os.Getenv("SPOTIFY_SECRET") == "" {
				log.Fatal("SPOTIFY_SECRET not set")
			}

			count, err := cmd.Flags().GetInt("count")
			if err != nil {
				log.Fatal(err)
			}

			term, err := cmd.Flags().GetString("term")
			if err != nil {
				log.Fatal(err)
			}

			if term != "short_term" && term != "medium_term" && term != "long_term" {
				log.Fatal("term must be short_term, medium_term or long_term")
			}

			s, err := sptfh.Login()
			if err != nil {
				log.Fatal(err)
			}

			a, err := s.GetTopXArtists(count, spotify.Range(term))
			if err != nil {
				log.Fatal(err)
			}

			for i, artist := range *a {
				fmt.Printf("%d: %s\n", i+1, artist)
			}
		},
	}

	cmd.PersistentFlags().IntVarP(&count, "count", "c", 10, "The number of artists to return")
	cmd.Flags().StringVarP(&term, "term", "t", "short_term", "The term to get the top artists for. short_term/medium_term/long_term")

	return cmd
}
