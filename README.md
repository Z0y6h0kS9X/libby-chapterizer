# Libby-Chapterizer

## Description

Go utility for converting downloaded audiobooks from Libby/Overdrive into individual chapters rather than the standard file split. Example output it can generate:
- /Author/Series/[\#]. Title (ASIN)/*.mp3 - ASIN Match, Book in Series, output multiple files by chapter
- /Author/Series/Title/Title.m4b - No ASIN

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

Currently the program will take the source mp3 files, alongside the openbook.json file - and output mp3 files split based on chapters (or whatever the openbook.json specifies the splits should be).  It will then look up the book using author, narrator, and title using audnexus.  In the event that there are mulitple returned, it will use duration to determine the correct ASIN.  If it is enumerated, it will pull the metadata for that entry and store it - if it cannot be, it will use the local metadata supplied by the openbook.json file.  

## Roadmap

- [x] Parse JSON
- [x] Split based on openbook splits
- [ ] Generate M3U file
- [ ] Fallback to use ISBN and pull metadata from that, if No ASIN or if specified
- [x] Convert mp3 files into chapterized M4b (optional)
- [x] Look up book in Audible to pull metadata
- [x] Write metadata to m4b file

## Usage

### Arguments

| Flag                   | Shorthand | Default | Description                                                          |
|------------------------|:---------:|---------|----------------------------------------------------------------------|
| --json                 |     -j    |    ""   | The path to the openbook.json file                                   |
| --out                  |     -o    |    ""   | The path to the directory you want to output the files to            |
| --use-audible-chapters |     -c    |  false  | Specifies to override default breaks and use audible markers instead |
| --single               |     -s    |  false  | Specifies to output a single file (MP3 or M4B), instead to chapters  |
| --format               |     -f    |   MP3   | Specifies to output mp3 or m4b files                                 |

#### Default (outputs in same directory as files)
./libby-chapterizer-windows.exe --json <'path to json'>

#### Custom (outputs in custom directory)
./libby-chapterizer-windows.exe --json <'path to json'> --out <'output directory path'>

#### Custom (outputs in custom directory, uses audible chapters instead of openbook)
./libby-chapterizer-windows.exe --json <'path to json'> --out <'output directory path'> --use-audible-chapters

#### Custom (outputs in custom directory as a single m4b file)
./libby-chapterizer-windows.exe --json <'path to json'> --out <'output directory path'> --single --format m4b

## Contributing

Feel free to fork or open pull requests to help me out.

## Authors

Z0y6h0kS9X
