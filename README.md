# projpak

`projpak` is a command-line utility written in Go that allows you to flatten files from a directory with specific extensions into a single file and reconstruct those files from the flattened format. This tool is useful for packaging files in a simple text-based format and later extracting them back into their original structure.

## Features

- **Flatten Files**: Combine multiple files with specified extensions into a single file using a custom format.
- **Reconstruct Files**: Extract files from the flattened format and recreate the original directory structure.

## Installation

Ensure you have Go installed on your system. You can install `projpak` using the following command:

```bash
go install github.com/nokusukun/projpak@latest
```

## Usage

`projpak` provides two main commands: `flatten` and `reconstruct`.

### Flatten Command

The `(f)latten` command combines files with specified extensions from a directory into a single output file.

#### Syntax

```bash
projpak f -dir <directory> -ext <extensions> -output <output_file>
```

#### Options

- `-dir`: The directory to search for files.
- `-ext`: Comma-separated list of file extensions to include (e.g., `.go,.txt`).
- `-output`: The output file to write the flattened contents to.

#### Example

```bash
projpak f -dir /path/to/directory -ext .go,.txt -output output.txt
```

This command will search `/path/to/directory` for `.go` and `.txt` files, flatten their contents into `output.txt`.

### Reconstruct Command

The `(r)econstruct` command reads a flattened file and reconstructs the original files in a specified directory.

#### Syntax

```bash
projpak r -file <flattened_file> -directory <output_directory>
```

#### Options

- `-file`: The flattened file to read and reconstruct files from.
- `-directory`: The directory where the reconstructed files will be unpacked.

#### Example

```bash
projpak r -file output.txt -directory /path/to/unpack
```

This command will read `output.txt` and recreate the original files inside `/path/to/unpack`.

## File Format

The flattened file uses a simple schema:

```
<file path=/path/to/filename.go>
{file_content}
</file>
<file path=/path/to/filename-2.go>
{file_content}
</file>
```

Each file's content is enclosed within `<file path=...>` and `</file>` tags, specifying the original path and content.

## License

This project is licensed under the MIT License.

## Contributing

Contributions are welcome! Please open issues or submit pull requests with improvements.