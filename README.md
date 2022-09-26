# nvim-updater

> Install `nightly` nvim with the cli.

## Installation

```
go get github.com/notjrbauer/nvim-updater
```

## Usage

### `--flavor`

#### Flavor of distro to fetch.

#### Defaults: `macos`

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
