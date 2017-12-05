package main

import (
	"fmt"
	"path/filepath"
	"os"
	"net/http"
	"log"
	"io/ioutil"
	"strconv"
)

var rootDir = "."

func main() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(getListenAddr(), nil))
}

func handler(rsp http.ResponseWriter, req *http.Request) {
	file := rootDir + filepath.FromSlash(req.RequestURI)
	fileInfo, err := os.Stat(file)
	if err != nil {
		if os.IsNotExist(err) || os.IsPermission(err) {
			rsp.WriteHeader(404)
			rsp.Write([]byte("<h1>File Not Found</h1>"))
		} else {
			rsp.WriteHeader(500)
			rsp.Write([]byte("<h1>Server Error</h1>"))
			fmt.Println(err.Error())
		}
		return
	}
	if !fileInfo.IsDir() {
		download(rsp, file)
		return
	} else {
		list(rsp, file)
	}
}

func download(rsp http.ResponseWriter, filename string) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		rsp.WriteHeader(500)
		rsp.Write([]byte("<h1>Server Error</h1>"))
		return;
	}
	rsp.Write(data)
}

func list(rsp http.ResponseWriter, dirName string)  {
	fileInfoList, err := ioutil.ReadDir(dirName)
	if err != nil {
		rsp.WriteHeader(500)
		rsp.Write([]byte("<h1>Server Error</h1>"))
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