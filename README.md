# g

---
[![GitHub stars](https://img.shields.io/github/stars/Equationzhao/g)](https://github.com/Equationzhao/g/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/Equationzhao/g)](https://github.com/Equationzhao/g/network)
[![GitHub issues](https://img.shields.io/github/issues/Equationzhao/g)](https://github.com/Equationzhao/g/issues)
[![GitHub license](https://img.shields.io/github/license/Equationzhao/g)](https://github.com/Equationzhao/g/blob/master/LICENSE)

A feature-rich, customizable, and cross-platform `ls` alternative.

Experience enhanced visuals with type-specific icons, various layout options, and git status integration.

---

## Key Features

1. **Customizable Display**: Icons and colors specific to file types, easy to customize.
2. **Multiple Layouts**: Choose from grid, across, byline, zero, comma, table, json, markdown, and tree layouts.
3. **Git Integration**: View file git-status/repo-status/repo-branch directly in your listings.
4. **Advanced Sorting**: Highly customizable sorting options like version-sort.
5. **Cross-Platform Compatibility**: Works seamlessly on Linux, Windows, and MacOS.
6. **Fuzzy Path Matching**: [`zoxide`](https://github.com/ajeetdsouza/zoxide) and [`fzf`](https://github.com/junegunn/fzf) like fuzzy path matching.
7. **Hyperlink support**: Open files/directories with a single click.

## Screenshots

![image](asset/screenshot_3.png)

## Usage

```bash
g path(s)
```

```bash
g --icon --long path(s) # show icons and long format
```

```bash
g --tree --long path(s) # show tree layout
```

## More options

[man.md](docs/man.md)

## Installation Guide

### Via package manager

#### Arch Linux (AUR)

```bash
yay -S g-ls
```

#### Homebrew

```bash
brew install g-ls
```
or use the homebrew tap:

```bash
brew tap equationzhao/core git@github.com:Equationzhao/homebrew-g.git
```

```bash
brew install g-ls
```

#### MacPort

```bash
sudo port install g-ls
```

#### Windows

windows scoop:

```powershell
scoop install https://raw.githubusercontent.com/Equationzhao/g/master/scoop/g.json
```

```powershell
# upgrade
scoop uninstall g # uninstall first
scoop install https://raw.githubusercontent.com/Equationzhao/g/master/scoop/g.json
```

#### Winget

TODO, see [issue](https://github.com/Equationzhao/g/issues/119)

### Pre-built executable

#### install script

##### install
```sh
bash -c "$(curl -fsSLk https://raw.githubusercontent.com/Equationzhao/g/master/script/install.sh)"
```

##### uninstall
```sh
curl -fsSLk https://raw.githubusercontent.com/Equationzhao/g/master/script/install.sh | bash /dev/stdin -r     
```

#### deb

download from [release](https://github.com/Equationzhao/g/releases) page

```bash
sudo dpkg -i g_$version_$arch.deb
```

#### tar.gz/zip

just download from [release page](https://github.com/Equationzhao/g/releases), extract the gzip and add the executable file to your `PATH`

### From source

Requires Go version >= 1.26.1

```bash
go install -ldflags="-s -w"  github.com/Equationzhao/g@latest
```

Alternatively, clone the repo for a dev version:

```bash
git clone github.com/Equationzhao/g
cd g
go build -ldflags="-s -w" 
# then add the executable file to your `PATH`
```

#### Build options
See [BuildOption.md](docs/BuildOption.md) for optional build tags to enable or disable features.

## Recommended terminal

macOS:
- [Iterm2](https://iterm2.com/)
- [Warp](https://www.warp.dev)

Windows:
- [Windows Terminal](https://github.com/microsoft/terminal)

cross-platform:
- [Hyper](https://hyper.is/)
- [WezTerm](https://wezfurlong.org/wezterm/index.html)


## Shell Integration

### completion

>> *if you install `g` through brew or the install script, the completion is usually installed already.*

#### zsh
```zsh
wget https://raw.githubusercontent.com/Equationzhao/g/master/completions/zsh/_g
```

install the file to your zsh completion directory, usually `/usr/local/share/zsh/site-functions` or `/usr/share/zsh/site-functions` (or anywhere in your $FPATH)

```zsh
mv _g ~/.zsh/completions
```

make sure `autoload -Uz compinit` and `compinit` are in the `~/.zshrc` or `~/.zprofile`

if not, add them to at least one of them.

```zsh
autoload -Uz compinit
compinit
```

#### bash

```bash
wget https://raw.githubusercontent.com/Equationzhao/g/master/completions/bash/g-completion.bash
```

add the following lines to your ~/.bashrc file:

```bash
source /path/to/g-completion.bash
```

#### fish

```fish
wget https://raw.githubusercontent.com/Equationzhao/g/master/completions/fish/g.fish
```

Install the file to your fish completion directory, usually ~/.config/fish/completions

```fish
mv g.fish ~/.config/fish/completions
```

Restart your terminal session or run the following command to immediately enable the completion functionality:

```fish
source ~/.config/fish/config.fish
```

### alias

Generate initialization scripts(alias) for various shells:

```bash
g -init bash/zsh/fish/pwsh
```

##### bash

```.bash
# add the following command to .bashrc
eval "$(g --init bash)"
# then `source ~/.bashrc`
```

##### zsh

```zsh
# add the following command to .zshrc
eval "$(g --init zsh)"
# then `source ~/.zshrc`
```

##### fish

```fish
#  add to fish config:
g --init fish | source
#  then `source ~/.config/fish/config.fish`
```

##### powershell

```powershell
# add the following line to your profile
Invoke-Expression (& { (g --init powershell | Out-String) })
```

use command `echo $profile` to find your profile path

##### nushell

the nushell has a nice built-in ls command, but if you wanna try `g` in nushell, you can do the following:

ps: the script is not guaranteed to work, if you have any problem, please [file an issue](https://github.com/Equationzhao/g/issues/new/choose)

```nu
# add the following to your $nu.env-path
^g --init nushell | save -f ~/.g.nu
# then add the following to your $nu.config-path
source ~/.g.nu

# if you want to replace nushell's g command with g
# add the following definition and alias to your $nu.config-path
#
# def nug [arg?] {
#     if ($arg == null) {
#         g $arg
#     } else {
#         g
#     }
# }
# alias g = ^g
```

## Custom theme

[theme](docs/Theme.md)

## CONTRIBUTING

Interested in contributing? Check out the [contributing guidelines](./CONTRIBUTING.md).

## Alternatives

`g` is highly inspired by following projects that you may wanna try!

- [exa](https://github.com/ogham/exa) or [eza](https://github.com/eza-community/eza)
- [lsd](https://github.com/lsd-rs/lsd)
- [ls-go](https://github.com/acarl005/ls-go)

|                | eza                                                                                           | g                                                                            |
|----------------|-----------------------------------------------------------------------------------------------|------------------------------------------------------------------------------|
| display mode   | oneline,grid,across,tree,recurse                                                              | oneline,grid,across,zero,comma,table,json,markdown,tree,recurse              |
| unique feature | -Z: list each file’s security context,-@: list each file’s extended attributes and sizes ...  | --mime: list each file's mime type, --charset: list each file's charset  ... |
| performance    | better                                                                                        | slower                                                                       |

## Star History

[![Star History Chart](https://api.star-history.com/svg?repos=Equationzhao/g&type=Date)](https://www.star-history.com/#Equationzhao/g&Date)
