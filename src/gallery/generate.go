package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//var carpetaDeFotos = "Visita-Luchetti"
var carpetasDeFotos = []string{
	"Visita-Luchetti",
	"Talleres-Matlab",
	}

const (
	relProjectPath = "../../"
	subsiteDir  = "/gallery/"
	photoDir    = "/gallery/photos/"
)
type photo struct {
	name string
	UiName string
	ProjectPath string
}
type album struct {
	Title      string
	Photos []photo
	FileProjectPath string
	relPath string
}

func main() {
	for _, carpeta := range carpetasDeFotos {
		album,err := genAlbum(filepath.Join(relProjectPath,photoDir,carpeta))
		tplAlbum, err := template.New("gallery").Parse(file2str("tplAlbum.html"))
		if err != nil {
			panic(err)
		}
		f, err := os.Create(filepath.Join(relProjectPath,subsiteDir,album.getFilename()))
		if err != nil {
			panic(err)
		}
		err = tplAlbum.Execute(f,album)
		if err != nil {
			panic(err)
		}
	}
}

const (
	bufsize = 1 << 8
)


func file2str(filename string) (content string) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	buf := make([]byte, bufsize)
	eof := false
	for !eof {
		n, err := f.Read(buf)
		content += string(buf[:n])
		if err != nil {
			eof = true
		}
	}
	return content
}

func (a *album)getFilename() string {
	aux := strings.ReplaceAll(a.Title," - ","-")
	return strings.ReplaceAll(strings.ToLower(aux)," ","-") + ".html"
}
func genAlbum(albumPath string) (*album, error) {
	var theAlbum album
	theAlbum.relPath = albumPath
	finfos, err := ioutil.ReadDir(albumPath)
	if err != nil {
		return &theAlbum, err
	}
	for _, f := range finfos {
		var p photo
		p.name = f.Name()
		ext := strings.ToLower(p.name[len(p.name)-3:])
		if ext != "jpg" && ext != "png" {
			continue
		}
		relativePath := filepath.Join(albumPath, p.name)
		p.ProjectPath =  strings.ReplaceAll(string(filepath.Separator) +localToProjectPath(relativePath),string(filepath.Separator),"/")
		p.setUiName()
		theAlbum.Photos = append(theAlbum.Photos,p)
	}
	if len(theAlbum.Photos) == 0 {
		return &theAlbum, fmt.Errorf("no album pictures found")
	}
	theAlbum.setTitle()
	theAlbum.FileProjectPath = subsiteDir + theAlbum.getFilename()
	return &theAlbum, nil
}

func (p *photo)setUiName() {
	p.UiName = strings.ReplaceAll(p.name[:strings.LastIndex(p.name,".")], "_"," ")
}

func (a *album)setTitle() {
	dir := filepath.Dir(a.Photos[0].ProjectPath)
	foldername := dir[1+strings.LastIndex(dir[:len(dir)-2],string(filepath.Separator)):]
	splitHyph := strings.Split(foldername,"-")
	titleName := strings.Join(splitHyph," - ")
	a.Title = strings.ReplaceAll(titleName,"_"," ")
}

func localToProjectPath(path string) string {
	prjpath,err:= filepath.Rel(relProjectPath, path)
	if err != nil {
		panic(err)
	}
	return prjpath
}