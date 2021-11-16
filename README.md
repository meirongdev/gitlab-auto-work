# gitlab-auto-work
automate some work about gitlab

# build project
- go get github.com/xanzy/go-gitlab
- go get -u github.com/spf13/cobra/cobra
- cobra init --pkg-name github.com/chengshidaomin/gitlab-auto-work

# Workflow
1. get all projects of the current user
2. new branch from tag or commit hash
3. replace .gitlab-ci.yml
4. update

# config.yml