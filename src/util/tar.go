package util

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

// TarGz 压缩
func TarGz(src string, dstTgz string, ow bool) error {
	// 清理路径字符串
	src = path.Clean(src)

	// 判断目标文件是否存在
	if FileExist(dstTgz) {
		if ow { // overwrite 已存在的文件
			if err := os.Remove(dstTgz); err != nil {
				// log.Printf("ERROR: %s\n", err.Error())
				return err
			}
		} else {
			return fmt.Errorf(ErrFileExisted, dstTgz)
		}
	}

	// 创建目标文件
	fw, err := os.Create(dstTgz)
	if err != nil {
		// log.Printf("ERROR: %s\n", err.Error())
		return err
	}
	defer fw.Close()

	gw := gzip.NewWriter(fw)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// 打包
	fi, err := os.Stat(src)
	if err != nil {
		// log.Printf("ERROR: %s\n", err.Error())
		return err
	}

	srcBase, srcRelative := path.Split(src)

	if fi.IsDir() {
		// tarGzDir(srcBase, srcRelative, tw, fi)
		fis, err := ioutil.ReadDir(src)
		if err != nil {
			// log.Printf("ERROR: %s\n", err.Error())
			return err
		}

		// 遍历 文件夹
		for _, fi := range fis {
			if fi.IsDir() {
				tarGzDir(src, fi.Name(), tw, fi)
			} else {
				tarGzFile(src, fi.Name(), tw, fi)
			}
		}
	} else {
		tarGzFile(srcBase, srcRelative, tw, fi)
	}

	return nil
}

// UnTarGz 解压
func UnTarGz(srcTarGz string, dstDir string) error {
	// 清理字符串路径
	dstDir = path.Clean(dstDir) + string(os.PathSeparator)

	// 打开要解压的包
	fr, err := os.Open(srcTarGz)
	if err != nil {
		// log.Printf("ERROR: %s\n", err.Error())
		return err
	}
	defer fr.Close()

	// 创建 Reader 来读取包中的内容
	gr, err := gzip.NewReader(fr)
	if err != nil {
		// log.Printf("ERROR: %s\n", err.Error())
		return err
	}
	defer gr.Close()

	tr := tar.NewReader(gr)

	// 遍历包中的文件
	for hdr, err := tr.Next(); err != io.EOF; hdr, err = tr.Next() {
		if err != nil {
			// log.Printf("ERROR: %s\n", err.Error())
			return err
		}

		// 文件头信息
		fi := hdr.FileInfo()
		// 解压路径
		dstFullPath := dstDir + hdr.Name

		if hdr.Typeflag == tar.TypeDir {
			// 创建目录
			os.MkdirAll(dstFullPath, fi.Mode().Perm())
		} else {
			// 创建文件所在的目录

			os.MkdirAll(path.Dir(dstFullPath), os.ModePerm)

			// 将 tr 中的数据写入文件中
			if err := unTarGzFile(dstFullPath, tr); err != nil {
				// log.Printf("ERROR: %s\n", err.Error())
				return err
			}
			os.Chmod(dstFullPath, fi.Mode().Perm())
		}
	}

	return nil
}

// 压缩文件夹
func tarGzDir(srcBase, srcRelative string, tw *tar.Writer, fi os.FileInfo) error {

	src := path.Join(srcBase, srcRelative)
	// 末尾添加 "/"
	last := len(srcRelative) - 1
	if srcRelative[last] != os.PathSeparator {
		srcRelative += string(os.PathSeparator)
	}

	fis, err := ioutil.ReadDir(src)
	if err != nil {
		// log.Printf("ERROR: %s\n", err.Error())
		return err
	}

	// 遍历 文件夹
	for _, fi := range fis {
		if fi.IsDir() {
			tarGzDir(srcBase, srcRelative+fi.Name(), tw, fi)
		} else {
			tarGzFile(srcBase, srcRelative+fi.Name(), tw, fi)
		}
	}

	// 当前目录信息写入 targz 文件中
	if len(srcRelative) > 0 {
		hdr, err := tar.FileInfoHeader(fi, "")
		if err != nil {
			// log.Printf("ERROR: %s\n", err.Error())
			return err
		}
		hdr.Name = srcRelative

		if err = tw.WriteHeader(hdr); err != nil {
			// log.Printf("ERROR: %s\n", err.Error())
			return err
		}
	}

	return nil
}

// 压缩文件
func tarGzFile(srcBase, srcRelative string, tw *tar.Writer, fi os.FileInfo) error {
	// 文件头信息写入 targz 中
	hdr, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		// log.Printf("ERROR: %s\n", err.Error())
		return err
	}
	hdr.Name = srcRelative

	if err = tw.WriteHeader(hdr); err != nil {
		// log.Printf("ERROR: %s\n", err.Error())
		return err
	}

	// 文件内容写入 targz 中
	fr, err := os.Open(path.Join(srcBase, srcRelative))
	if err != nil {
		// log.Printf("ERROR: %s\n", err.Error())
		return err
	}
	defer fr.Close()

	// 将文件数据写入 tw 中
	if _, err = io.Copy(tw, fr); err != nil {
		// log.Printf("ERROR: %s\n", err.Error())
		return err
	}

	return nil
}

// 解压某个文件
func unTarGzFile(dstFile string, tr *tar.Reader) error {
	// 创建空文件, 准备写入解压后的数据
	fw, err := os.Create(dstFile)
	if err != nil {
		// log.Printf("ERROR: %s\n", err.Error())
		return err
	}
	defer fw.Close()

	// 空文件中写入数据
	_, err = io.Copy(fw, tr)
	if err != nil {
		// log.Printf("ERROR: %s\n", err.Error())
		return err
	}

	return nil
}
