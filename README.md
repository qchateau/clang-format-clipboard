# clang-format-clipboard

clang-format-clipboard will run clang-format on your clipboard, optionally
keeping the indentation and leading newlines.

This in intended to format code in editors that do not support it:
copy your code, format the clipboard, paste it back.

## Build and install it

```bash
go get .
go build .
sudo cp clang-format-clipboard /usr/local/bin
```

You may also require to install a clipboard interface such as xclip

## Use in your DE

Configure a new shortcut in your DE.

Put a `.clang-format` file where your DE calls the executable. On Ubuntu 18.04, putting it in your home folder seems to work.
