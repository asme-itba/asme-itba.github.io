package main

import (
	"fmt"
	"github.com/spf13/viper"
	"html/template"
	"os"
	"path/filepath"
)

//            <!--  search button     <form class="form-inline mt-2 mt-md-0">
//                    <input class="form-control mr-sm-2" type="text" placeholder="Search" aria-label="Search">
//                    <button class="btn btn-outline-success my-2 my-sm-0" type="submit">Search</button>
//                  </form> -->
const configPath = "generate.yaml"

type button struct {
	Name       string
	Href       string
	IsDropdown bool
	SubButtons []button
}

type headlines struct {
	Headline    string
	SubHeadline string
	Content     string
	LinkText    string
	Href        string
	PhotoURL    string
	PhotoAlt    string
}

type config struct {
	Team                team
	TeamFilePath        string
	PhotoAlbumPath      string
	AlbumTemplatePath   string
	GalleryTemplatePath string
	GalleryFilePath     string
	Albums              []album
	RelativeProjectPath string
	Index               struct {
		TemplatePath string
		Carousel     []headlines
		NavBar       []button
		Tiles        []headlines
	}
	Gallery struct {
		TemplatePath string
		Title        string
		Filepath     string
		NavBar       []button
		Path string
	}
}

func init() {
	viper.AddConfigPath(".")
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		fmt.Print(err)
	}
}

// init() runs first
func main() {
	if err := run(); err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	fmt.Println("Done building site!")
}

func run() error {
	var theConfig config
	err := viper.Unmarshal(&theConfig)
	if err != nil {
		return err
	}
	// shorthand the relative project path. if run from /src/ rel path is ../
	rpp := theConfig.RelativeProjectPath
	// Produce TEAM.HTML
	teamTemplate, err := template.New("team").Parse(file2str(filepath.Join(rpp, theConfig.Team.TemplatePath)))
	if err != nil {
		return err
	}
	//teamFile := filepath.Join(rpp, theConfig.TeamPath, theConfig.TeamFilename)
	fteam, err := os.Create(filepath.Join(rpp, theConfig.TeamFilePath))
	if err != nil {
		return err
	}
	defer fteam.Close()
	if len(theConfig.Team.Members) == 1 {
		theConfig.Team.Members[0].Style = `margin-left: auto;margin-right: auto;`
	}
	theConfig.Team.sanitize()
	err = teamTemplate.Execute(fteam, theConfig.Team)
	if err != nil {
		return err
	}
	// PRODUCE ALBUMS
	for i, a := range theConfig.Albums {
		err = a.completeAlbum(filepath.Join(rpp, theConfig.PhotoAlbumPath, a.Folder), rpp, theConfig.Gallery.Path)
		if err != nil {
			return err
		}
		f, err := os.Create(filepath.Join(rpp, a.FileProjectPath))
		if err != nil {
			return err
		}
		albumTemplate, err := template.New(a.Title).Parse(file2str(filepath.Join(rpp, theConfig.AlbumTemplatePath)))
		if err != nil {
			return err
		}
		err = albumTemplate.Execute(f, a)
		if err != nil {
			return err
		}
		theConfig.Albums[i] = a
	}
	// PRODUCE GALLERY
	galleryTemplate, err := template.New("gallery").Parse(file2str(filepath.Join(rpp, theConfig.Gallery.TemplatePath)))
	if err != nil {
		return err
	}
	fgallery, err := os.Create(filepath.Join(rpp, theConfig.GalleryFilePath))
	if err != nil {
		return err
	}
	err = galleryTemplate.Execute(fgallery, theConfig)
	if err != nil {
		return err
	}
	// PRODUCE INDEX.HTML this should be at the end!
	indexTemplate, err := template.New("index").Parse(file2str(filepath.Join(rpp, theConfig.Index.TemplatePath)))
	if err != nil {
		return err
	}
	findex, err := os.Create(filepath.Join(rpp, "index.html"))
	if err != nil {
		return err
	}
	defer findex.Close()
	err = indexTemplate.Execute(findex, theConfig)
	if err != nil {
		return err
	}
	return nil
}

const (
	bufsize = 1 << 8 // 1 << 8 == 256 buffer size. good size i think
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
