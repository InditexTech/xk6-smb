import xk6smb from "k6/x/smb";
import { group, check, sleep } from "k6";

export const options = {
    iterations: 1,
    vus: 1,
};

const client = xk6smb.newClient("localhost:445", "user", "pwd", "User Volume");

export default function () {
    const fileName = `test_${__ITER}_${__VU}.txt`;
    const dirName = `test-dir_${__ITER}_${__VU}`;

    group("Get Shares", function () {
        const shares = client.getShares();
        console.log("getShares result:", shares);
        check(shares, {
            "shares retrieved": (s) => s.success === true,
        });
    });

    group("Write File", function () {
        const writeResult = client.appendString(fileName, "Hello, World!");
        console.log("appendString result:", writeResult);
        check(writeResult, {
            "file written": (w) => w.success === true,
        });
    });

    group("File exists", function () {
        const isFile = client.fileExists(fileName);
        console.log("fileExists result:", isFile);
        check(isFile, {
            "is File": (d) => d === true,
        });
    });

    group("Read File", function () {
        const readResult = client.readFile(fileName);
        console.log("readFile result:", readResult);
        check(readResult, {
            "file read": (r) => r !== null,
        });
    });

    group("Delete File", function () {
        const deleteResult = client.deleteFile(fileName);
        console.log("deleteFile result:", deleteResult);
        check(deleteResult, {
            "file deleted": (d) => d.success === true,
        });
    });

    group("Create Dir", function () {
        const createResult = client.createDir(dirName);
        console.log("createDir result:", createResult);
        check(createResult, {
            "dir created": (d) => d.success === true,
        });
    });

    group("List Dir", function () {
        const listResult = client.listFilesInDir(dirName);
        console.log("listFilesInDir result:", listResult);
        check(listResult, {
            "dir listed": (d) => d.length > 0,
        });
    });

    group("Is Dir?", function () {
        const isDir = client.isDir(dirName);
        console.log("isDir result:", isDir);
        check(isDir, {
            "is Dir": (d) => d === true,
        });
    });

    group("Dir exists?", function () {
        const dirExists = client.dirExists(dirName);
        console.log("dirExists result:", dirExists);
        check(dirExists, {
            "is Dir": (d) => d === true,
        });
    });

    group("Delete file in dir", function () {
        const deleteResult = client.deleteFile(`${dirName}/${fileName}`);
        console.log("deleteFile result:", deleteResult);
        check(deleteResult, {
            "file deleted": (d) => d.success === true,
        });
    });

    group("Delete dir", function () {
        const deleteResult = client.deleteDir(dirName);
        console.log("deleteDir result:", deleteResult);
        check(deleteResult, {
            "dir deleted": (d) => d.success === true,
        });
    });

    group("Copy file", function () {
        const copyResult = client.copyFile("./examples/test-data/test-file.txt", "copy.txt");
        console.log("copyFile result:", copyResult);
        check(copyResult, {
            "file copied": (d) => d.success === true,
        });
    });

    sleep(0.3);
}

export function teardown() {
    client.close();
}