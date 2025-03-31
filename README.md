# xk6-smb

An extension for the k6 load testing tool that adds support for SMB (Samba) operations. This extension allows you to perform operations such as uploading, downloading and managing files on a Samba server as part of your load testing scripts.

## Install

To build a `k6` binary with this extension, first ensure you have the prerequisites:

- [Go toolchain](https://go101.org/article/go-toolchain.html)
- Git

Then:

1. Download [xk6](https://github.com/grafana/xk6):
```bash
go install go.k6.io/xk6/cmd/xk6@latest
```

2. [Build](https://github.com/grafana/xk6#command-usage) the k6 binary:
```bash
xk6 build --with github.com/InditexTech/xk6-smb@latest
```

### Development

For building a `k6` binary with the extension from the local code, you can run:

```bash
make build
```

For testing and running the extension locally, a Samba server is required. The default target in the Makefile will:

- Run an SMB server with Docker Compose (make sure you have [Docker](https://docs.docker.com/engine/install/) & [Docker Compose](https://docs.docker.com/compose/install/) installed in your system).
- Download the dependencies.
- Format your code.
- Run the integration tests.
- Run the [example](examples/main.js) script.

```bash
git clone git@github.com:InditexTech/xk6-smb.git
cd xk6-smb
make
```

## Usage

This extension provides the following JS methods for interacting with the SMB server:

```javascript
import xk6smb from "k6/x/smb";

// Create a new SMB client
const client = xk6smb.newClient("host:port", "username", "password", "Share name");

export default function() {
  // Retrieve the list of shares from the server (regardless of the Share provided in the client).
  const shares = client.getShares();
  // Output example: {"success":true,"message":"","shares":["Mount","User Volume","IPC$"]}

  // Copy a file from a local path to a remote path in the server share.
  const copyResult = client.copyFile("./examples/test-data/test-file.txt", "folder/copied-file.txt");
  // Output example: {"success":true,"message":"File copied successfully"}

  // Append a string to a file in the server share (if the file does not exist, it will be created).
  const writeResult = client.appendString("folder/file.txt", "Lorem ipsum");
  // Output example: {"success":true,"message":"Bytes appended successfully"}

  // Check if a file exists in the server share.
  const isFile = client.fileExists("folder/file.txt");
  // Output: true/false

  // Read the content of a file from the server share.
  const readResult = client.readFile("folder/file.txt");
  // Output: the content of the file

  // Create a directory in the server share.
  const createDirResult = client.createDir("foo/bar/baz");
  // Output example: {"success":true,"message":"Folder created successfully"}

  // List the files existing in directory in the server share.
  const listResult = client.listFilesInDir("foo/bar");
  // Output: array with the existing files names. Example: ["test_0_1.txt"]

  // Check if a given path is a directory in the server share.
  const isDir = client.isDir("foo/bar");
  // Output: true/false

  // Check if a directory exists in the server share.
  const dirExists = client.dirExists("foo/bar/baz");
  // Output: true/false

  // Delete a single file from the server share.
  const deleteFileResult = client.deleteFile("folder/file.txt");
  // Output example: {"success":true,"message":"File deleted successfully"}

  // Delete a directory from the server share
  const deleteDirResult = client.deleteDir("foo/bar/baz");
  // Output example: {"success":true,"message":"Folder deleted successfully"}
}

export function teardown() {
  // Close the SMB connection
  client.close();
}
```

See the [examples](./examples) folder for a more detailed usage example.

## Contributing

Contributions are welcome! Please see the [CONTRIBUTING.md](CONTRIBUTING.md) file for guidelines on how to contribute to this project.

## License

This project is licensed under the terms of the [AGPL-3.0-only](LICENSE) license.

© 2025 INDUSTRIA DE DISEÑO TEXTIL S.A. (INDITEX S.A.)
