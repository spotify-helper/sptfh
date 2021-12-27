# sptfh

## Install

```sh
go install github.com/spotify-helper/sptfh/cmd/sptfh@latest
```

## Usage

```sh
./sptfh --help
Spotify refuses to add basic things to their application, so here is a helper tool

Usage:
  spth [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  liked       Creates a playlist of your liked songs called "Liked"

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
