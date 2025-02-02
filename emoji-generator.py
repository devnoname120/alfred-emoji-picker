#!/usr/bin/env python3
import os
import sys
import json
import requests
from io import BytesIO
from PIL import Image
from fontTools.ttLib import TTFont

OUTPUT_DIR = os.path.join(os.getcwd(), "emojis")

def download_json(url):
    response = requests.get(url, timeout=10)
    response.raise_for_status()
    return response.json()

def get_sbix_image_for_emoji(char, font, strike):
    cmap = font.getBestCmap()

    # Try full emoji sequence first
    codepoints = [ord(c) for c in char]
    glyph_name = cmap.get(tuple(codepoints))  # Try full sequence

    if not glyph_name:
        # If full sequence fails, try individual codepoints (single character fallback)
        glyph_name = cmap.get(codepoints[0])

    if not glyph_name:
        raise ValueError(f"Glyph not found for codepoint(s) {', '.join(f'U+{cp:X}' for cp in codepoints)}")

    if glyph_name not in strike.glyphs:
        raise ValueError(f"No sbix record for glyph '{glyph_name}'")

    glyph = strike.glyphs[glyph_name]

    if not glyph.imageData:
        raise ValueError(f"No image data for glyph '{glyph_name}'")

    return Image.open(BytesIO(glyph.imageData))

def main():
    emoji_data_url = "https://github.com/muan/unicode-emoji-json/raw/refs/heads/main/data-by-emoji.json"
    emoji_data = download_json(emoji_data_url)
    ordered_emoji = list(emoji_data.keys())

    font_path = "/System/Library/Fonts/Apple Color Emoji.ttc"
    try:
        font = TTFont(font_path, fontNumber=0)
    except Exception as e:
        sys.exit(f"Failed to load font from {font_path}: {e}")

    if "sbix" not in font:
        sys.exit("Font does not contain an sbix table for color emoji images.")

    sbix_table = font["sbix"]
    if 64 in sbix_table.strikes:
        strike = sbix_table.strikes[64]
    else:
        available = list(sbix_table.strikes.keys())
        closest = min(available, key=lambda s: abs(s - 64))
        strike = sbix_table.strikes[closest]
        print(f"Using sbix strike for size {closest} instead of 64.")

    os.makedirs(OUTPUT_DIR, exist_ok=True)

    for char in ordered_emoji:
        data = emoji_data.get(char)
        try:
            img = get_sbix_image_for_emoji(char, font, strike)
            output_path = os.path.join(OUTPUT_DIR, f"{data['slug']}.png")
            img.save(output_path)
            print(f"Saved icon: {output_path}")
        except Exception as e:
            emoji_name = data.get("name", char) if data else char
            print(f'Skipped "{emoji_name}": {e}', file=sys.stderr)

if __name__ == "__main__":
    main()
