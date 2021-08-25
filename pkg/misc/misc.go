package misc

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
)

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
