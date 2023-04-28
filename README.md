# tolino-readwise-sync
Syncs the highlights and annotations you made on your Tolino to Readwise.

## How to run
1. Using binary: Download the latest binary ([link](https://github.com/ssppooff/tolino-readwise-sync/releases/latest)) and run it with the following flags (or see help message
- `-t TOKEN_FILE`, where `TOKEN_FILE` is the path to the file containing your Readwise API token, see [Usage](#usage)
- `-n NOTES_FILE`, where `NOTES_FILE` is the path to the file you got off of your Tolino (usually `notes.txt`, located at the root of the Tolino's filesystem)

2. Using Docker image on GHCR: `docker pull ghcr.io/ssppooff/tolino-readwise-sync/tolino-sync:latest`
- Default command uses the flags `-t /files/token -n /files/notes.txt`, ie., it looks for the files `token` and `notes.txt` inside the container file system at `/folder`; the files contain the API token and notes & highlights, respectively.
- The command to run it using docker and a bind mount from the current working directory (assuming `token` & `notes.txt` inside `local_files/`) would be `docker run -v "$(pwd)"/local_files:/files tolino-sync:latest`

## Problem
Readwise can import your highlights and annotations from many other services, including Amazon's Kindle.

For Kindle devices, it supports importing directly from your [highlights page on amazon.com](https://read.amazon.com/notebook), but only for books purchased directly through Amazon. For all other books and documents, there is a file on Kindle devices called `My Clippings.txt` with all your highlights and annotations. Readwise let's you import that file either via email or uploading it to Readwise.

However, there is no import mechanism for Tolino devices, which are available in many parts of Europe (see [Wikipedia DE - Tolino](https://de.wikipedia.org/wiki/Tolino), and especially popular in German speaking countries (namely Germany, Austria, parts of Switzerland). Even though Tolino devices also use a local file to save highlights and annotations, its structure is different from Kindle's `My Clippings.txt`. And there is a cloud & app system available similar to Amazon's, however, it does not present your highlights and annotations on one page.

However, Readwise provides an [API](https://readwise.io/api_deets) which manages an account's highlights.

## Current state
- Only works if `notes.txt` file is in English (probably depends on which language is chosen for the Tolino)
- Mass upload of highlights & notes, irrespective of whether they have been already added
 or changed on the Tolino (Readwise does de-duplication based on `title/author/text/source_url`)
 
## Usage
- Get your API token [readwise.io/access_token](https://readwise.io/access_token) (subscription or trial necessary)
- Demo data Tolino `notes.txt` included in repo
- Read help message from CLI for which flags to use

## Roadmap/Milestones
1. Upload highlights & notes mode on your Tolino (done)
2. Don't add highlights & notes if they are already on Readwise
3. Like previous version, but if a note or highlight has been marked as modified on the Tolino, add corresponding metadata to Readwise

### Possible Features, depending on feasability
- Detect if a highlight has been deleted on the Tolino, delete it on Readwise (or make it an option to do so)
- Add metadata inside Tolino's `notes.txt`, not sure if it would mess with the Tolino
- Add highlights/notes as well as books/documents to Readwise's new [Reader](https://readwise.io/read)

## TODO
- Improve test coverage
