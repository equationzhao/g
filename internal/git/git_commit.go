package git

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	strftime "github.com/itchyny/timefmt-go"
)

type CommitInfo struct {
	Hash string `json:"h"`

	Committer     string `json:"c"`
	CommitterDate string `json:"cd"`

	Author     string `json:"a"`
	AuthorDate string `json:"ad"`
}

func (c CommitInfo) GetCommitterDateInFormat(format string) string {
	t, err := time.Parse(time.RFC3339, c.CommitterDate)
	if err != nil {
		return ""
	}
	return t.Format(format)
}

func (c CommitInfo) GetAuthorDateInFormat(format string) string {
	t, err := time.Parse(goParseFormat, c.AuthorDate)
	if err != nil {
		return ""
	}
	if after, ok := strings.CutPrefix(format, "+"); ok {
		return strftime.Format(t, after)
	}
	return t.Format(format)
}

var NoneCommitInfo = CommitInfo{"-", "-", "-", "-", "-"}

// https://github.com/chaqchase/lla/blob/main/plugins/last_git_commit/src/lib.rs
func GetLastCommitInfo(path string) (*CommitInfo, error) {
	return getLastCommitInfo(path)
}

const (
	gitDateFormat = `format:"%Y-%m-%d %H:%M:%S.%9N %z"`
	goParseFormat = time.RFC3339
)

func getLastCommitInfo(path string) (*CommitInfo, error) {
	cmd := exec.Command("git", "log", "-1", `--pretty=format:{"h":"%h","a":"%an","c":"%cn","ad":"%aI","cd":"%cI"}`, fmt.Sprintf(`--date=%s`, gitDateFormat), path)

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	if len(output) == 0 {
		return &NoneCommitInfo, nil
	}
	var info CommitInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, err
	}

	return &info, nil
}
