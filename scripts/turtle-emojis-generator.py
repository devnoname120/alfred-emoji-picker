#!/usr/bin/env python3
import requests

# How to use: uvx --from ./scripts ./scripts/turtle-emojis-generator.py > emojis.go

def download_json(url):
    response = requests.get(url, timeout=10)
    response.raise_for_status()
    return response.json()

def escape_string(s):
    return '"' + s.replace('\\', '\\\\').replace('"', '\\"') + '"'

TEMPLATE = """package turtle

// Emoji holds info about an emoji character
type Emoji struct {{
    Name     string   `json:"name" toml:"name"`
    Slug     string   `json:"slug" toml:"slug"`
    Category string   `json:"category" toml:"category"`
    Char     string   `json:"char" toml:"char"`
    Keywords []string `json:"keywords" toml:"keywords"`
}}

// String implementation for Emoji
func (e Emoji) String() string {{
    return e.Slug
}}

// emojis holds all available Emoji
var emojis = []*Emoji{{
{emoji_entries}
}}
"""

ENTRY_TEMPLATE = """    {{
        Name:     {name},
        Slug:     {slug},
        Category: {category},
        Char:     {char},
        Keywords: []string{{{keywords}}},
    }},"""

def main():
    emoji_keywords = download_json("https://github.com/muan/emojilib/raw/refs/heads/main/dist/emoji-en-US.json")
    emoji_data = download_json("https://github.com/muan/unicode-emoji-json/raw/refs/heads/main/data-by-emoji.json")

    for symbol, data in emoji_data.items():
        if symbol in emoji_keywords:
            data['keywords'] = emoji_keywords[symbol]

    emoji_entries = "\n".join(
        ENTRY_TEMPLATE.format(
            slug=escape_string(info.get("slug", "")),
            name=escape_string(info.get("name", "")),
            category=escape_string(info.get("group", "")),
            char=escape_string(char),
            keywords=", ".join(escape_string(kw) for kw in info.get("keywords", []))
        )
        for char, info in emoji_data.items()
    )
    go_code = TEMPLATE.format(emoji_entries=emoji_entries)
    print(go_code)

if __name__ == "__main__":
    main()
