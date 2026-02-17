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
    _index.md                  # Landing page (tagline, install, highlights)
    download/
      _index.md                # Download page (platform list)
  layouts/
    index.html                 # Landing page template
    404.html                   # Custom 404 page
    partials/
      head.html                # <head> with meta, OG tags, favicon
      navbar.html              # Sticky nav
      hero.html                # ASCII logo, tagline, install command, onboarding video
      showcase.html            # Video showcase grid (4 videos)
      highlights.html          # Stat cards grid
      closing.html             # CTA section
      separator.html           # Section separator (unused in current layout)
      footer.html              # Copyright, author, license
    download/
      list.html                # Download page layout
  static/
    CNAME                      # GitHub Pages custom domain
    css/style.css              # All styles
    js/tabs.js                 # Install tab switching + copy button
    favicon.svg                # SVG favicon
    robots.txt                 # Search engine directives
    install.sh                 # Symlink to repo root install.sh
    video/
      onboarding.mp4           # Hero demo video
      play-menu.mp4            # Mode selection showcase
      gameplay.mp4             # Gameplay showcase
      practice.mp4             # Practice mode showcase
      statistics.mp4           # Statistics showcase
      posters/                 # Fallback poster images for each video
```

## Videos

5 MP4 videos encoded as H.264 Baseline Level 4.0 (1280px wide, 30fps) for broad mobile compatibility. Each `<video>` tag uses `autoplay loop muted playsinline` with a `poster` fallback image from `/video/posters/`.

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
- `install.curl` — install command shown in hero
- `highlights` — stat cards
