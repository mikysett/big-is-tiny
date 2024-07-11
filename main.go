package main

type BigChange struct {
	domains  []Domain
	settings Settings
}

type Settings struct {
	isDryRun   bool
	isDraftPrs bool
	verbose    bool
	repoPath   string
	platform   Platform
}

type Domain struct {
	name  string
	path  string
	teams []Team
}

type Team struct {
	teamUrl  string
	teamType Comunication
}

type Comunication int

const (
	Slack Comunication = iota
	Teams
)

type Platform int

const (
	Azure Platform = iota
	GitHub
)

func main() {

}
