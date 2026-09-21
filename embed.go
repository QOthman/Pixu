package pixu

import (
	_ "embed"
)

// Embedded default bitmap font atlases so Pixu works out-of-the-box
// with zero external asset dependencies.

//go:embed font/font_atlas.png
var defaultFontAtlasData []byte

//go:embed font/font_atlas_bold.png
var defaultFontAtlasBoldData []byte
