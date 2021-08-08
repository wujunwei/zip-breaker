package main

import (
	"flag"
	"fmt"
	"github.com/yeka/zip"
	"io"
	"os"
)

var Length int
var PasswordSet int
var FilePath string
var password string

const (
	BigAlp = 1 << iota
	SmallAlp
	Num
	Sign
)

func init() {
	flag.IntVar(&Length, "l", 6, "the length of the zip password.")
	flag.IntVar(&PasswordSet, "s", 2, "the source set of the password to guess.")
	flag.StringVar(&FilePath, "f", "/Users/adam/Downloads/test.zip", "the path to the zip file to get out")

}

func main() {
	var source []byte
	if PasswordSet&BigAlp != 0 {
		source = append(source, "ABCDEFGHIGKLMNOPQRSTUVWXYZ"...)
	}
	if PasswordSet&SmallAlp != 0 {
		source = append(source, "abcdefghigklmnopqrstuvwxyz"...)
	}
	if PasswordSet&Num != 0 {
		source = append(source, "1234567890"...)
	}
	if PasswordSet&Sign != 0 {
		source = append(source, ".,/"...)
	}
	fmt.Println("use passwordSet:" + string(source))
	r, err := zip.OpenReader(FilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, f := range r.File {
		fmt.Println(" start to unzip :" + f.Name)
		if f.IsEncrypted() {
			guess(source, func(s string) bool {
				fmt.Println(s)
				f.SetPassword(s)
				read, _ := f.Open()
				_, err = io.ReadAll(read)
				if err != nil {
					return false
				}
				return true
			})
			fmt.Println(password)
			os.Exit(0)
		}
	}
}

func guess(scr []byte, judge func(string) bool) {
	if password != "" {
		return
	}
	var dfs func(try []byte, selected map[int]bool) bool
	dfs = func(try []byte, selected map[int]bool) bool {
		if len(try) == Length {
			if judge(string(try)) {
				password = string(try)
				return true
			}
			return false
		}
		for i := 0; i < len(scr); i++ {
			if selected[i] {
				continue
			}
			selected[i] = true
			if dfs(append(try, scr[i]), selected) {
				return true
			}
			selected[i] = false
		}
		return false
	}
	dfs([]byte{}, map[int]bool{})
}
