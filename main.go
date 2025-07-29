package main

import (
	"flag"
	"fmt"
	"strings"
)

func main() {
	action := flag.String("action", "", "操作:create_pr,merge_pr")

	//merge_pr需要
	tag := flag.String("tag", "", "[merge_pr]合并pr时是否创建tag")

	//create_pr需要
	branch := flag.String("branch", "", "[create_pr|merge_pr]分支名|pr的源分支名")
	prType := flag.String("type", "release", "[create_pr]release:创建pr到release分支,master:到master分支")
	owner := flag.String("owner", "huanjutang", "[create_pr|merge_pr]仓库所有者")
	flag.Parse()

	fmt.Println(*action, *branch, *prType, *owner, *tag)
	err := run(*action, *branch, *prType, *owner, *tag)
	if err != nil {
		fmt.Println("run error", err)
		return
	}
}

func run(action, branch, prType, owner, tag string) error {
	repos, err := searchRepos(owner, branch)
	if err != nil {
		fmt.Println("搜索仓库失败")
		return err
	}

	switch action {
	case "create_pr":
		newBranch := "master"
		if prType == "release" {
			newBranch = strings.Replace(branch, "develop", "release", 1)
		}
		for _, repo := range repos {
			_, err := createBranch(owner, repo, newBranch, "master")
			if err != nil {
				fmt.Println("创建分支失败:", repo, err)
				continue
			}
			prUrl, err := createPr(owner, repo, branch, newBranch)
			if err != nil {
				fmt.Println("创建PR失败:", repo, err)
				continue
			}
			fmt.Println(repo, " PR 创建成功:", prUrl)
		}
	case "merge_pr":
		for _, repo := range repos {
			list, err := prList(owner, repo, branch)
			if err != nil {
				fmt.Println("获取PR列表失败:", repo, err)
				continue
			}

			//合并pr
			for _, prId := range list {
				err := reviewPr(owner, repo, "review", prId)
				if err != nil {
					fmt.Println("审核PR失败:", repo, prId, err)
					continue
				}
				err = reviewPr(owner, repo, "test", prId)
				if err != nil {
					fmt.Println("测试PR失败:", repo, prId, err)
					continue
				}
				err = mergePR(owner, repo, prId)
				if err != nil {
					fmt.Println("合并PR失败:", repo, prId, err)
					continue
				}
				if tag != "" {
					err = createTag(owner, repo, tag, "master")
					if err != nil {
						fmt.Println("创建Tag失败:", repo, tag, err)
					}
				}
				fmt.Println("PR合并处理完成:", repo, prId)
			}
		}
	}

	return nil
}
