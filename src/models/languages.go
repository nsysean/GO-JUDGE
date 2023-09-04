package models

type Languages struct {
	Languages []Language `json:"languages"`
}

type Language struct {
	Name string `json:"name"`
	Build [][]string `json:"build"`
	Run []string `json:"run"`
}