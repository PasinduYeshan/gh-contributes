# gh-contributes

A minimal [GitHub CLI](https://cli.github.com/) extension that shows:
- **Daily contributions for the last 5 days**
- **Yearly contribution stats** (commits, PRs, issues, and more)

## Installation

```bash
gh extension install PasinduYeshan/gh-contributes
```

Then simply run:
```bash
gh contributes
```

## Usage

Sample output:
```
------------------------------------------
👋 Your GitHub Contributions (Last 5 days):
  2025-01-24: contributions: 3
  2025-01-23: contributions: 27
  2025-01-22: contributions: 17
  2025-01-21: contributions: 4
  2025-01-20: contributions: 13
------------------------------------------
👋 Your GitHub Contributions (Last Year):
 • Total Commits:              437
 • Total Pull Requests:        117
 • Total Pull Request Reviews: 44
 • Total Issues:               26
 • Total Repositories:         24
 • Private Contributions:      332
 • Overall Contributions:      980

------------------------------------------
```

## Development

1. Clone this repo and build:
   ```bash
   go build
   ```
2. Install locally to test:
   ```bash
   gh extension install .
   gh contributes
   ```

Enjoy! ✨