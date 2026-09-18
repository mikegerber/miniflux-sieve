# miniflux-sieve

Mark unwanted entries read in your Miniflux RSS reader. The elderly AKA Gen X would call it a
killfile, but for Miniflux.

There a few solutions for this, but I wanted more expressive filters, to do something like
"mark everything older than 7 days as read, unless the title contains any of 'ernie', 'bert',
'oscar' or 'tiffy'".

## Progress

Currently working functionality:
- Filter based on
    - age (e.g. anything older than 7 days)
    - tags (e.g. "Sponsor")
- and combinations

In progress:
- Filter based on
    - article URL (e.g. filter "/sport/" articles)
    - article title
    - content (e.g. filter newspaper "plus" articles)
- Documentation
- More examples
- Daemon mode
- Docker container


## Example rules

```yaml
# ~/.config/miniflux-sieve/rules.yml
version 1:
rules:
  - name: Expire old Ars Technica entries
    feed: "http://feeds.arstechnica.com/arstechnica/index"
    when:
      older_than: 14d
  - name: Not interested in Patreon
    feed: "https://example.com/feed"
    when:
      tagged: "patreon"
```
