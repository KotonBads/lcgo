# LC Go

A Lightweight LunarClient launcher written in Go.
(currently for Linux and macOS only)

---

# Installation

## Go install
```sh
go install github.com/KotonBads/lcgo/lcgo@latest
```
or if you want a specific commit
```sh
go install github.com/KotonBads/lcgo/lcgo@commit hash
```

## Git
- Clone the repo
```sh
git clone https://github.com/KotonBads/lcgo.git lcgo
```
### VSCode:
- Open the folder in VSCode
- Run the build task included

### Manual:
- Requires Go 1.20 or newer
- Run this command in the terminal
```sh
mkdir build && go build -o build/lcgo lcgo/lcgo.go
```

LC Go should now be built inside the `build/` folder

---
# Usage
- Rename `example.config.json` to `config.json` or to any name
- Delete any keys you don't need which are marked as `optional`
- Fill in any required fields
- Run this command in the terminal
```sh
./path/to/lcgo/binary -config /path/to/config.json -version <insert minecraft version>
```
### Optional:
If you want the debug output, run this instead
```sh
./path/to/lcgo/binary -config /path/to/config.json -version <insert minecraft version> -debug=true
```

---