package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/google/go-github/github"
	h "github.com/rodaine/hclencoder"
	"golang.org/x/oauth2"
)

type config struct {
	collaborator *[]collaborator `hcl:"resource github_repository_collaborator"`
}

type collaborator struct {
	name       string `hcl:",key"`
	repository string `hcl:",repository"`
	username   string `hcl:"username"`
	permission string `hcl:"permission"`
}

func main() {
	repoNames, err := getOrgRepo("StreetEasy")
	if err != nil {
		fmt.Println(err)
	}
	genTerraform(repoNames)

}

func getcollaborators(repoName string) []collaborator {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "TOKEN"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	var allcollaborator []collaborator

	opt := &github.ListCollaboratorsOptions{Affiliation: "direct"}

	collaborators, _, err := client.Repositories.ListCollaborators(ctx, "StreetEasy", repoName, opt)

	if err != nil {
		fmt.Println(err)
	}

	for i := range collaborators {
		perm := mapSearch(collaborators[i].Permissions)
		// fmt.Printf("%v:%v,%v\n", *collaborators[i].Login, perm, collaborators[i].Permissions)
		user := new(collaborator)
		user.name = repoName + "-" + *collaborators[i].Login
		user.repository = repoName
		user.username = *collaborators[i].Login
		user.permission = perm
		allcollaborator = append(allcollaborator, *user)
	}
	list := &github.ListOptions{}

	// teams, _, err := client.Repositories.ListTeams(ctx, "StreetEasy", repoName, list)
	//
	// if err != nil {
	// 	fmt.Println(err)
	// }
	//
	// for i := range teams {
	// 	// fmt.Printf("%v:%v\n", *teams[i].Name, *teams[i].Permission)
	// 	// perm := mapSearch(*teams[i].Permission)
	// 	user := new(collaborator)
	// 	user.name = repoName + "-" + *teams[i].Name
	// 	user.repository = repoName
	// 	user.username = *teams[i].Name
	// 	user.permission = *teams[i].Permission
	// 	allcollaborator = append(allcollaborator, *user)
	// }

	invite, _, err := client.Repositories.ListInvitations(ctx, "StreetEasy", repoName, list)

	if err != nil {
		fmt.Println(err)
	}

	for i := range invite {

		newPerm := invite[i].Permissions

		switch *newPerm {
		case "read":
			*newPerm = "pull"

		case "write":
			*newPerm = "push"
		}
		user := new(collaborator)
		user.name = repoName + "-" + *invite[i].Invitee.Login
		user.repository = repoName
		user.username = *invite[i].Invitee.Login
		user.permission = *newPerm
		allcollaborator = append(allcollaborator, *user)
	}
	return allcollaborator

}

func mapSearch(check *map[string]bool) string {
	count := 0
	for _, value := range *check {
		if value == true {
			count++

		}
	}

	switch count {
	case 1:
		return "pull"
	case 2:
		return "push"
	case 3:
		return "admin"

	}

	return ""

}

func genHCL(allcollaborator []collaborator) (string, error) {

	input := config{
		collaborator: &allcollaborator,
	}

	hcl, err := h.Encode(input)
	if err != nil {
		log.Fatal("unable to encode: ", err)
	}

	// fmt.Print(string(hcl))
	return string(hcl), nil

}

func genTerraform(repoNames []string) {
	fmt.Println(len(repoNames))
	path := "github-private-repos-collaborators.tf"
	os.Remove(path)
	createFile(path)

	for _, i := range repoNames {
		hcltemp := getcollaborators(i)
		hcl, err := genHCL(hcltemp)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(hcl)
		writeFile(path, hcl)

		for j := range hcltemp {
			cmd := "terragrunt import github_repository_collaborator." + hcltemp[j].name + " " + hcltemp[j].repository + ":" + hcltemp[j].username
			fmt.Println(string(runcmd(cmd, true)))
		}
	}

}

func createFile(path string) {
	// detect if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
	}
}

func writeFile(path, text string) {
	// open file using READ & WRITE permission
	var file, err = os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	// write some text to file
	_, err = file.WriteString(text)
	if err != nil {
		fmt.Println(err)
	}

	// save changes
	err = file.Sync()
	if err != nil {
		fmt.Println(err)
	}
}

func runcmd(cmd string, shell bool) []byte {
	if shell {
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			log.Fatal(err)
			panic("some error found")
		}
		return out
	}
	out, err := exec.Command(cmd).Output()
	if err != nil {
		log.Fatal(err)
	}
	return out
}

func getOrgRepo(org string) ([]string, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "TOKEN"},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.RepositoryListByOrgOptions{
		Type:        "all",
		ListOptions: github.ListOptions{PerPage: 30},
	}

	var allRepos []*github.Repository
	var reponames []string
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, org, opt)
		if err != nil {
			fmt.Println(err)
		}
		allRepos = append(allRepos, repos...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	for i := range allRepos {
		reponames = append(reponames, *allRepos[i].Name)
	}

	return reponames, nil
}
