# clang-format-clipboard

clang-format-clipboard will strip whitespaces, call clang-format and restore the original indentation on whatever is in your clipboard.

## Build and install it

```bash
go build .
sudo cp clang-format-clipboard /usr/local/bin
```

You may also require to install a clipboard interface (for instance: xclip)

## Use in your DE

Configure a new shortcut in your DE.

I bind `Ctrl+Shift+F` to `clang-format-clipboard -executable=clang-format-9`.

Put a `.clang-format` file where your DE calls the executable. On Ubuntu 18.04, putting it in your home folder seems to work.
