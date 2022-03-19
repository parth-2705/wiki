# wiki
golang scraper for wikipedia

Parses p tags from content body of wikipedia articles and stores them in a badger db. Also recursively visits all the internal links of articles in the content body. If page is already visited, then doesn't visit the page again.
