# Overview

[Hodler](https://en.wikipedia.org/wiki/Ferdinand_Hodler) converts
[iTerm 2](https://www.iterm2.com) color scheme [Property
Lists](https://en.wikipedia.org/wiki/Property_list) into forms that
the [Suckless Simple Terminal "st"](http://st.suckless.org) and
[X resources](https://en.wikipedia.org/wiki/X_resources)-based terminal
emulaters (e.g., [XTerm](http://invisible-island.net/xterm/)) and
[Alacritty](https://github.com/jwilm/alacritty) can consume.

# Installation

Hodler is built using the [Go Programming Language](https://golang.org).  Go
is required to build and modify the tool.

    $ go get -u github.com/matttproud/hodler/cmd/hodler

Go generates staticly linked binaries, so users of Hodler needn't have Go
installed for casual use.

# Usage

Users of st can generate a fragment to embed into their local `config.h`.

    $ hodler -in adio.itermcolors -out config.h -output_format Suckless

Users of XTerm and other X resources systems will fancy this:

    $ hodler -in adio.itermcolors -out Xresources -output_format Xresources

Users of Alacritty can generate a YAML fragment to embed into their local `alacritty.yml`:

    $ hodler -in adio.itermcolors -out alacritty.yml -output_format Alacritty

# Examples

Borland theme in ST:


![Demo of ST with Borland Theme](Demo_ST_Borland.png)


Borland theme in XTerm:


![Demo of XTerm with Borland Theme](Demo_XTerm_Borland.png)

Borland theme in Alacritty:


![Demo of Alacritty with Borland Theme](Demo_Alacritty_Borland.png)


Reference case of Borland theme in iTerm 2:


![Demo of iTerm 2 with Borland Theme](Demo_iTerm2_Borland.png)


# Build Status

[![Build Status](https://travis-ci.org/matttproud/hodler.svg)](https://travis-ci.org/matttproud/hodler)
