// Command profile gets the public profile information about a Spotify user.
package spotify

import (
	"context"
	"fmt"
	"log"
	"net/http"

	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/zmb3/spotify/v2"
)

const redirectURI = "http://localhost:8080/callback"

var (
	auth = spotifyauth.New(
		spotifyauth.WithRedirectURL(redirectURI),
		spotifyauth.WithScopes(
			spotifyauth.ScopeUserLibraryRead,
			spotifyauth.ScopePlaylistModifyPrivate,
			spotifyauth.ScopePlaylistReadPrivate,
			spotifyauth.ScopeUserTopRead,
		),
	)
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

type Spotify struct {
	client        *spotify.Client
	currentUserID string
}

type ISpotify interface {
	GetUserLikedSongs() (*[]string, error)
	GetUserLikedSongsId() (*[]spotify.ID, error)
	GetCurrentUserId() string
	GetTopXArtists(x int, term spotify.Range) (*[]string, error)
	GetTopXTracks(x int, term spotify.Range) (*[]string, error)

	CreatePlaylistLikedSongs() (playlistID spotify.ID, err error)

	AddSongsToLikedPlaylist(songIDS *[]spotify.ID, playlistID spotify.ID) error

	RemoveDuplicateSongsFromList(songIDS *[]spotify.ID, playlistID spotify.ID) (*[]spotify.ID, error)
}

func Login() (s ISpotify, string error) {

	// first start an HTTP server
	http.HandleFunc("/callback", completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client := <-ch

	user, err := client.CurrentUser(context.Background())
	if err != nil {
		return nil, err
	}

	if user.ID == "" {
		return nil, fmt.Errorf("user id is empty")
	}

	return &Spotify{client: client, currentUserID: user.ID}, nil
}

// gets the current users liked songs
func (s *Spotify) GetUserLikedSongs() (*[]string, error) {

	ls, err := s.client.CurrentUsersTracks(context.Background())
	if err != nil {
		return nil, err
	}

	songs := []string{}

	for page := 1; ; page++ {
		for _, item := range ls.Tracks {
			songs = append(songs, item.Name)
		}
		err = s.client.NextPage(context.Background(), ls)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	return &songs, nil
}

// creates a playlist for the current user
// if it doesnt exist, otherwise it returns the id of the existing playlist
func (s *Spotify) CreatePlaylistLikedSongs() (playlistID spotify.ID, err error) {

	existingPlaylists, err := s.client.GetPlaylistsForUser(context.Background(), s.GetCurrentUserId())
	if err != nil {
		return "", err
	}

	for page := 1; ; page++ {
		for _, item := range existingPlaylists.Playlists {
			if item.Name == "Liked" {
				return item.ID, nil
			}
		}

		err = s.client.NextPage(context.Background(), existingPlaylists)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			return "", err
		}
	}

	playlist, err := s.client.CreatePlaylistForUser(context.Background(), s.GetCurrentUserId(), "Liked", "", false, false)
	if err != nil {
		return "", err
	}
	return playlist.ID, nil
}

// adds the liked songs to the playlist
func (s *Spotify) AddSongsToLikedPlaylist(songIDS *[]spotify.ID, playlistID spotify.ID) (err error) {

	// add tracts to playlist every 100 elements

	songIDS, err = s.RemoveDuplicateSongsFromList(songIDS, playlistID)
	if err != nil {
		return err
	}

	for i := 0; i < len(*songIDS); i += 100 {
		end := i + 100
		if i+100 > len(*songIDS) {
			end = len(*songIDS)
		}
		t := (*songIDS)[i:end]
		_, err := s.client.AddTracksToPlaylist(context.Background(), playlistID, t...)
		if err != nil {
			return err
		}
	}

	return nil
}

// gets the current users liked songs ids
func (s *Spotify) GetUserLikedSongsId() (*[]spotify.ID, error) {

	ls, err := s.client.CurrentUsersTracks(context.Background())
	if err != nil {
		return nil, err
	}

	songIDS := []spotify.ID{}

	for page := 1; ; page++ {
		for _, item := range ls.Tracks {
			songIDS = append(songIDS, item.ID)
		}
		err = s.client.NextPage(context.Background(), ls)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	return &songIDS, nil
}

func (s *Spotify) GetTopXArtists(x int, term spotify.Range) (*[]string, error) {

	if term != spotify.LongTermRange && term != spotify.ShortTermRange && term != spotify.MediumTermRange {
		return nil, fmt.Errorf("term must be one of LongTermRange, ShortTermRange, MediumTermRange")
	}

	topArtists, err := s.client.CurrentUsersTopArtists(
		context.Background(),
		spotify.Timerange(term),
	)
	if err != nil {
		return nil, err
	}

	artists := []string{}
	for page := 1; ; page++ {
		for _, item := range topArtists.Artists {
			artists = append(artists, item.Name)
		}
		err = s.client.NextPage(context.Background(), topArtists)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	// return only x number of artists
	if len(artists) < x {
		return &artists, nil
	}

	artists = artists[:x]
	return &artists, nil
}

// gets the current users id
func (s *Spotify) GetCurrentUserId() string {

	return s.currentUserID
}

// removes duplicate songs from the existing list
func (s *Spotify) RemoveDuplicateSongsFromList(songIDS *[]spotify.ID, playlistID spotify.ID) (*[]spotify.ID, error) {

	existingSongs, err := s.client.GetPlaylistTracks(context.Background(), playlistID)
	if err != nil {
		return nil, err
	}

	for page := 1; ; page++ {
		for _, item := range existingSongs.Tracks {
			for i, song := range *songIDS {
				if item.Track.ID == song {
					*songIDS = append((*songIDS)[:i], (*songIDS)[i+1:]...)
				}
			}
		}
		err = s.client.NextPage(context.Background(), existingSongs)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	return songIDS, nil
}

func (s *Spotify) GetTopXTracks(x int, term spotify.Range) (*[]string, error) {

	if term != spotify.LongTermRange && term != spotify.ShortTermRange && term != spotify.MediumTermRange {
		return nil, fmt.Errorf("term must be one of LongTermRange, ShortTermRange, MediumTermRange")
	}

	topSongs, err := s.client.CurrentUsersTopTracks(
		context.Background(),
		spotify.Timerange(term),
	)
	if err != nil {
		return nil, err
	}

	songs := []string{}
	for page := 1; ; page++ {
		for _, item := range topSongs.Tracks {
			songs = append(songs, item.Name)
		}
		err = s.client.NextPage(context.Background(), topSongs)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			return nil, err
		}
	}

	// return only x number of songs
	if len(songs) < x {
		return &songs, nil
	}

	songs = songs[:x]
	return &songs, nil
}

func completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := spotify.New(auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}
