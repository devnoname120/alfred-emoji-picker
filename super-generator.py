#!/usr/bin/env python3
import json
import requests

def download_json(url):
    response = requests.get(url, timeout=10)
    response.raise_for_status()
    return response.json()

def get_code_point(emoji):
    return " ".join(f"{ord(char):X}" for char in emoji)

def escape_string(s):
    return '"' + s.replace('\\', '\\\\').replace('"', '\\"') + '"'

TEMPLATE = """package turtle

// Emoji holds info about an emoji character
type Emoji struct {{
    Name     string   `json:"name" toml:"name"`
    Category string   `json:"category" toml:"category"`
    Char     string   `json:"char" toml:"char"`
    Keywords []string `json:"keywords" toml:"keywords"`
}}

// String implementation for Emoji
func (e Emoji) String() string {{
    return e.Char
}}

// emojis holds all available Emoji
var emojis = []*Emoji{{
{emoji_entries}
}}
"""

ENTRY_TEMPLATE = """    {{
        Name:     {name},
        Category: {category},
        Char:     {char},
        Keywords: []string{{{keywords}}},
    }},"""

def main():
    emoji_keywords_url = "https://github.com/muan/emojilib/raw/refs/heads/main/dist/emoji-en-US.json"
    emoji_data_url = "https://github.com/muan/unicode-emoji-json/raw/refs/heads/main/data-by-emoji.json"

    emoji_keywords = download_json(emoji_keywords_url)
    emoji_data = download_json(emoji_data_url)

    emoji_info = {}
    for emoji_symbol, keywords in emoji_keywords.items():
        if emoji_symbol in emoji_data:
            emoji_data[emoji_symbol]['codepoint'] = get_code_point(emoji_symbol)
            emoji_data[emoji_symbol]['keywords'] = keywords
            emoji_info[emoji_symbol] = emoji_data[emoji_symbol]

    emoji_entries = "\n".join(
        ENTRY_TEMPLATE.format(
            name=escape_string(info.get("slug", "")),
            category=escape_string(info.get("group", "")),
            char=escape_string(char),
            keywords=", ".join(escape_string(kw) for kw in info.get("keywords", []))
        )
        for char, info in emoji_info.items()
    )
    go_code = TEMPLATE.format(emoji_entries=emoji_entries)
    print(go_code)

if __name__ == "__main__":
    main()
