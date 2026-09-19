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
    - title match
    - title non-match, as in "but not if title matches"
    - article URL (e.g. filter "/sport/" articles)
- and combinations

In progress:
- Filter based on
    - content (e.g. filter newspaper "plus" articles)
- Documentation
- More examples
- Daemon mode
- Container image


## Example config

```yaml
# ~/.config/miniflux-sieve/config.yml
MINIFLUX_URL: "https://miniflux.example.com"
MINIFLUX_API_KEY: "My secret token"  # Settings → API Keys → Create a new API key
```


## Example rules

```yaml
# ~/.config/miniflux-sieve/rules.yml
version: 1
rules:
  - name: Expire old Ars Technica entries
    feed: "http://feeds.arstechnica.com/arstechnica/index"
    when:
      older_than: 14d
  - name: Not interested in Patreon
    feed: "https://example.com/feed"
    when:
      tagged: "patreon"
  - name: Expire old HN entries, with exceptions
    feed: "https://hnrss.org/newest?points=200"
    when:
      older_than: 30d
      not_title_matches:
        - '(?-i)Elmo\b'  # Elmo is case-sensitive, and must end with a word boundary
        - ernie
        - bert
        - oscar
  - name: Not interested in sport news
    feed: "https://www.tagesspiegel.de/contentexport/feed/home"
    when:
      url_matches: "/sport/"
```
