#!/usr/bin/env python3
import json
import sys

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

def escape_string(s):
    """Escape backslashes and double quotes for Go string literals."""
    return '"' + s.replace('\\', '\\\\').replace('"', '\\"') + '"'

def main():
    with open('emoji.pack.json') as emo:
        data = json.load(emo)
        emojis = data.get("emoji", {})

        emoji_entries = "\n".join(
            ENTRY_TEMPLATE.format(
                name=escape_string(info.get("slug", "")),
                category=escape_string(info.get("group", "")),
                char=escape_string(char),
                keywords=", ".join(escape_string(kw) for kw in info.get("keywords", []))
            )
            for char, info in emojis.items()
        )

        print(TEMPLATE.format(emoji_entries=emoji_entries))

if __name__ == "__main__":
    main()
