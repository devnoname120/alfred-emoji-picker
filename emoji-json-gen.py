#!/usr/bin/env python3
import os
import json
import requests

def download_json(url, filename):
    """
    Downloads a JSON file from the given URL and saves it with the specified filename.
    """
    print(f"Downloading {url}...")
    response = requests.get(url, timeout=10)
    response.raise_for_status()  # Raise an error for bad responses (4xx, 5xx)

    with open(filename, 'w', encoding='utf-8') as f:
        json.dump(response.json(), f, ensure_ascii=False, indent=2)

    print(f"Saved as {filename}")

def get_code_point(emoji):
    """
    Returns the Unicode code points of an emoji string as uppercase hex values,
    joined by spaces if it's a multi-character sequence.
    """
    return " ".join(f"{ord(char):X}" for char in emoji)

# Define URLs and file paths
emoji_keywords_url = "https://github.com/muan/emojilib/raw/refs/heads/main/dist/emoji-en-US.json"
emoji_data_url = "https://github.com/muan/unicode-emoji-json/raw/refs/heads/main/data-by-emoji.json"

emoji_keywords_file = "emoji-en-US.json"
emoji_data_file = "data-by-emoji.json"

# Download JSON files
download_json(emoji_keywords_url, emoji_keywords_file)
download_json(emoji_data_url, emoji_data_file)

# Load the downloaded JSON files
with open(emoji_keywords_file, 'r', encoding='utf-8') as f:
    emoji_keywords = json.load(f)

with open(emoji_data_file, 'r', encoding='utf-8') as f:
    emoji_data = json.load(f)

emoji_info = {}

# Combine data from both JSON sources
for emoji_symbol, keywords in emoji_keywords.items():
    if emoji_symbol in emoji_data:
        emoji_data[emoji_symbol]['codepoint'] = get_code_point(emoji_symbol)
        emoji_data[emoji_symbol]['keywords'] = keywords
        emoji_info[emoji_symbol] = emoji_data[emoji_symbol]

packed_emoji = {
    "emoji": emoji_info
}

with open('emoji.pack.new.json', 'w', encoding='utf-8') as f:
    json.dump(packed_emoji, f, ensure_ascii=False, indent=2)

print(f"Emoji data saved to emoji.pack.new.json")
