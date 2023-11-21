# Libby-Chapterizer

## Description

Go utility for converting downloaded audiobooks from Libby/Overdrive into individual chapters rather than the standard file split.  It will place the files in ../Author/Series/[\#]. Title (ASIN)/*

## Table of Contents

- [Installation](#installation)
- [Status](#status)
- [Usage](#usage)
- [Roadmap](#roadmap)
- [Contributing](#contributing)
- [Authors](#authors)

## Installation

Download the binary for the platform you want to use, then navigate the terminal of choice to the directory containing the binary.

## Status

Currently the program will take the source mp3 files, alongside the openbook.json file - and output mp3 files split based on chapters (or whatever the openbook.json specifies the splits should be).  It will generate a VLC-compatible m3u playlist file as well.

## Roadmap

- [x] Parse JSON
- [x] Split based on openbook splits
- [x] Generate M3U file
- [ ] Convert mp3 files into chapterized M4b (optional)
- [x] Look up book in Audible to pull metadata
- [ ] Write metadata to m4b file

## Usage

### Arguments

| Flag                   | Shorthand | Default | Description                                                          |
|------------------------|:---------:|---------|----------------------------------------------------------------------|
| --json                 |     -j    |    ""   | The path to the openbook.json file                                   |
| --out                  |     -o    |    ""   | The path to the directory you want to output the files to            |
| --use-audible-chapters |     -c    |  false  | Specifies to override default breaks and use audible markers instead |

#### Default (outputs in same directory as files)
./libby-chapterizer-windows.exe --json <'path to json'>

#### Custom (outputs in custom directory)
./libby-chapterizer-windows.exe --json <'path to json'> --out <'output directory path'>

#### Custom (outputs in custom directory, uses audible chapters instead of openbook)
./libby-chapterizer-windows.exe --json <'path to json'> --out <'output directory path'> --use-audible-chapters

## Contributing

Feel free to fork or open pull requests to help me out.

## Authors

Z0y6h0kS9X
