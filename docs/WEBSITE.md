# Website

Custom Hugo site with no external theme. Dark theme, monospace-accented, minimal.

## Commands

```bash
cd website
hugo server          # Local dev server on localhost:1313
hugo --minify        # Production build to public/
```

## Structure

```
website/
  hugo.toml                    # Site config, menu, latestVersion
  content/
    _index.md                  # Landing page (tagline, features, highlights, story)
    changelog/
      _index.md                # Changelog list page
      v0.1.0.md                # Release entry
    download/
      _index.md                # Download page (Windows platform list)
  layouts/
    index.html                 # Landing page template
    partials/
      head.html                # <head> with meta, OG tags, favicon
      navbar.html              # Sticky nav
      hero.html                # ASCII logo, tagline, install tabs
      story.html               # "Built for the AI era" section
      features.html            # Feature list
      highlights.html          # Stat cards grid
      footer.html              # Copyright, author, license
    changelog/
      list.html                # Changelog list layout
      single.html              # Individual release layout
    download/
      list.html                # Download page layout
  static/
    css/style.css              # All styles
    js/tabs.js                 # Install tab switching
    favicon.svg                # SVG favicon
    robots.txt                 # Search engine directives
    install.sh                 # Install script (copied from repo root)
  archetypes/
    changelog.md               # Template for new changelog entries
```

## Common Tasks

### Add a changelog entry

```bash
cd website
hugo new changelog/v0.2.0.md
```

Then edit `content/changelog/v0.2.0.md` with the release notes.

### Update the latest version

Edit `hugo.toml` — change one line:

```toml
latestVersion = '0.2.0'
```

This updates the download page URLs across the site.

### Edit landing page copy

All landing page content lives in `content/_index.md`:

- `tagline` — hero subtitle
- `features` — feature list items
- `highlights` — stat cards
- `story` — "Built for the AI era" section

### Update the install script

The install script at `static/install.sh` is a copy of the repo root `install.sh`. If the root script changes, copy it:

```bash
cp install.sh website/static/install.sh
```
