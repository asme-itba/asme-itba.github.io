package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type photo struct {
	name        string
	UiName      string
	ProjectPath string
}
type album struct {
	Title           string
	Description     string
	Photos          []photo
	FileProjectPath string
	relPath         string
	Folder          string
}

func (a *album) getFilename() string {
	aux := strings.ReplaceAll(a.Folder, " - ", "-")
	return strings.ReplaceAll(strings.ToLower(aux), " ", "-") + ".html"
}

func (a *album) completeAlbum(albumPath, relativeProjectPath, galleryPath string) error {
	a.relPath = albumPath
	finfos, err := ioutil.ReadDir(albumPath)
	if err != nil {
		return err
	}
	for _, f := range finfos {
		var p photo
		p.name = f.Name()
		ext := strings.ToLower(p.name[len(p.name)-3:])
		if ext != "jpg" && ext != "png" {
			continue
		}
		relativePath := filepath.Join(albumPath, p.name)
		p.ProjectPath = strings.ReplaceAll(string(filepath.Separator)+localToProjectPath(relativePath, relativeProjectPath), string(filepath.Separator), "/")
		p.setUiName()
		a.Photos = append(a.Photos, p)
	}
	if len(a.Photos) == 0 {
		return fmt.Errorf("no album pictures found")
	}
	if a.Title == "" {
		a.setTitle()
	}
	a.FileProjectPath = galleryPath + a.getFilename()
	return nil
}

func (p *photo) setUiName() {
	p.UiName = strings.ReplaceAll(p.name[:strings.LastIndex(p.name, ".")], "_", " ")
}

func (a *album) setTitle() {
	dir := filepath.Dir(a.Photos[0].ProjectPath)
	foldername := dir[1+strings.LastIndex(dir[:len(dir)-2], string(filepath.Separator)):]
	splitHyph := strings.Split(foldername, "-")
	titleName := strings.Join(splitHyph, " - ")
	a.Title = strings.ReplaceAll(titleName, "_", " ")
}

func localToProjectPath(path, rpp string) string {
	prjpath, err := filepath.Rel(rpp, path)
	if err != nil {
		panic(err)
	}
	return prjpath
}
