# sptfh

## Install

```sh
go install github.com/spotify-helper/sptfh/cmd/sptfh@latest
```

## Usage

```sh
$ sptfh --help
Spotify refuses to add basic things to their application, so here is a helper tool

Usage:
  spth [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  liked       Creates a playlist of your liked songs called "Liked"
  top-artists Get the top x artists
  top-tracks  Get the top x number of tracks for your user

Flags:
  -h, --help     help for spth
  -t, --toggle   Help message for toggle

Use "spth [command] --help" for more information about a command.
```

1. Follow [this](https://developer.spotify.com/documentation/general/guides/authorization/app-settings/) guide to get `SPOTIFY_ID` and `SPOTIFY_SECRET`
2. Set `SPOTIFY_ID`, and `SPOTIFY_SECRET`
```
export SPOTIFY_ID=asdasdasdasdasdasd
export SPOTIFY_SECRET=jsdliglksfdjglksgjfd
```
3. Add `http://localhost:8080` as the redirect URI on the spotify application you created.


## Commands

### `liked`
```
sptfh liked --help
Creates a playlist of your liked songs called "Liked"

Usage:
  spth liked [flags]

Flags:
  -h, --help   help for liked
```

### `top-artists`
```
$ sptfh  top-artists --help
Get the top x artists

Usage:
  spth top-artists [flags]

Flags:
  -c, --count int     The number of artists to return (default 10)
  -h, --help          help for top-artists
  -t, --term string   The term to get the top artists for. short_term/medium_term/long_term (default "short_term")
```
```
$ sptfh  top-artists -c 10 -t long_term
Please log in to Spotify by visiting the following page in your browser: https://accounts.spotify.com/authorize?....
1: Grateful Dead
2: Miles Davis
3: Neil Young
4: Weather Report
5: Crosby, Stills, Nash & Young
6: Steely Dan
7: Jefferson Airplane
8: Sly & The Family Stone
9: Talking Heads
10: Allman Brothers Band
```

### `top-tracks`
```
$ sptfh  top-tracks --help
Get the top x number of tracks for your user

Usage:
  spth top-tracks [flags]

Flags:
  -c, --count int     The number of artists to return (default 10)
  -h, --help          help for top-tracks
  -t, --term string   The term to get the top artists for. short_term/medium_term/long_term (default "short_term")
```
```
$ sptfh  top-tracks --term long_term --count 10
Please log in to Spotify by visiting the following page in your browser: https://accounts.spotify.com/authorize?....
1: Wharf Rat - Live, June 26/28, 1974
2: Tomorrow Never Knows - Remastered 2009
3: Powderfinger - 2016 Remaster
4: Teach Your Children
5: What a Day That Was - Live; Edit
6: 856
7: Ohio
8: Caution (Do Not Stop on Tracks) - Live at Wembley Empire Pool, April 1972
9: Turn Blue
10: Whipping Post - Live at the Atlanta International Pop Festival July 5, 1970
```
