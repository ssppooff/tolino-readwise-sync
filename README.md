# tolino-readwise-sync
Syncs the highlights and annotations you made on your Tolino to Readwise.

## Problem
Readwise can import your highlights and annotations from many other services, including Amazon's Kindle.

For Kindle devices, it supports importing directly from your [highlights page on amazon.com](https://read.amazon.com/notebook), but only for books purchased directly through Amazon. For all other books and documents, there is a file on Kindle devices called `My Clippings.txt` with all your highlights and annotations. Readwise let's you import that file either via email or uploading it to Readwise.

However, there is no import mechanism for Tolino devices, which are available in many parts of Europe (see [Wikipedia DE - Tolino](https://de.wikipedia.org/wiki/Tolino), and especially popular in German speaking countries (namely Germany, Austria, parts of Switzerland). Even though Tolino devices also use a local file to save highlights and annotations, its structure is different from Kindle's `My Clippings.txt`. And there is a cloud & app system available similar to Amazon's, however, it does not present your highlights and annotations on one page.

## Impelemtation
Readwise provides an [API](https://readwise.io/api_deets) which managing of an accounts highlights.

Possible Features, besides syncing your Tolino highlights to Readwise (one way for now)
- Add a specific tag to all highlights created by this app
- Don't add a highlight twice (Readwise does de-dupe based on 'title/author/text/source_url')
- (not sure if possible) in combination with Readwise's Reader, upload the original document alongside the highlights
- Detect if a highlight has been deleted on the Tolino, delete it on Readwise (or make it an option to do so)
- Verify that access token is valid (see API docs)

Possibly useful API commands
- Highlight CREATE -> keep note of each id under `modified_highlights`
- Highlight UPDATE
- Highlight DETAIL
- Highlights LIST, especially "filter by last updated datetime"
- Books LIST <- more like documents (categories include `books`, `articles`, `tweets`, `supplementals`, `podcasts`)
- Highlight Tag CREATE & Highlight Tags LIST


### Implementation Questions
- Is it possible to add an id-string to the note's file on the Tolino?
- Is API command `Books CREATE` missing? or does a book get created automatically if necessary?

### External Depencies

## Roadmap / Milestones
1. Check correct API token (done)
2. Get List of Highlites or Books

## TODO
## Concepts & Go things I need to learn
- How to create HTTPS POST & GET requests
- How to handle JSON objects
