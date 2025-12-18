// Package sources declares static seed metadata used to populate files.
package sources

// Seed describes a source file used to populate generated files.
type Seed struct {
	URL       string
	FileName  string
	Extension string
	Size      int64
}

// All seeds are static to allow accurate planning before downloads.
var All = []Seed{
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/img_01.jpg",
		FileName:  "img_01.jpg",
		Extension: ".jpg",
		Size:      2624144,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/img_02.jpg",
		FileName:  "img_02.jpg",
		Extension: ".jpg",
		Size:      1304804,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/img_03.jpg",
		FileName:  "img_03.jpg",
		Extension: ".jpg",
		Size:      881435,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/img_04.jpg",
		FileName:  "img_04.jpg",
		Extension: ".jpg",
		Size:      2052754,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/img_05.jpg",
		FileName:  "img_05.jpg",
		Extension: ".jpg",
		Size:      581189,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/img_06.jpg",
		FileName:  "img_06.jpg",
		Extension: ".jpg",
		Size:      1460410,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/img_07.jpg",
		FileName:  "img_07.jpg",
		Extension: ".jpg",
		Size:      843609,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/img_500kB.webp",
		FileName:  "img_500kB.webp",
		Extension: ".webp",
		Size:      517842,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/img_50kB.webp",
		FileName:  "img_50kB.webp",
		Extension: ".webp",
		Size:      50408,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/opendoc_100kB.odt",
		FileName:  "opendoc_100kB.odt",
		Extension: ".odt",
		Size:      116076,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/portable_doc_150kB.pdf",
		FileName:  "portable_doc_150kB.pdf",
		Extension: ".pdf",
		Size:      142786,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/portable_doc_500_kB.pdf",
		FileName:  "portable_doc_500_kB.pdf",
		Extension: ".pdf",
		Size:      469513,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/powerpoint.ppt",
		FileName:  "powerpoint.ppt",
		Extension: ".ppt",
		Size:      1028608,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/richtext_300kB.rtf",
		FileName:  "richtext_300kB.rtf",
		Extension: ".rtf",
		Size:      295392,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/sound.mp3",
		FileName:  "sound.mp3",
		Extension: ".mp3",
		Size:      1059386,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/sound.ogg",
		FileName:  "sound.ogg",
		Extension: ".ogg",
		Size:      1032948,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/spreadsheet_01.xlsx",
		FileName:  "spreadsheet_01.xlsx",
		Extension: ".xlsx",
		Size:      5425,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/spreadsheet_02.xlsx",
		FileName:  "spreadsheet_02.xlsx",
		Extension: ".xlsx",
		Size:      9299,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/spreadsheet_03.xlsx",
		FileName:  "spreadsheet_03.xlsx",
		Extension: ".xlsx",
		Size:      188887,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/video.mp4",
		FileName:  "video.mp4",
		Extension: ".mp4",
		Size:      3114374,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/word_100kB.docx",
		FileName:  "word_100kB.docx",
		Extension: ".docx",
		Size:      111303,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/word_1MB.docx",
		FileName:  "word_1MB.docx",
		Extension: ".docx",
		Size:      1026736,
	},
	{
		URL:       "https://github.com/thorstenkramm/fillfs/raw/refs/heads/main/samples/word_500kB.doc",
		FileName:  "word_500kB.doc",
		Extension: ".doc",
		Size:      503296,
	},
}

// SeedsByExtension returns all seeds with the given extension.
func SeedsByExtension(ext string) []Seed {
	result := make([]Seed, 0)
	for _, s := range All {
		if s.Extension == ext {
			result = append(result, s)
		}
	}
	return result
}
