# nvim-updater

> Install `nightly` `macOS` (linux soon) nvim with the cli.

## Installation

```
go get github.com/notjrbauer/nvim-updater
```

## Usage

### `--source`

#### Directory to download neovim.

#### Defaults: `cwd`

### `--destination`

#### Directory to symlink neovim.

#### Defaults: `/usr/local/bin`, no need to specify nvim directly.

## Quick Setup

To quickly set up:

```bash
# Install to `$HOME` and symlink to `/usr/local/bin`
cd $HOME && nvim-updater
```
