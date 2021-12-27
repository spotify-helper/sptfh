package liked_playist

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spotify-helper/sptfh/pkg/spotify"
)

func NewCmdLikedPlaylist() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "liked",
		Short: "Creates a playlist of your liked songs called \"Liked\"",
		Run: func(cmd *cobra.Command, args []string) {

			if os.Getenv("SPOTIFY_ID") == "" {
				log.Fatal("SPOTIFY_ID not set")
			}

			if os.Getenv("SPOTIFY_SECRET") == "" {
				log.Fatal("SPOTIFY_SECRET not set")
			}

			s, err := spotify.Login()
			if err != nil {
				log.Fatal(err)
			}

			i, err := s.GetUserLikedSongsId()

			if err != nil {
				log.Fatal(err)
			}

			playlist, err := s.CreatePlaylistLikedSongs()
			if err != nil {
				log.Fatal(err)
			}

			err = s.AddSongsToLikedPlaylist(i, playlist)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("Liked Playlist created, and songs added")
		},
	}

	return cmd
}
