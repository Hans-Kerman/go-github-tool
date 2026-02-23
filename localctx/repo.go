package localctx

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v6"
)

func FindRepo() (repo *git.Repository, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		err = fmt.Errorf("获取当前目录失败: %w", err)
		return
	}
	for {
		repo, err = git.PlainOpen(cwd)
		if err == nil {
			return
		} else if errors.Is(err, git.ErrRepositoryNotExists) {
			parent := filepath.Dir(cwd)
			if parent == cwd {
				err = fmt.Errorf("未找到git仓库")
				return
			} else {
				cwd = parent
			}
		} else {
			err = fmt.Errorf("回溯git仓库时出现错误: %w", err)
		}
	}
}
