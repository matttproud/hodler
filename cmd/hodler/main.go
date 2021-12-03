// hodler converts iTerm2 color schemes into various formats.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"howett.net/plist"
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

	kernel = template.Must(
		template.New("Xresources").Funcs(template.FuncMap{"H": kernelHex}).Parse(
			`# Kernel Command Line: Generated with Hodler (http://github.com/matttproud/hodler)
vt.default_red={{H .ANSI0.Red}},{{H .ANSI1.Red}},{{H .ANSI2.Red}},{{H .ANSI3.Red}},{{H .ANSI4.Red}},{{H .ANSI5.Red}},{{H .ANSI6.Red}},{{H .ANSI7.Red}},{{H .ANSI8.Red}},{{H .ANSI9.Red}},{{H .ANSI10.Red}},{{H .ANSI11.Red}},{{H .ANSI12.Red}},{{H .ANSI13.Red}},{{H .ANSI14.Red}},{{H .ANSI15.Red}} vt.default_grn={{H .ANSI0.Green}},{{H .ANSI1.Green}},{{H .ANSI2.Green}},{{H .ANSI3.Green}},{{H .ANSI4.Green}},{{H .ANSI5.Green}},{{H .ANSI6.Green}},{{H .ANSI7.Green}},{{H .ANSI8.Green}},{{H .ANSI9.Green}},{{H .ANSI10.Green}},{{H .ANSI11.Green}},{{H .ANSI12.Green}},{{H .ANSI13.Green}},{{H .ANSI14.Green}},{{H .ANSI15.Green}} vt.default_blu={{H .ANSI0.Blue}},{{H .ANSI1.Blue}},{{H .ANSI2.Blue}},{{H .ANSI3.Blue}},{{H .ANSI4.Blue}},{{H .ANSI5.Blue}},{{H .ANSI6.Blue}},{{H .ANSI7.Blue}},{{H .ANSI8.Blue}},{{H .ANSI9.Blue}},{{H .ANSI10.Blue}},{{H .ANSI11.Blue}},{{H .ANSI12.Blue}},{{H .ANSI13.Blue}},{{H .ANSI14.Blue}},{{H .ANSI15.Blue}}
#
# No support for background; would be {{.Background}}.
#   Try using vt.color=0xXY where X and Y refer to color
#   offsets for background and foreground respectively, or potentially using
#   setterm -background with closely matching color.
#   X = {{.BackgroundIndex}}; Y = {{.ForegroundIndex}}
#
# No support for foreground; would be {{.Foreground}}.
#   Try using vt.color=0xXY where X and Y refer to color
#   offsets for background and foreground respectively, or potentially using
#   setterm -foreground with closely matching color.
#   X = {{.BackgroundIndex}}; Y = {{.ForegroundIndex}}
#
# No support for cursor color; would be {{.Cursor}}.
# No support for selected text color; would be {{.SelectedText}}.
# No support for selection color; would be {{.Selection}}.
# No support for cursor text coloring; would be {{.CursorText}}.
#
# No support for bold coloring; would be {{.Bold}}.
#   Try using vt.italic=X or vt.underline=X where X refers to
#   color offsets.
#   X = {{.BoldIndex}}
`))
)

func kernelHex(v float64) string {
	n := normalize(v)
	return fmt.Sprintf("0x%.2x", n)
}

type def struct {
	Blue  float64 `plist:"Blue Component"`
	Green float64 `plist:"Green Component"`
	Red   float64 `plist:"Red Component"`
}

func normalize(f float64) uint8 { return uint8(255 * f) }

func (d def) String() string {
	r := normalize(d.Red)
	g := normalize(d.Green)
	b := normalize(d.Blue)
	return fmt.Sprintf("%.2x%.2x%.2x", r, g, b)
}

type Table struct {
	ANSI0        def `plist:"Ansi 0 Color"`
	ANSI1        def `plist:"Ansi 1 Color"`
	ANSI2        def `plist:"Ansi 2 Color"`
	ANSI3        def `plist:"Ansi 3 Color"`
	ANSI4        def `plist:"Ansi 4 Color"`
	ANSI5        def `plist:"Ansi 5 Color"`
	ANSI6        def `plist:"Ansi 6 Color"`
	ANSI7        def `plist:"Ansi 7 Color"`
	ANSI8        def `plist:"Ansi 8 Color"`
	ANSI9        def `plist:"Ansi 9 Color"`
	ANSI10       def `plist:"Ansi 10 Color"`
	ANSI11       def `plist:"Ansi 11 Color"`
	ANSI12       def `plist:"Ansi 12 Color"`
	ANSI13       def `plist:"Ansi 13 Color"`
	ANSI14       def `plist:"Ansi 14 Color"`
	ANSI15       def `plist:"Ansi 15 Color"`
	Background   def `plist:"Background Color"`
	Bold         def `plist:"Bold Color"`
	Cursor       def `plist:"Cursor Color"`
	CursorText   def `plist:"Cursor Text Color"`
	Foreground   def `plist:"Foreground Color"`
	SelectedText def `plist:"Selected Text Color"`
	Selection    def `plist:"Selection Color"`
}

func (t *Table) findIndex(c def) (n int, ok bool) {
	switch c {
	case t.ANSI0:
		return 0, true
	case t.ANSI1:
		return 1, true
	case t.ANSI2:
		return 2, true
	case t.ANSI3:
		return 2, true
	case t.ANSI4:
		return 4, true
	case t.ANSI5:
		return 5, true
	case t.ANSI6:
		return 6, true
	case t.ANSI7:
		return 7, true
	case t.ANSI8:
		return 8, true
	case t.ANSI9:
		return 9, true
	case t.ANSI10:
		return 10, true
	case t.ANSI11:
		return 11, true
	case t.ANSI12:
		return 12, true
	case t.ANSI13:
		return 13, true
	case t.ANSI14:
		return 14, true
	case t.ANSI15:
		return 15, true
	default:
		return 0, false
	}
}

func (t *Table) ForegroundIndex() string {
	i, ok := t.findIndex(t.Foreground)
	if !ok {
		return "infeasible"
	}
	return fmt.Sprintf("%.1x", i)
}

func (t *Table) BackgroundIndex() string {
	i, ok := t.findIndex(t.Background)
	if !ok {
		return "infeasible"
	}
	return fmt.Sprintf("%.1x", i)
}

func (t *Table) BoldIndex() string {
	i, ok := t.findIndex(t.Bold)
	if !ok {
		return "infeasible"
	}
	return fmt.Sprint(i)
}

func DecodeInput(r io.ReadSeeker) (*Table, error) {
	d := plist.NewDecoder(r)
	var data Table
	return &data, d.Decode(&data)
}

func Output(w io.Writer, t *Table, tmpl *template.Template) error { return tmpl.Execute(w, t) }

type outputFormat string

const (
	unknown    = outputFormat("")
	Suckless   = outputFormat("Suckless")
	Xresources = outputFormat("Xresources")
	Alacritty  = outputFormat("Alacritty")
	Kernel     = outputFormat("Kernel")
)

type UnknownFormatError struct {
	Name string
}

func (err *UnknownFormatError) Error() string {
	if err == nil {
		return ""
	}
	return fmt.Sprintf("unknown format: %q", err.Name)
}

func (f outputFormat) String() string { return string(f) }
func (f *outputFormat) Set(v string) error {
	switch outputFormat(v) {
	case Suckless:
		*f = Suckless
	case Xresources:
		*f = Xresources
	case Alacritty:
		*f = Alacritty
	case Kernel:
		*f = Kernel
	default:
		return &UnknownFormatError{Name: v}
	}
	return nil
}

func (f outputFormat) template() *template.Template {
	switch f {
	case Suckless:
		return suckless
	case Xresources:
		return xresources
	case Alacritty:
		return alacritty
	case Kernel:
		return kernel
	default:
		panic("unhandled")
	}
}

func main() {
	var in, out string
	var format outputFormat
	flag.StringVar(&in, "in", "", "input source")
	flag.StringVar(&out, "out", "/dev/stdout", "output destination")
	flag.Var(&format, "output_format", "output format: 'Xresources' or 'Suckless' or 'Alacritty' or 'Kernel'")
	flag.Parse()
	if in == "" || out == "" {
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
	tmpl := format.template()
	if err := Output(fout, tab, tmpl); err != nil {
		log.Fatalln(err)
	}
}
