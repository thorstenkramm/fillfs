// Package registry registers all extension generators.
package registry

import (
	"github.com/thorstenkramm/fillfs/internal/generator"
	"github.com/thorstenkramm/fillfs/pkg/ext/doc"
	"github.com/thorstenkramm/fillfs/pkg/ext/docx"
	"github.com/thorstenkramm/fillfs/pkg/ext/jpg"
	"github.com/thorstenkramm/fillfs/pkg/ext/mp3"
	"github.com/thorstenkramm/fillfs/pkg/ext/mp4"
	"github.com/thorstenkramm/fillfs/pkg/ext/odt"
	"github.com/thorstenkramm/fillfs/pkg/ext/ogg"
	"github.com/thorstenkramm/fillfs/pkg/ext/pdf"
	"github.com/thorstenkramm/fillfs/pkg/ext/ppt"
	"github.com/thorstenkramm/fillfs/pkg/ext/rtf"
	"github.com/thorstenkramm/fillfs/pkg/ext/webp"
	"github.com/thorstenkramm/fillfs/pkg/ext/xlsx"
)

// Generators returns all registered extension generators.
func Generators() []generator.Generator {
	return []generator.Generator{
		doc.New(),
		docx.New(),
		jpg.New(),
		mp3.New(),
		mp4.New(),
		odt.New(),
		ogg.New(),
		pdf.New(),
		ppt.New(),
		rtf.New(),
		webp.New(),
		xlsx.New(),
	}
}
