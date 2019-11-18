package main

import (
	"bufio"
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type SCENELIST struct {
	XMLName xml.Name `xml:"SCENELIST"`
	Scene []Tags `xml:"SCENE"`
	Region SCENE `xml:"REGION"`
}

type SCENE struct {
	Scene []Tags `xml:"SCENE"`
}

type Tags struct {
	ID string `xml:"ID,attr"`
	Desc string `xml:"desc,attr"`
	Config string `xml:"config,attr"`
	Value  string `xml:",chardata"`
}

func main() {

	xmlPath := flag.String("d", "notset", "Path to XMLs")
	outTxt := flag.String("o", "notset", "Path to output txt")
	flag.Parse()

	if *xmlPath == "notset" {
		fmt.Println("XMLs path not set")
		os.Exit(1)
	}
	if *outTxt == "notset" {
		fmt.Println("Output path not set")
		os.Exit(1)
	}

	parseDir(*xmlPath, *outTxt)
}

func parseDir(dir string, outTxt string){
	var scenes []string
	fileList := []string{}
	_ = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})

	for _, file := range fileList {
		isXML := strings.HasSuffix(strings.ToLower(file), ".xml")
		if isXML {
			fmt.Println(file)
			sceneList := new(SCENELIST)

			data, _ := ioutil.ReadFile(file)
			_ = xml.Unmarshal([]byte(data), &sceneList)

			for _, s := range sceneList.Scene{
				fmt.Println(s.Config)
				scenes = append(scenes, s.Config + "\n")
			}

			for _, s := range sceneList.Region.Scene{
				fmt.Println(s.Config)
				scenes = append(scenes, s.Config + "\n")
			}

			//if sceneList.Region.Scene[file].Config != ""{
			//	fmt.Println(sceneList.Scene.Config)
			//	scenes = append(scenes, sceneList.Region.Scene.Config + "\n")
			//}
			//if sceneList.Scene.Config != "" {
			//	fmt.Println(sceneList.Scene.Config)
			//	scenes = append(scenes, sceneList.Scene.Config + "\n")
			//}
		}
	}

	err := WriteToFile(outTxt, scenes)
	if err != nil {
		log.Fatal(err)
	}
}

func WriteToFile(filename string, data []string) error {
	hasTxt := strings.HasSuffix(filename, ".txt")
	if !hasTxt {
		filename = filename + ".txt"
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	for _, line := range data {
		fmt.Fprint(w, line)
	}
	return w.Flush()
}