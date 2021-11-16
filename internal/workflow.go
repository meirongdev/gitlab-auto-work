package internal

import (
	"fmt"
	"os"

	_ "embed"

	"github.com/chengshidaomin/gitlab-auto-work/internal/config"
	"github.com/rs/zerolog"
	"github.com/xanzy/go-gitlab"
	"golang.org/x/mod/module"
)

type Workflow struct {
	WorkConfig config.WorkConfig
	Client     *gitlab.Client
	Log        *zerolog.Logger
}

//go:embed gitlab/.gitlab-ci.yml
var defaultCiFile string

func (wf *Workflow) Run() {
	projects, err := wf.getCurrentUserProjects()
	if err != nil {
		wf.Log.Error().AnErr("GetCurrentUserProjects Error", err)
		os.Exit(2)
	}

	if len(projects) == 0 {
		wf.Log.Info().Msg("The user has no project")
		os.Exit(0)
	}
	var projectMap = make(map[string]*gitlab.Project, len(projects))
	for _, project := range projects {
		projectMap[project.WebURL] = project
	}
	fmt.Println(projects[2].WebURL)

	for _, repo := range wf.WorkConfig.Repositories {
		wf.Log.Info().Msgf("Processing repo[%s]- [%s] now", repo.Name, repo.Url)
		project := projectMap[repo.Url]
		wf.Log.Info().Msgf("%v", project)
		for _, version := range repo.Versions {
			if err := wf.processVersion(project, repo, version); err != nil {
				wf.Log.Error().Stack().AnErr("process version"+repo.Url+"-"+version+" error", err)
				continue
			}
			wf.Log.Info().Msgf("version %s-%s success", repo.Url, version)
		}
	}

}

func (wf *Workflow) processVersion(project *gitlab.Project, repo config.Repository, version string) error {
	var err error
	if module.IsPseudoVersion(version) {
		version, err = module.PseudoVersionRev(version)
		if err != nil {
			return err
		}
	}

	wf.Log.Debug().Msgf("project:[%v] repo.Url: [%s]", project, repo.Url)
	newBranchName := wf.WorkConfig.BranchPrefix + version
	wf.Log.Debug().Msgf("new branch name is %s", newBranchName)
	_, resp, err := wf.Client.Branches.GetBranch(project.ID, newBranchName)
	if err != nil {
		wf.Log.Info().Msgf("Get Branch %s-%s error code %d", project.WebURL, newBranchName, resp.StatusCode)
	}
	if resp.StatusCode == 404 {
		wf.Log.Info().Msgf("%s not exists, create a new branch", newBranchName)

		_, _, err = wf.Client.Branches.CreateBranch(project.ID, &gitlab.CreateBranchOptions{
			Branch: gitlab.String(newBranchName),
			Ref:    gitlab.String(version),
		})
		if err != nil {
			return err
		}
	}

	// check file
	content, _, err := wf.Client.RepositoryFiles.GetRawFile(project.ID, *gitlab.String(".gitlab-ci.yml"), &gitlab.GetRawFileOptions{
		Ref: gitlab.String(newBranchName),
	})
	if err != nil {
		return err
	}
	wf.Log.Debug().Msgf(".gitlab-ci.yml content\n%s", string(content))

	fileInfo, resp, err := wf.Client.RepositoryFiles.UpdateFile(project.ID, *gitlab.String(".gitlab-ci.yml"), &gitlab.UpdateFileOptions{
		Branch:        gitlab.String(newBranchName),
		AuthorEmail:   gitlab.String(wf.WorkConfig.UserEmail),
		AuthorName:    gitlab.String(wf.WorkConfig.UserName),
		Content:       gitlab.String(defaultCiFile),
		CommitMessage: gitlab.String(wf.WorkConfig.CommitMsg),
	})
	if err != nil {
		return err
	}
	wf.Log.Info().Msgf("version process success %s, %v", fileInfo.FilePath, resp.Body)
	return nil
}

func (wf *Workflow) getCurrentUserProjects() ([]*gitlab.Project, error) {
	user, _, err := wf.Client.Users.CurrentUser()
	if err != nil {
		return nil, err
	}

	projects, _, err := wf.Client.Projects.ListUserProjects(user.ID, &gitlab.ListProjectsOptions{
		MinAccessLevel: gitlab.AccessLevel(gitlab.DeveloperPermissions),
		Simple:         gitlab.Bool(true),
	})
	if err != nil {
		return nil, err
	}
	return projects, nil
}
