// SPDX-FileCopyrightText: 2025 INDUSTRIA DE DISEÑO TEXTIL S.A. (INDITEX S.A.)
//
// SPDX-License-Identifier: AGPL-3.0-only

package xk6smb

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/hirochachacha/go-smb2"
)

type Client struct {
}

type OperationResult struct {
	Success bool
	Message string
}

type GetSharesResult struct {
	Success bool
	Message string
	Shares  []string
}

type SmbClient struct {
	conn    net.Conn
	dialer  *smb2.Dialer
	session *smb2.Session
	share   *smb2.Share
}

func (*Client) NewClient(addressWithPort string, username string, psw string, shareName string) *SmbClient {
	var err error
	c := new(SmbClient)
	c.conn, err = initConn(addressWithPort)
	if err != nil {
		logger.Error(err)
		return nil
	}
	c.dialer = initDialer(username, psw)
	c.session, err = c.initSession()
	if err != nil {
		logger.Error(err)
		return nil
	}
	c.share, err = c.mount(shareName)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return c
}

func initConn(addressWithPort string) (net.Conn, error) {
	conn, err := net.Dial("tcp", addressWithPort)
	if err != nil {
		return nil, formatErr(err)
	}
	return conn, nil
}

func initDialer(user string, psw string) *smb2.Dialer {
	dialer := &smb2.Dialer{
		Initiator: &smb2.NTLMInitiator{
			User:     user,
			Password: psw,
		},
	}
	return dialer
}

func (c *SmbClient) initSession() (*smb2.Session, error) {
	if c == nil || c.dialer == nil {
		return nil, formatErr(errors.New("Client Dialer not initialized"))
	}
	s, err := c.dialer.Dial(c.conn)
	if err != nil {
		fmt.Println(err)
		return nil, formatErr(err)
	}
	return s, formatErr(err)
}

func (c *SmbClient) mount(shareName string) (*smb2.Share, error) {
	if c == nil || c.session == nil {
		return nil, formatErr(errors.New("Client not initialized"))
	}

	fs, err := c.session.Mount(shareName)
	if err != nil {
		return nil, formatErr(err)
	}
	return fs, nil
}

func (c *SmbClient) IsConnected() bool {
	if c != nil && c.conn != nil && c.dialer != nil && c.share != nil {
		result := c.GetShares()
		return result.Success
	}
	return false
}

func (c *SmbClient) Close() {
	if c != nil {
		if c.share != nil {
			defer c.share.Umount() //nolint:errcheck
		}
		if c.session != nil {
			defer c.session.Logoff() //nolint:errcheck
		}
		if c.conn != nil {
			defer c.conn.Close() //nolint:errcheck
		}
	}

}

func (c *SmbClient) GetShares() GetSharesResult {
	if c == nil || c.session == nil {
		return GetSharesResult{
			Success: false,
			Message: "Client not initialized",
			Shares:  nil,
		}
	}
	names, err := c.session.ListSharenames()
	if err != nil {
		return GetSharesResult{
			Success: false,
			Message: err.Error(),
			Shares:  nil,
		}
	}

	return GetSharesResult{
		Success: true,
		Message: "",
		Shares:  names,
	}
}

func (c *SmbClient) AppendLine(fileName string, strToWrite string) OperationResult {
	return c.AppendBytes(fileName, []byte(strToWrite), true)
}

func (c *SmbClient) AppendString(fileName string, strToWrite string) OperationResult {
	return c.AppendBytes(fileName, []byte(strToWrite), false)
}

func (c *SmbClient) AppendBytes(fileName string, bytes []byte, newLine bool) OperationResult {
	f, err := openOrCreate(c, fileName)
	if err != nil {
		fmt.Println(err)
		return OperationResult{Success: false, Message: err.Error()}
	}
	defer f.Close()
	stats, errStats := f.Stat()
	if errStats != nil {
		return OperationResult{Success: false, Message: errStats.Error()}
	}

	if stats.Size() > 0 && newLine {
		bytes = []byte("\n" + string(bytes))
	}
	_, err = f.WriteAt(bytes, stats.Size())
	if err != nil {
		fmt.Println(err)
		return OperationResult{Success: false, Message: err.Error()}
	}
	return OperationResult{Success: true, Message: "Bytes appended successfully"}
}

func (c *SmbClient) CopyFile(srcFileName string, destFileName string) OperationResult {
	srcFile, err := os.Open(srcFileName)
	if err != nil {
		return OperationResult{Success: false, Message: err.Error()}
	}
	defer srcFile.Close()

	destFile, err := c.share.Create(destFileName)
	if err != nil {
		return OperationResult{Success: false, Message: err.Error()}
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return OperationResult{Success: false, Message: err.Error()}
	}

	return OperationResult{Success: true, Message: "File copied successfully"}
}

func (c *SmbClient) ReadFile(fileName string) string {
	f, err := openFile(c, fileName)
	if err != nil {
		logger.Error(err)
		return ""
	}

	defer f.Close()
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		logger.Error(err)
		return ""
	}
	bs, err := io.ReadAll(f)
	if err != nil {
		logger.Error(err)
		return ""
	}

	return string(bs)
}

func (c *SmbClient) RemoveFile(fileName string) OperationResult {
	if c == nil || c.share == nil {
		return OperationResult{Success: false, Message: "Client not initialized"}
	}
	err := c.share.Remove(fileName)
	if err != nil {
		return OperationResult{Success: false, Message: err.Error()}
	}
	return OperationResult{Success: true, Message: "File removed successfully"}
}

func (c *SmbClient) RenameFile(pathOld string, pathNew string) OperationResult {
	if c == nil || c.share == nil {
		return OperationResult{Success: false, Message: "Client not initialized"}
	}
	err := c.share.Rename(pathOld, pathNew)
	if err != nil {
		return OperationResult{Success: false, Message: err.Error()}
	}
	return OperationResult{Success: true, Message: "File renamed successfully"}
}

func (c *SmbClient) FileExists(name string) bool {
	if c == nil || c.share == nil {
		logger.Error("Client not initialized")
		return false
	}
	_, err := c.share.Lstat(name)
	if err != nil {
		logger.Error(err)
		return false
	}
	return true
}

func (c *SmbClient) DeleteFile(name string) OperationResult {
	if c == nil || c.share == nil {
		return OperationResult{Success: false, Message: "Client not initialized"}
	}
	err := c.share.Remove(name)
	if err != nil {
		return OperationResult{Success: false, Message: err.Error()}
	}
	return OperationResult{Success: true, Message: "File deleted successfully"}
}

func (c *SmbClient) CreateDir(name string) OperationResult {
	if c == nil || c.share == nil {
		return OperationResult{Success: false, Message: "Client not initialized"}
	}
	err := c.share.Mkdir(name, os.ModeDir)
	if err != nil {
		return OperationResult{Success: false, Message: err.Error()}
	}
	return OperationResult{Success: true, Message: "Folder created successfully"}
}

func (c *SmbClient) RenameDir(oldPath string, newPath string) OperationResult {
	if c == nil || c.share == nil {
		return OperationResult{Success: false, Message: "Client not initialized"}
	}
	err := c.share.Rename(oldPath, newPath)
	if err != nil {
		return OperationResult{Success: false, Message: err.Error()}
	}
	return OperationResult{Success: true, Message: "Folder renamed successfully"}
}

func (c *SmbClient) DirExists(name string) bool {
	if c == nil || c.share == nil {
		logger.Error("Client not initialized")
		return false
	}
	_, err := c.share.ReadDir(name)
	if err != nil {
		logger.Error(err)
		return false
	}
	return true
}

func (c *SmbClient) DeleteDir(name string) OperationResult {
	if c == nil || c.share == nil {
		return OperationResult{Success: false, Message: "Client not initialized"}
	}
	fInfo, err := c.share.ReadDir(name)
	if err != nil {
		return OperationResult{Success: false, Message: err.Error()}
	}
	if len(fInfo) > 0 {
		return OperationResult{Success: false, Message: "directory is not empty"}
	}
	err = c.share.Remove(name)
	if err != nil {
		return OperationResult{Success: false, Message: err.Error()}
	}
	return OperationResult{Success: true, Message: "Folder deleted successfully"}
}

func (c *SmbClient) IsDir(name string) bool {
	if c == nil || c.share == nil {
		logger.Error("Client not initialized")
		return false
	}
	fInfo, err := c.share.Stat(name)
	if err != nil {
		logger.Error(err)
		return false
	}
	return fInfo.IsDir()
}

func (c *SmbClient) ListFilesInDir(name string) []string {
	if c == nil || c.share == nil {
		logger.Error("Client not initialized")
		return nil
	}
	fInfo, err := c.share.ReadDir(name)
	if err != nil {
		logger.Error(err)
		return nil
	}
	if len(fInfo) == 0 {
		logger.Error("No files found in folder")
		return nil
	}
	var retArray []string
	for _, info := range fInfo {
		retArray = append(retArray, info.Name())
	}
	return retArray
}

func openOrCreate(c *SmbClient, fileName string) (*smb2.File, error) {
	if c == nil || c.share == nil {
		return nil, formatErr(errors.New("Client not initialized"))
	}
	f, err := c.share.OpenFile(fileName, os.O_APPEND, os.ModeAppend)
	if err != nil {
		f, err = c.share.Create(fileName)
	}
	return f, formatErr(err)
}

func openFile(c *SmbClient, fileName string) (*smb2.File, error) {
	if c == nil || c.share == nil {
		return nil, formatErr(errors.New("Client not initialized"))
	}
	f, err := c.share.OpenFile(fileName, os.O_APPEND, os.ModeAppend)
	if err != nil {
		return nil, formatErr(err)
	}
	return f, err
}

// formatErr add #ERROR# prefix to identify error in retur string
func formatErr(err error) error {
	if err != nil && (!strings.Contains(err.Error(), "#ERROR#")) {
		return errors.New("#ERROR# " + err.Error())
	}
	return err
}
