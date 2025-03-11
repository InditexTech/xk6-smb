# xk6-smb

The `xk6-smb` extension is a plugin for the k6 load testing tool that adds support for SMB (Samba) operations.

## Install

### Pre-built binaries

``` sh
make run
```

### Build from source

``` sh
make build
```

## Examples

See [examples](./examples) folder.

## Extension API

- `newClient`: creates a new SMB client.
- `getShares`: retrieve the list of shares from the server.
- `appendString`: append a string to a file existing in the server (if not exists, it will be created).
- `fileExists`: retrieve whether a file exists in the server.
- `readFile`: read a file contents from the server.
- `deleteFile`: delete a file from the server.
- `createDir`: create a directory in the server.
- `listFilesInDir`: list files in a directory in the server.
- `isDir`: check if a path is a directory.
- `dirExists`: check if a directory exists.
- `deleteFile`: delete a file from the server.
- `deleteDir`: delete a directory from the server, only if is empty.
- `copyFile`: copy a file from a local path to the server.
