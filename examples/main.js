import xk6smb from "k6/x/smb";
import { group, check, sleep } from "k6";

const client = xk6smb.newClient("localhost:445", "user", "pwd", "User Volume");

export default function () {
    const fileName = `test_${__ITER}_${__VU}.txt`;
    const dirName = `test-dir_${__ITER}_${__VU}`;

    group("Get Shares", function () {
        const shares = client.getShares();
        check(shares, {
            "shares retrieved": (s) => s.success === true,
        });
    });

    group("Write File", function () {
        const writeResult = client.appendString(fileName, "Hello, World!");
        check(writeResult, {
            "file written": (w) => w.success === true,
        });
    });

    group("File exists", function () {
        const isFile = client.fileExists(fileName);
        check(isFile, {
            "is File": (d) => d === true,
        });
    });

    group("Read File", function () {
        const readResult = client.readFile(fileName);
        check(readResult, {
            "file read": (r) => r !== null,
        });
    });

    group("Delete File", function () {
        const deleteResult = client.deleteFile(fileName);
        check(deleteResult, {
            "file deleted": (d) => d.success === true,
        });
    });

    group("Create Dir", function () {
        const createResult = client.createDir(dirName);
        check(createResult, {
            "dir created": (d) => d.success === true,
        });
    });

    client.appendString(`${dirName}/${fileName}`, "Hello, World!");

    group("List Dir", function () {
        const listResult = client.listFilesInDir(dirName);
        check(listResult, {
            "dir listed": (d) => d.length > 0,
        });
    });

    group("Is Dir?", function () {
        const isDir = client.isDir(dirName);
        check(isDir, {
            "is Dir": (d) => d === true,
        });
    });

    group("Dir exists", function () {
        const isDir = client.dirExists(dirName);
        check(isDir, {
            "is Dir": (d) => d === true,
        });
    });

    group("Delete file in dir", function () {
        const deleteResult = client.deleteFile(`${dirName}/${fileName}`);
        check(deleteResult, {
            "file deleted": (d) => d.success === true,
        });
    });

    group("Delete dir", function () {
        const deleteResult = client.deleteDir(dirName);
        check(deleteResult, {
            "dir deleted": (d) => d.success === true,
        });
    });

    group("Copy file", function () {
        const copyResult = client.copyFile("./examples/test-data/test-file.txt", "copy.txt");
        console.log(copyResult);
        check(copyResult, {
            "file copied": (d) => d.success === true,
        });
    });

    sleep(0.3);
}

export function teardown() {
    client.close();
}