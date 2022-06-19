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

# TODO

- [ ] Support for multiple words fuzzy search
- [ ] Add scoring on results (exact match > partial match at beginning > partial match > keywords, categories, etc…)
- [ ] Add a script to scrap https://emojis.wiki/apple/ to regenerate the emoji images
- [ ] Support for auto-updates
- [ ] Support for skin tones
