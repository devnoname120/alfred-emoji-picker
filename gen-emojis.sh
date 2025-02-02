#!/usr/bin/env bash

set -exuo pipefail

# Optional: -28% in size (76.7 MB → 54.9 MB) lossy but cannot see any difference visually
pngquant --verbose --force --skip-if-larger --speed=1 --quality=95-100 --ext=.png emojis/*
.png || true

# Optional: -13% in size (54.9 MB → 47.8 MB) lossless so why not I guess
oxipng --dir emojis --strip safe --interlace 0 --opt=max --alpha emojis/*.png || true

# ⇒ Total: -38% in size (76.7 MB → 47.8 MB) lossy but cannot see any difference visually
# It seems to be a tiny bit faster to load in Alfred as well, possibly thanks to the PNG palette complexity reduction which makes them parse faster?
