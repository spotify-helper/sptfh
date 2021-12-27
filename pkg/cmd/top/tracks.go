package top

import (
	"fmt"
	"log"

	sptfh "github.com/spotify-helper/sptfh/pkg/spotify"
	"github.com/zmb3/spotify/v2"

	"github.com/spf13/cobra"
)

func NewCmdTopTracks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "top-tracks",
		Short: "Get the top x number of tracks for your user",
		Run: func(cmd *cobra.Command, args []string) {
			count, err := cmd.Flags().GetInt("count")
			if err != nil {
				log.Fatal(err)
			}

			term, err := cmd.Flags().GetString("term")
			if err != nil {
				log.Fatal(err)
			}

			if term != "short_term" && term != "medium_term" && term != "long_term" {
				log.Fatal("term must be short_term, medium_term, or long_term")
			}

			s, err := sptfh.Login()
			if err != nil {
				log.Fatal(err)
			}

			songs, err := s.GetTopXTracks(count, spotify.Range(term))
			if err != nil {
				log.Fatal(err)
			}
			for i, song := range *songs {
				fmt.Printf("%d: %s\n", i+1, song)
			}

		},
	}

	cmd.Flags().IntVarP(&count, "count", "c", 10, "The number of artists to return")
	cmd.Flags().StringVarP(&term, "term", "t", "short_term", "The term to get the top artists for. short_term/medium_term/long_term")

	return cmd
}
