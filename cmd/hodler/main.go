// hodler converts iTerm2 color schemes to X resources or Suckless ST terminal config.h definitions.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	plist "github.com/DHowett/go-plist"
)

var (
	xresources = template.Must(
		template.New("Xresources").Parse(
			`! X Resources: Generated with Hodler (http://github.com/matttproud/hodler)
*color0: #{{.ANSI0}}
*color1: #{{.ANSI1}}
*color2: #{{.ANSI2}}
*color3: #{{.ANSI3}}
*color4: #{{.ANSI4}}
*color5: #{{.ANSI5}}
*color6: #{{.ANSI6}}
*color7: #{{.ANSI7}}
*color8: #{{.ANSI8}}
*color9: #{{.ANSI9}}
*color10: #{{.ANSI10}}
*color11: #{{.ANSI11}}
*color12: #{{.ANSI12}}
*color13: #{{.ANSI13}}
*color14: #{{.ANSI14}}
*color15: #{{.ANSI15}}
*background: #{{.Background}}
*foreground: #{{.Foreground}}
*cursorColor: #{{.Cursor}}
! See "highlightColorMode" and "hm" options in XTerm manual page.
*highlightTextColor: #{{.SelectedText}}
*highlightColor: #{{.Selection}}
! No support for cursor text coloring; would be #{{.CursorText}}.
! No support for bold coloring; would be #{{.Bold}}.
`))

	suckless = template.Must(
		template.New("Suckless").Parse(
			`/* Suckless ST config.h Fragment
 * Generated with Hodler (http://github.com/matttproud/hodler)
 */
static const char *colorname[] = {
	"#{{.ANSI0}}",		/* 0: ANSI Color 0 */
	"#{{.ANSI1}}",		/* 1: ANSI Color 1 */
	"#{{.ANSI2}}",		/* 2: ANSI Color 2 */
	"#{{.ANSI3}}",		/* 3: ANSI Color 3 */
	"#{{.ANSI4}}",		/* 4: ANSI Color 4 */
	"#{{.ANSI5}}",		/* 5: ANSI Color 5 */
	"#{{.ANSI6}}",		/* 6: ANSI Color 6 */
	"#{{.ANSI7}}",		/* 7: ANSI Color 7 */
	"#{{.ANSI8}}",		/* 8: ANSI Color 8 */
	"#{{.ANSI9}}",		/* 9: ANSI Color 9 */
	"#{{.ANSI10}}",		/* 10: ANSI Color 10 */
	"#{{.ANSI11}}",		/* 11: ANSI Color 11 */
	"#{{.ANSI12}}",		/* 12: ANSI Color 12 */
	"#{{.ANSI13}}",		/* 13: ANSI Color 13 */
	"#{{.ANSI14}}",		/* 14: ANSI Color 14 */
	"#{{.ANSI15}}",		/* 15: ANSI Color 15 */
	[255] = 0,
	[256] = "#{{.Background}}",		/* 256: Background */
	[257] = "#{{.Foreground}}",		/* 257: Foreground */
	[258] = "#{{.Cursor}}",		/* 258: Cursor */
	[259] = "#{{.CursorText}}",		/* 259: Cursor Text */
	/* No support for text highlight coloring; would be #{{.SelectedText}}. */
	/* No support for highlight coloring; would be #{{.Selection}}. */
	/* No support for bold coloring; would be #{{.Bold}}. */
};

static unsigned int defaultfg  = 257;
static unsigned int defaultbg  = 256;
static unsigned int defaultcs  = 258;
static unsigned int defaultrcs = 259;
`))

	alacritty = template.Must(
		template.New("Alacritty").Parse(
			`# Alacritty Colors
# Generated with Hodler (http://github.com/matttproud/hodler)
colors:
  primary:
    background: '0x{{.Background}}'
    foreground: '0x{{.Foreground}}'

  cursor:
    text: '0x{{.CursorText}}'
    cursor: '0x{{.Cursor}}'

  normal:
    black:   '0x{{.ANSI0}}'
    red:     '0x{{.ANSI1}}'
    green:   '0x{{.ANSI2}}'
    yellow:  '0x{{.ANSI3}}'
    blue:    '0x{{.ANSI4}}'
    magenta: '0x{{.ANSI5}}'
    cyan:    '0x{{.ANSI6}}'
    white:   '0x{{.ANSI7}}'

  bright:
    black:   '0x{{.ANSI8}}'
    red:     '0x{{.ANSI9}}'
    green:   '0x{{.ANSI10}}'
    yellow:  '0x{{.ANSI11}}'
    blue:    '0x{{.ANSI12}}'
    magenta: '0x{{.ANSI13}}'
    cyan:    '0x{{.ANSI14}}'
    white:   '0x{{.ANSI15}}'
# No support for text highlight coloring; would be 0x{{.SelectedText}}.
# No support for highlight coloring; would be 0x{{.Selection}}.
# No support for bold coloring; would be 0x{{.Bold}}.
`))
)

type Defn struct {
	Blue  float64 `plist:"Blue Component"`
	Green float64 `plist:"Green Component"`
	Red   float64 `plist:"Red Component"`
}

func Normalize(f float64) uint8 { return uint8(255 * f) }

func (d Defn) String() string {
	r := Normalize(d.Red)
	g := Normalize(d.Green)
	b := Normalize(d.Blue)
	return fmt.Sprintf("%.2x%.2x%.2x", r, g, b)
}

type Table struct {
	ANSI0        Defn `plist:"Ansi 0 Color"`
	ANSI1        Defn `plist:"Ansi 1 Color"`
	ANSI2        Defn `plist:"Ansi 2 Color"`
	ANSI3        Defn `plist:"Ansi 3 Color"`
	ANSI4        Defn `plist:"Ansi 4 Color"`
	ANSI5        Defn `plist:"Ansi 5 Color"`
	ANSI6        Defn `plist:"Ansi 6 Color"`
	ANSI7        Defn `plist:"Ansi 7 Color"`
	ANSI8        Defn `plist:"Ansi 8 Color"`
	ANSI9        Defn `plist:"Ansi 9 Color"`
	ANSI10       Defn `plist:"Ansi 10 Color"`
	ANSI11       Defn `plist:"Ansi 11 Color"`
	ANSI12       Defn `plist:"Ansi 12 Color"`
	ANSI13       Defn `plist:"Ansi 13 Color"`
	ANSI14       Defn `plist:"Ansi 14 Color"`
	ANSI15       Defn `plist:"Ansi 15 Color"`
	Background   Defn `plist:"Background Color"`
	Bold         Defn `plist:"Bold Color"`
	Cursor       Defn `plist:"Cursor Color"`
	CursorText   Defn `plist:"Cursor Text Color"`
	Foreground   Defn `plist:"Foreground Color"`
	SelectedText Defn `plist:"Selected Text Color"`
	Selection    Defn `plist:"Selection Color"`
}

func DecodeInput(r io.ReadSeeker) (*Table, error) {
	d := plist.NewDecoder(r)
	var data Table
	return &data, d.Decode(&data)
}

func GetTmpl(n string) *template.Template {
	switch n {
	case "Xresources":
		return xresources
	case "Suckless":
		return suckless
	case "Alacritty":
		return alacritty
	default:
		panic("unhandled")
	}

}

func Output(w io.Writer, t *Table, tmpl *template.Template) error { return tmpl.Execute(w, t) }

func main() {
	var in, out, outputFormat string
	flag.StringVar(&in, "in", "", "input source")
	flag.StringVar(&out, "out", "/dev/stdout", "output destination")
	flag.StringVar(&outputFormat, "output_format", "", "output format: 'Xresources' or 'Suckless' or 'Alacritty'")
	flag.Parse()
	if in == "" || out == "" || (outputFormat != "Suckless" && outputFormat != "Xresources" && outputFormat != "Alacritty") {
		flag.Usage()
		os.Exit(1)
	}
	fin, err := os.Open(in)
	if err != nil {
		log.Fatalln(err)
	}
	defer fin.Close()
	fout, err := os.Create(out)
	if err != nil {
		log.Fatalln(err)
	}
	defer fout.Close()
	tab, err := DecodeInput(fin)
	if err != nil {
		log.Fatalln(err)
	}
	tmpl := GetTmpl(outputFormat)
	if err := Output(fout, tab, tmpl); err != nil {
		log.Fatalln(err)
	}
}
