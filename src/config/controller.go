package config

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
	"github.com/nwaples/rardecode/v2"
)

// 接口,检查压缩包
type ZipHandle interface {
	Check() bool
	Close() error
	Extract(dst string) error
}

// zip ------------------------------------------------------------------
type ZipChecker struct {
	reader *zip.ReadCloser
}

func NewZipChecker(path string) (*ZipChecker, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	return &ZipChecker{reader: r}, nil
}

func (z *ZipChecker) Close() error { return z.reader.Close() }

func (z *ZipChecker) Check() bool {
	for _, f := range z.reader.File {
		// 1. 检查加密 (Flags 第 0 位)
		if f.Flags&1 == 1 {
			continue // 如果某个文件加密，跳过或返回 false（取决于你的业务）
		}
		// 2. 匹配后缀
		if strings.HasSuffix(f.Name, ".model3.json") {
			rc, err := f.Open()
			if err != nil {
				return false // 无法打开，可能是损坏或伪装加密
			}
			rc.Close()
			return true
		}
	}
	return false
}

func (z *ZipChecker) Extract(dst string) error {
	for _, f := range z.reader.File {

		// 1. 跳过目录
		if f.FileInfo().IsDir() {
			continue
		}

		// 2. 打开文件
		rc, err := f.Open()
		if err != nil {
			return err
		}

		// 3. 创建目标文件路径
		outPath := filepath.Join(dst, f.Name)

		// 防止 zip slip（重要！）
		if !strings.HasPrefix(outPath, filepath.Clean(dst)+string(os.PathSeparator)) {
			rc.Close()
			return fmt.Errorf("illegal file path: %s", f.Name)
		}

		// 4. 创建目录
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			rc.Close()
			return err
		}

		// 5. 写文件
		outFile, err := os.Create(outPath)
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// 7z --------------------------------------------------------------------------
type SevenZipChecker struct {
	reader *sevenzip.ReadCloser
}

func NewSevenZipChecker(path string) (*SevenZipChecker, error) {
	r, err := sevenzip.OpenReader(path)
	if err != nil {
		return nil, err // 如果7z整体加密了文件名，这里通常就会直接报错
	}
	return &SevenZipChecker{reader: r}, nil
}

func (s *SevenZipChecker) Close() error { return s.reader.Close() }

func (s *SevenZipChecker) Check() bool {
	for _, f := range s.reader.File {
		if !strings.HasSuffix(f.Name, ".model3.json") {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			continue
		}

		buf := make([]byte, 1)
		_, err = rc.Read(buf)
		rc.Close()

		if err != nil && err != io.EOF {
			continue
		}

		return true
	}
	return false
}

func (s *SevenZipChecker) Extract(dst string) error {
	for _, f := range s.reader.File {

		if f.FileInfo().IsDir() {
			continue
		}

		rc, err := f.Open()
		if err != nil {
			continue
		}

		outPath := filepath.Join(dst, f.Name)

		// 防 zip-slip
		if !strings.HasPrefix(outPath, filepath.Clean(dst)+string(os.PathSeparator)) {
			rc.Close()
			return fmt.Errorf("illegal path: %s", f.Name)
		}

		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			rc.Close()
			return err
		}

		outFile, err := os.Create(outPath)
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// rar ----------------------------------------------------------------------------
type RarChecker struct {
	path string
}

func NewRarChecker(path string) (*RarChecker, error) {
	// RAR 采用流式读取，这里先只存路径
	return &RarChecker{path: path}, nil
}

func (r *RarChecker) Close() error { return nil }

func (r *RarChecker) Check() bool {
	rr, err := rardecode.OpenReader(r.path)
	if err != nil {
		return false
	}
	defer rr.Close()

	for {
		header, err := rr.Next()
		if err != nil {
			break
		}

		if !strings.HasSuffix(header.Name, ".model3.json") {
			continue
		}

		buf := make([]byte, 1)
		_, err = rr.Read(buf)
		if err != nil && err != io.EOF {
			continue
		}

		return true
	}

	return false
}

func (r *RarChecker) Extract(dst string) error {
	rr, err := rardecode.OpenReader(r.path)
	if err != nil {
		return err
	}
	defer rr.Close()

	for {
		header, err := rr.Next()
		if err != nil {
			break
		}

		if header.IsDir {
			continue
		}

		outPath := filepath.Join(dst, header.Name)

		// 防 zip-slip
		if !strings.HasPrefix(outPath, filepath.Clean(dst)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal path: %s", header.Name)
		}

		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			return err
		}

		outFile, err := os.Create(outPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, rr)

		outFile.Close()

		if err != nil && err != io.EOF {
			return err
		}
	}

	return nil
}
