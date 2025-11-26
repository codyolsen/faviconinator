# faviconinator

Generate a full favicon set (PNG, ICO, SVG) from a single SVG source.

## Requirements

- Go 1.22+
- ImageMagick 7+ with the `magick` binary on your `PATH` (install guide: [imagemagick.org/script/download.php](https://imagemagick.org/script/download.php))

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
- `-color`: optional hex color (e.g. `#ff6600`) to tint the icon
- `-v` / `-verbose`: verbose logging
- `-version`: print version and exit
- `-jobs`: number of concurrent workers (default: CPU count)
- `-json`: print summary (files, outputs, workers, duration) as JSON

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
- `favicon.svg` (with white fills stripped to keep transparency)
- Android/Apple/Microsoft sized PNGs (32px–512px)

Generated PNG sizes:

- Android/Chrome: 36, 48, 72, 96, 144, 192, 192, 512
- Apple touch: 57, 60, 72, 76, 114, 120, 144, 152, 180 (also `apple-icon.png`, `apple-icon-precomposed.png`, `apple-touch-icon.png`)
- Favicons: 32, 96, 96
- Microsoft tiles: 70, 144, 150, 310

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
