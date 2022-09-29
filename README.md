# ls

> TODO: Rewrite that in Zig.

A little command line utility that switches `-a` and `-A` flags before proxying standard I/O to ls.

It also tweaks the help command output to match the new flags.

```diff
...
- -a, --all                  do not ignore entries starting with .
- -A, --almost-all           do not list implied . and ..
+ -a, --almost-all           do not list implied . and ..
+ -A, --all                  do not ignore entries starting with .
...

```

## Installation
```bash
curl -L https://github.com/felixdorn/ls/releases/latest/download/ls -o /usr/bin/ls-proxy
chmod +x /usr/bin/ls-proxy

# In your ~/.bashrc / ~/.zshrc / whatever
alias ls="/usr/bin/ls-proxy"
```


## Building

```bash
go build -ldflags="-s -w" -gcflags=all='-l'
upx ls # if you're into that
```
