package misc

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"time"
)

///
type logWriter struct{}

func SetLogTimeFmt() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
}

func (writer *logWriter) Write(bts []byte) (int, error) {
	return fmt.Print(time.Now().Format(time.RFC3339) + " " + string(bts))
}

func FileCopy(src, dst string) (err error) {
	var in, out *os.File

	if in, err = os.Open(src); err != nil {
		return err
	}
	defer in.Close()

	if out, err = os.Create(dst); err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return nil
}

func CmdMd5(cmd []string) string {
	bts, _ := json.Marshal(map[string][]string{"commandline": cmd})

	hash := md5.Sum(bts)
	return hex.EncodeToString(hash[:])
}

func FileExists(p string) (yes bool, err error) {
	var info fs.FileInfo

	if info, err = os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return !info.IsDir(), nil
}

func DirExists(p string) (yes bool, err error) {
	var info fs.FileInfo

	if info, err = os.Stat(p); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return info.IsDir(), nil
}
