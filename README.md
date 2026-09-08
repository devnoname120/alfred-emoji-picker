# Alfred emoji picker

Input emojis from Alfred — at a blazing-fast speed!

<p align="center">
    <img src="https://user-images.githubusercontent.com/2824100/174484132-c76cf892-27e8-4d8c-bec7-76745016fe1a.png" data-canonical-src="https://user-images.githubusercontent.com/2824100/174484132-c76cf892-27e8-4d8c-bec7-76745016fe1a.png" width="400"/>
</p>

# Install

- Download the workflow from the [latest release](https://github.com/devnoname120/alfred-emoji-picker/releases/latest).
- Open the file and import it into Alfred.
- **Click on the workflow in Alfred and define a hotkey**.

👉 I recommend using <kbd>Command ⌘</kbd> <kbd>Control ⌃</kbd> <kbd>Space</kbd>

# Build

```shell
go install

./build.sh
```

Copy the executable in the Alfred workflow directory and export the new workflow from Alfred.

# Frequently used emoji

The binary now supports a small local usage database in Alfred's workflow data directory:

- Empty query shows only emoji the user has actually selected before.
- The amount shown on empty query is controlled by the Alfred variable `frequent_emoji_limit`.
- When the query is empty, previously selected emoji are ordered by usage count.
- Type `reset` in the picker and select **Reset frequent emoji** to clear the history.

To reset history from a terminal, run the binary from the workflow directory and
set `alfred_workflow_data` to the workflow's data directory:

```shell
alfred_workflow_data="$HOME/Library/Application Support/Alfred/Workflow Data/com.github.devnoname120.alfred-emoji-picker" \
  ./alfred-emoji-picker --reset-frequent
```

Both `--reset-frequent` and `--record EMOJI` use `alfred_workflow_data` and run
without Alfred's other environment variables. Alfred sets this directory
automatically when running the workflow. For a custom data location, set the
variable to that directory. Missing configuration or storage errors are reported
on stderr with a nonzero exit status.

# Update emojis

1) Emoji metadata (names, slugs, keywords) lives in the [`turtle`](https://github.com/devnoname120/turtle) module. Check the README on how to regenerate `emojis.go` to make it up-to-date, and push a new tag.

2) Bump the dependency of `alfred-emoji-picker` via `go get github.com/devnoname120/turtle@latest`.

3) Re-render the PNGs from the current Apple Color Emoji font and re-optimize them (macOS only; requires [uv](https://docs.astral.sh/uv/), `pngquant`, and `oxipng`):

    ```shell
    uv run --project scripts scripts/emojis-generator.py
    ./gen-emojis.sh
    ```

# TODO

- [x] Restore clipboard after pasting emoji
- [x] Support frequently used emoji
- [ ] Support for multiple words fuzzy search
- [x] Add scoring on results (exact match > partial match at beginning > partial match > keywords, categories, etc…)
- [x] Add scripts to update the emoji database
- [x] Support for auto-updates
- [ ] Support for skin tones (note: `Turtle` doesn't support them)
