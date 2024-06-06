#!/bin/bash

set -eEuo pipefail

artist=$(aquamarine artists ls |
	fzf --with-nth '2' -d '\t' \
		--border=top \
		--border-label='Play music' \
		--preview-window nohidden:border-sharp:right:50% \
		--height 100% \
		--min-height 30 \
		--bind 'enter:become(echo {1})' \
		--preview 'echo Album count: {3}')

if [[ $artist == "" ]]; then
	exit
fi

album=$(aquamarine albums ls "$artist" |
	fzf --with-nth '2' -d '\t' \
		--border=top \
		--border-label='Play music' \
		--preview-window nohidden:border-sharp:right:50% \
		--height 100% \
		--min-height 30 \
		--bind 'enter:become(echo {1})' \
		--header=$'\e[1;34m<ctrl-a>\e[0m queue album\n\e[1;34m<ctrl-s>\e[0m scrobble album\n\e[1;34m<enter>\e[0m list tracks\n\n' \
		--bind 'ctrl-a:become(aquamarine queue {1})' \
		--bind 'ctrl-s:become(aquamarine scrobble {1})' \
		--preview "echo 'Year: {3}\nGenre: {4}\nSong count: {5}\nDuration: {6}'")

if [[ $album == "" ]]; then
	exit
fi

aquamarine songs ls "$album" |
	fzf --with-nth '2' -d '\t' \
		--border=top \
		--border-label='Play music' \
		--preview-window nohidden:border-sharp:right:50% \
		--no-sort \
		--height 100% \
		--min-height 30 \
		--header=$'\e[1;34m<ctrl-a>\e[0m queue album\n\e[1;34m<ctrl-s>\e[0m scrobble track\n\e[1;34m<enter>\e[0m queue track\n\n' \
		--bind 'ctrl-a:become(aquamarine queue {9})' \
		--bind 'ctrl-s:become(aquamarine scrobble {1})' \
		--bind 'enter:become(aquamarine queue {1})' \
		--preview "echo 'Track: {3}\nDisc: {4}\nType: {5}\nBit-rate: {6}\nDuration: {7}\nSize: {8}'"
