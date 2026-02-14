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
    download/
      _index.md                # Download page (platform list)
  layouts/
    index.html                 # Landing page template
    404.html                   # Custom 404 page
    partials/
      head.html                # <head> with meta, OG tags, favicon
      navbar.html              # Sticky nav
      hero.html                # ASCII logo, tagline, install tabs
      separator.html           # Section separator
      story.html               # "Built for the AI era" section
      features.html            # Feature list
      highlights.html          # Stat cards grid
      footer.html              # Copyright, author, license
    download/
      list.html                # Download page layout
  static/
    CNAME                      # GitHub Pages custom domain
    css/style.css              # All styles
    js/tabs.js                 # Install tab switching
    favicon.svg                # SVG favicon
    robots.txt                 # Search engine directives
    install.sh                 # Symlink to repo root install.sh
```

## Common Tasks

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
