# faviconinator ![CI](https://github.com/codyolsen/faviconinator/actions/workflows/ci.yml/badge.svg)

## Why do I need this?

Generate a full favicon set (PNG, ICO, SVG) from a single SVG source. 

Step right up, folks! One gleaming SVG rolls in, a full parade of browser-ready icons rolls out—every size, every platform, even the ICO—fresh off the line in one smooth command, CI-stamped and ready to bolt onto your site’s chrome.
GAD ZOOKS!

## Speedy Boi

Benchmarks (darwin/arm64, ImageMagick on PATH):
- workers=1: ~4.25s/op
- workers=16: ~0.47s/op
- Speedup: ~9.0x (≈88.9% time reduction)

## Requirements

- Go 1.22+
- ImageMagick 7+ with the `magick` binary on your `PATH` (install guide: [imagemagick.org/script/download.php](https://imagemagick.org/script/download.php))
- Optional (better SVG rendering): `rsvg-convert` from librsvg. macOS: `brew install librsvg`.

## Install

```sh
go install github.com/codyolsen/faviconinator/cmd/faviconinator@latest
```

More details and release notes live on GitHub: [github.com/codyolsen/faviconinator](https://github.com/codyolsen/faviconinator).

## Usage

```sh
faviconinator [flags] input.svg
```

Flags:

- `-out` (default `build/<input basename>`): directory to write generated files
- `-color`: (reserved) hex color to tint the icon
- `-v` / `-verbose`: verbose logging
- `-version`: print version and exit
- `-jobs`: number of concurrent workers (default: CPU count)
- `-json`: print summary (files, outputs, workers, duration) as JSON
- `-renderer`: PNG renderer (`auto` default prefers `rsvg-convert` if present, else `magick`)

Example:

```sh
faviconinator -out dist -color "#3366ff" bl_square.svg
```

Platform reference links:

- Favicons: [MDN link relation `icon`](https://developer.mozilla.org/en-US/docs/Web/HTML/Link_types/icon)
- Android / Chrome: [Web App Manifest icons](https://developer.mozilla.org/en-US/docs/Web/Manifest/icons)
- Apple touch icons: [`apple-touch-icon` link relation](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/link#attr-apple-touch-icon)
- Microsoft tiles: [`msapplication-*` meta tags](https://developer.mozilla.org/en-US/docs/Web/HTML/Element/meta/name#msapplication-tileimage)

Stats: command prints total assets, destination, duration, and worker count after completion.

Outputs:

- `favicon.ico`
- `favicon.svg` (copied verbatim from source)
- Android/Apple/Microsoft sized PNGs (32px–512px)

Generated PNG sizes:

- Android/Chrome: 36, 48, 72, 96, 144, 192, 192, 512
- Apple touch: 57, 60, 72, 76, 114, 120, 144, 152, 180 (also `apple-icon.png`, `apple-icon-precomposed.png`, `apple-touch-icon.png`)
- Favicons: 32, 96, 96
- Microsoft tiles: 70, 144, 150, 310

All generated PNGs (smallest to largest):
- 32: `favicon-32x32.png`
- 36: `android-icon-36x36.png`
- 48: `android-icon-48x48.png`
- 57: `apple-icon-57x57.png`
- 60: `apple-icon-60x60.png`
- 70: `ms-icon-70x70.png`
- 72: `android-icon-72x72.png`, `apple-icon-72x72.png`
- 76: `apple-icon-76x76.png`
- 96: `android-icon-96x96.png`, `favicon-96x96.png`, `favicon.png`
- 114: `apple-icon-114x114.png`
- 120: `apple-icon-120x120.png`
- 144: `android-icon-144x144.png`, `apple-icon-144x144.png`, `ms-icon-144x144.png`
- 150: `ms-icon-150x150.png`
- 152: `apple-icon-152x152.png`
- 180: `apple-icon-180x180.png`, `apple-icon.png`, `apple-icon-precomposed.png`, `apple-touch-icon.png`
- 192: `android-icon-192x192.png`, `android-chrome-192x192.png`
- 310: `ms-icon-310x310.png`
- 512: `android-chrome-512x512.png`

Add to your site (HTML example):

```html
<!-- General -->
<link rel="icon" type="image/png" sizes="32x32" href="/favicon-32x32.png">
<link rel="icon" type="image/png" sizes="96x96" href="/favicon-96x96.png">
<link rel="icon" type="image/x-icon" href="/favicon.ico">
<link rel="manifest" href="/site.webmanifest">

<!-- Apple -->
<link rel="apple-touch-icon" sizes="180x180" href="/apple-icon-180x180.png">
<link rel="apple-touch-icon" href="/apple-touch-icon.png">

<!-- Android / Chrome -->
<link rel="icon" type="image/png" sizes="192x192" href="/android-chrome-192x192.png">
<link rel="icon" type="image/png" sizes="512x512" href="/android-chrome-512x512.png">

<!-- Microsoft tiles -->
<meta name="msapplication-TileColor" content="#ffffff">
<meta name="msapplication-TileImage" content="/ms-icon-144x144.png">
```

## Development

```sh
make fmt
make test
make build
```

`make build` writes the binary to `bin/faviconinator`.

<p>
  <img alt="icon 32px" src="internal/generate/testdata/complexahexagon.svg" width="32" height="32" />
  <img alt="icon 36px" src="internal/generate/testdata/complexahexagon.svg" width="36" height="36" />
  <img alt="icon 48px" src="internal/generate/testdata/complexahexagon.svg" width="48" height="48" />
  <img alt="icon 57px" src="internal/generate/testdata/complexahexagon.svg" width="57" height="57" />
  <img alt="icon 60px" src="internal/generate/testdata/complexahexagon.svg" width="60" height="60" />
  <img alt="icon 70px" src="internal/generate/testdata/complexahexagon.svg" width="70" height="70" />
  <img alt="icon 72px" src="internal/generate/testdata/complexahexagon.svg" width="72" height="72" />
  <img alt="icon 76px" src="internal/generate/testdata/complexahexagon.svg" width="76" height="76" />
  <img alt="icon 96px" src="internal/generate/testdata/complexahexagon.svg" width="96" height="96" />
  <img alt="icon 114px" src="internal/generate/testdata/complexahexagon.svg" width="114" height="114" />
  <img alt="icon 120px" src="internal/generate/testdata/complexahexagon.svg" width="120" height="120" />
  <img alt="icon 144px" src="internal/generate/testdata/complexahexagon.svg" width="144" height="144" />
  <img alt="icon 150px" src="internal/generate/testdata/complexahexagon.svg" width="150" height="150" />
  <img alt="icon 152px" src="internal/generate/testdata/complexahexagon.svg" width="152" height="152" />
  <img alt="icon 180px" src="internal/generate/testdata/complexahexagon.svg" width="180" height="180" />
  <img alt="icon 192px" src="internal/generate/testdata/complexahexagon.svg" width="192" height="192" />
  <img alt="icon 310px" src="internal/generate/testdata/complexahexagon.svg" width="310" height="310" />
  <img alt="icon 512px" src="internal/generate/testdata/complexahexagon.svg" width="512" height="512" />
</p>
