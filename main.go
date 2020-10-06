package main

import (
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	flag.Usage = func() {
		fmt.Printf(`Usage of md-embed:
	md-embed -o [Output filename] <Input filename>
	-o string
		Output file name (default "out.md")
`,
		)
	}

	var nameOut = flag.String("o", "out.md", "Output file name")
	flag.Parse()
	if len(flag.Args()) == 0 {
		panic("请指定文件名")
	}

	err := embedMarkdown(flag.Arg(0), *nameOut)
	if err != nil {
		panic(err)
	}

}

var dataUrls = make(map[string]string)

//翻译  ![图片alt](图片地址 "图片title") / ![图片alt](图片地址)
func embedMarkdown(fn, outName string) (err error) {
	r := regexp.MustCompile(`\!\[([^\)]*)\]\(([^ \)]+)([^\)]*)\)`)
	fp, err := os.Open(fn)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer fp.Close()
	fout, err := os.Create(outName)
	if err != nil {
		log.Println(err.Error())
		return
	}
	defer fout.Close()

	reader := bufio.NewReader(fp)
	for {
		var l1 []byte
		l1, _, err = reader.ReadLine()
		if err != nil {
			err = nil
			break
		}
		line1 := strings.TrimSpace(string(l1))
		m := r.FindAllStringSubmatch(line1, -1)
		if m == nil {
			fout.WriteString(line1 + "\n")
			continue
		}
		for i := range m {
			fmt.Fprintf(os.Stderr, "\n%d: %#v\n", i+1, m[i])
			datName := genDataUrl(m[i][2])
			fmt.Fprintf(fout, "![%s][%s]", m[i][1], datName)
			dataUrls[datName] = genDataSection(getImagePath(fn, m[i][2]))
		}
	}
	for k := range dataUrls {
		fmt.Fprintf(fout, "\n[%s]:data:%s\n", k, dataUrls[k])
	}
	return
}

func getImagePath(mdFn, imgPath string) string {
	d := filepath.Dir(mdFn)
	return filepath.Join(d, imgPath)
}

func genDataUrl(p string) string {
	r := md5.Sum([]byte(p))
	return fmt.Sprintf("dat_%x", r)
}

func genDataSection(fn string) string {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return ""
	}
	encData := base64.StdEncoding.EncodeToString(data)
	typ := mime.TypeByExtension(path.Ext(fn))
	return typ + ";base64," + encData
}
