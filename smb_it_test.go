// SPDX-FileCopyrightText: 2025 INDUSTRIA DE DISEÑO TEXTIL S.A. (INDITEX S.A.)
//
// SPDX-License-Identifier: AGPL-3.0-only

package xk6smb

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const (
	username  = "user"
	password  = "pwd"
	host      = "localhost:445"
	shareName = "User Volume"
)

func TestConnect(t *testing.T) {
	c := new(Client)
	client := c.NewClient(host, username, password, shareName)
	if client == nil {
		t.Fail()
	}
	client.Close()
}

func TestGetShares(t *testing.T) {
	c := new(Client)
	client := c.NewClient(host, username, password, shareName)
	if client == nil {
		t.Fail()
	}
	defer client.Close()

	shares := client.GetShares()
	if !shares.Success {
		fmt.Println(shares.Message)
		t.Fail()
	}
	fmt.Println("Shares:", shares.Message)
}

func TestFileOperations(t *testing.T) {
	c := new(Client)
	client := c.NewClient(host, username, password, shareName)
	require.NotNil(t, client)
	defer client.Close()

	// Test file creation and existence check
	randomName := fmt.Sprintf("testfile_%d.txt", time.Now().UnixNano())
	createResult := client.AppendString(randomName, "Hello, World!")
	require.True(t, createResult.Success, createResult.Message)

	existsResult := client.FileExists(randomName)
	require.True(t, existsResult)

	// Test file reading
	content := client.ReadFile(randomName)
	require.Equal(t, "Hello, World!", content)

	// Test file deletion
	deleteResult := client.DeleteFile(randomName)
	require.True(t, deleteResult.Success, deleteResult.Message)

	existsResult = client.FileExists(randomName)
	require.False(t, existsResult)

	// Test creating a new text file in /tmp folder
	tmpFileName := "/tmp/testfile.txt"
	fileContent := []byte("Temporary file content")
	err := os.WriteFile(tmpFileName, fileContent, 0644)
	require.NoError(t, err)

	// Copy the temporary file to the SMB share
	copyResult := client.CopyFile(tmpFileName, randomName)
	require.True(t, copyResult.Success, copyResult.Message)

	// Verify the file exists on the SMB share
	existsResult = client.FileExists(randomName)
	require.True(t, existsResult)

	// Clean up the copied file from the SMB share
	deleteResult = client.DeleteFile(randomName)
	require.True(t, deleteResult.Success, deleteResult.Message)

	// Clean up the temporary file from the local filesystem
	err = os.Remove(tmpFileName)
	require.NoError(t, err)

}

func TestFolderOperations(t *testing.T) {
	c := new(Client)
	client := c.NewClient(host, username, password, shareName)
	require.NotNil(t, client)
	defer client.Close()

	// Test folder creation and existence check
	folderName := fmt.Sprintf("testfolder_%d", time.Now().UnixNano())
	createResult := client.CreateDir(folderName)
	require.True(t, createResult.Success, createResult.Message)

	existsResult := client.DirExists(folderName)
	require.True(t, existsResult)

	// Test folder renaming
	newFolderName := fmt.Sprintf("renamedfolder_%d", time.Now().UnixNano())
	renameResult := client.RenameDir(folderName, newFolderName)
	require.True(t, renameResult.Success, renameResult.Message)

	existsResult = client.DirExists(newFolderName)
	require.True(t, existsResult)

	// Test folder deletion
	deleteResult := client.DeleteDir(newFolderName)
	require.True(t, deleteResult.Success, deleteResult.Message)

	existsResult = client.DirExists(newFolderName)
	require.False(t, existsResult)
}
