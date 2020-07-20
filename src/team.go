package main

const placeHolderGray = "https://upload.wikimedia.org/wikipedia/commons/thumb/8/86/Solid_grey.svg/512px-Solid_grey.svg.png"

type team struct {
	Members      []member
	TemplatePath string
	Filename     string
	NavBar       []button
}

type member struct {
	Name        string
	Description string
	PhotoURL    string
	LinkText    string
	LinkURL     string
	Style       string
}

func (t *team) sanitize() {
	for i, m := range t.Members {
		if m.PhotoURL == "" {
			t.Members[i].PhotoURL = placeHolderGray
		}
	}
}
