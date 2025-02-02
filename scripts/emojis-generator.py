#!/usr/bin/env python3
import shutil
from pathlib import Path
import sys
import requests
from AppKit import (
    NSImage,
    NSFont,
    NSAttributedString,
    NSColor,
    NSMakeRect,
    NSRectFill,
    NSBitmapImageRep,
    NSPNGFileType,
)

# Define output directory relative to the script's location
OUTPUT_DIR = Path.cwd() / "emojis"

try:
    shutil.rmtree(OUTPUT_DIR)
except FileNotFoundError: # Python 3.11+: Only thrown if the passed directory doesn't exist
    pass

OUTPUT_DIR.mkdir(parents=True, exist_ok=True)

# Note: This actually renders as 128x128 due to Retina automatic 2x resolution rendering
IMAGE_SIZE = 64

def download_json(url):
    response = requests.get(url, timeout=10)
    response.raise_for_status()
    return response.json()

def render_emoji_image(emoji_str, size=IMAGE_SIZE):
    """Render an emoji as a PNG image using Apple's Color Emoji font."""
    font = NSFont.fontWithName_size_("Apple Color Emoji", size)
    if font is None:
        raise ValueError("Could not load Apple Color Emoji font.")

    attributes = {"NSFont": font, "NSForegroundColor": NSColor.blackColor()}
    attributed_str = NSAttributedString.alloc().initWithString_attributes_(emoji_str, attributes)
    text_size = attributed_str.size()

    # Note: the canvas is actually (size*2, size*2) because of retina rendering
    image = NSImage.alloc().initWithSize_((size, size))
    image.lockFocus()
    NSColor.clearColor().set()
    NSRectFill(NSMakeRect(0, 0, size, size))

    x = (size - text_size.width) / 2.0
    y = (size - text_size.height) / 2.0 - 2
    attributed_str.drawAtPoint_((x, y))

    image.unlockFocus()
    tiff_data = image.TIFFRepresentation()
    bitmap_rep = NSBitmapImageRep.imageRepWithData_(tiff_data)
    png_data = bitmap_rep.representationUsingType_properties_(NSPNGFileType, None)

    return bytes(png_data)

def main():
    emoji_data_url = "https://github.com/muan/unicode-emoji-json/raw/refs/heads/main/data-by-emoji.json"
    emoji_data = download_json(emoji_data_url)

    for emoji_char, data in emoji_data.items():
        slug = data.get("slug", "unknown")
        try:
            png_data = render_emoji_image(emoji_char, IMAGE_SIZE)
            output_path = OUTPUT_DIR / f"{slug}.png"
            output_path.write_bytes(png_data)
            print(f"Saved icon: {output_path}")
        except Exception as e:
            emoji_name = data.get("name", emoji_char)
            print(f'Skipped "{emoji_name}": {e}', file=sys.stderr)

if __name__ == "__main__":
    main()
