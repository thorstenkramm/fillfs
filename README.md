# fillfs

Fill my file system with realistic random files.

Files contain real but demo content. All content is harmless.

Fillfs is for testing scenarios where you quickly need many files.

## Usage

```bash
./fillfs --dest ./fakefs \
  --cache-dir /tmp/fillfs \
  --clean-cache false \
  --folders 10 \
  --files-per-folder 100 -\
  -depths 5
```

The above command will do the following:

1. Download sample source files from [github.com/thorstenkramm/fillfs](https://github.com/thorstenkramm/fillfs/samples)
2. Stores the source files in the cache directory `/tmp/fillfs`. Directory will be created if missing.
3. Create 10 folders and 100 files in the directory `./fakefs`. Directory will be created if missing.
4. Changes into each of the 10 directories at top level and repeats step 3 five times (depth).

How many files and how many folders will this command create?

## fillfs File and Folder Calculation

Based on the described behavior (recursive into each folder at each depth):

### Example

`--folders 10 --files-per-folder 100 --depths 5`

| Depth | Folders     | Files |
|-------|-------------|-------|
| 1     | 10          | 1,000 |
| 2     | 100         | 10,000 |
| 3     | 1,000       | 100,000 |
| 4     | 10,000      | 1,000,000 |
| 5     | 100,000     | 10,000,000 |
|Total  | **111,110** | **11,111,000** |

### Formulas

```
Total folders = f × (f^d - 1) / (f - 1)
Total files   = Total folders × n
```

Where:
- `f` = folders (--folders)
- `d` = depths (--depths)
- `n` = files-per-folder (--files-per-folder)

> [!CAUTION]
> This exponential growth can produce very large file counts quickly.
> Depth accepts floats, so consider a depth of 2.2, for example, and fewer files per folder.

## Default settings

If you invoke `./fillfs` without any arguments, the following default settings will apply:

- `dest`: current directory
- `cache-dir`: default temporary folder of your operating system, with a fallback to `/tmp`.
- `clean-cache`: false
- `folders`: 2
- `files-per-folder`: 20
- `depth`: 1

## Behaviour

Before doing anything, the command-line utility `fillfs` first calculates the resulting number of files
and folders and the approximate disk space required. If there is not enough free disk space the program exits
with an error (code 3).

After presenting the above calculation results, the user must confirm that they wish to continue.
You can bypass the interactive confirmation with `--yes`.

The cache directory will be "marked" as fillfs-created by storing a hidden empty file `.fillfs` at top level.
The command-line utility will exit with an error (code 4) if you attempt to use an existing directory as a cache directory that has not been created by `fillfs`.

If the destination directory exists and it's not empty, `./fillfs` will exit with an error (code 5). You can change this behaviour by using `--wipe-dest` which will cause fillfs to delete all files and folders from the destination first.

## Using as a go module

You can use fillfs directly in your Go project and inside your Go unit tests.
See the example below:

```go
// fill me
```
