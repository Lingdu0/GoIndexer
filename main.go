package main

import (
	"fmt"
	"path/filepath"
	"os"
	"net/http"
	"log"
	"io/ioutil"
	"strconv"
	"net/url"
)

var rootDir = "."

func main() {
	fmt.Printf("It works on %v\r\n", getListenAddr())
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(getListenAddr(), nil))
}

func handler(rsp http.ResponseWriter, req *http.Request) {
	fmt.Printf("%v %v\r\n", req.Method, req.URL)
	uri, err := url.QueryUnescape(req.RequestURI)
	if err != nil {
		http.NotFound(rsp, req)
	}
	file := rootDir + filepath.FromSlash(uri)
	fileInfo, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) || os.IsPermission(err) {
			http.NotFound(rsp, req)
		} else {
			http.Error(rsp, "internal error", 500)
			fmt.Println(err.Error())
		}
		return
	}
	if !fileInfo.IsDir() {
		f, err := os.Open(file)
		if err != nil {
			http.NotFound(rsp, req)
		}
		http.ServeContent(rsp, req, file, fileInfo.ModTime(), f)
		return
	} else {
		list(rsp, file)
		return
	}
}

func list(rsp http.ResponseWriter, dirName string)  {
	fileInfoList, err := ioutil.ReadDir(dirName)
	if err != nil {
		http.Error(rsp, "internal error", 500)
		fmt.Println(err.Error())
		return
	}
	rsp.Write(parseDirListHtml(fileInfoList))
}

func parseDirListHtml(fileInfoList []os.FileInfo) []byte {
	html := `<html>
<head><title>GoIndexer</title></head>
<body bgcolor="white">
<h1>GoIndexer</h1><hr><pre>
<a href="../">../</a>
`
	fileLink := "<a href=\"%v\">%v</a>                               %v             %v"

	for _, fileInfo := range fileInfoList {
		fileSize := strconv.FormatInt(fileInfo.Size(), 10)
		fileName := fileInfo.Name()
		fileTime := fileInfo.ModTime().Format("2006-01-02 15:04:05")
		if fileInfo.IsDir() {
			fileSize = "-"
			fileName += "/"
		}
		html += "\r\n" + fmt.Sprintf(fileLink, fileName, fileName, fileTime, fileSize)
	}

	html += `</pre><hr></body>
</html>
`
	return []byte(html)
}

func getListenAddr() string {
	address := ":8000"
	if len(os.Args) == 2 {
		address = os.Args[1]
	}
	return address
}