package main

import (
	"context"
	"encoding/json"
	"fmt"
	"gitee-cmd/utils"
)

//
// searchRepos

// @Description: 查询指定指定的仓库列表
// @param owner
// @param branchName
// @return []string
// @return error
func searchRepos(owner, branchName string) ([]string, error) {
	uri := fmt.Sprintf("/enterprises/%s/repos", owner)

	type repoItemResult struct {
		FullName string `json:"full_name"`
	}

	page := 1
	size := 100
	var repos []string
	ctx := context.Background()
	for {
		// 构造请求参数
		params := map[string]interface{}{
			"page":     page,
			"per_page": size,
		}
		var repoItems []repoItemResult

		// 发送GET请求获取仓库列表
		result, err := utils.NewGiteeClient("GET", uri, utils.WithContext(ctx), utils.WithQuery(params)).Do()
		if err != nil {
			return nil, err
		}
		body, err := result.GetRespBody()
		if err != nil {
			return nil, err
		}

		json.Unmarshal(body, &repoItems)
		if len(repoItems) == 0 {
			break
		}

		var semaphore chan struct{}
		var resultChan chan string
		resultChan = make(chan string, len(repoItems))
		semaphore = make(chan struct{}, 10)

		for _, item := range repoItems {
			if len(item.FullName) == 0 {
				continue
			}
			explode := utils.StrExplode(item.FullName, "/")
			if len(explode) < 2 {
				continue
			}
			repo := explode[1]

			go func(owner, repo, branchName string) {
				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				// 检查仓库中是否存在指定分支
				branchExists := checkBranchExists(owner, repo, branchName)
				if branchExists {
					resultChan <- repo
				} else {
					resultChan <- ""
				}
			}(owner, repo, branchName)
		}

		//收集结果
		for i := 0; i < len(repoItems); i++ {
			select {
			case repoName := <-resultChan:
				if repoName != "" {
					repos = append(repos, repoName)
				}
			}
		}

		page++
	}

	return repos, nil
}

// checkBranchExists
//
//	@Description: 仓库是否存在指定分支
//	@param owner
//	@param repo
//	@param branchName
//	@return bool
func checkBranchExists(owner, repo, branchName string) bool {
	uri := fmt.Sprintf("/repos/%s/%s/branches/%s", owner, repo, branchName)
	params := map[string]interface{}{}
	ctx := context.Background()

	// 发送GET请求获取仓库列表
	_, err := utils.NewGiteeClient("GET", uri, utils.WithContext(ctx), utils.WithPayload(params)).Do()
	if err != nil {
		return false
	}

	return true
}

//
// createBranch

// @Description: 创建分支
// @param owner
// @param repo
// @param branchName
// @param refs
// @return string
// @return error
func createBranch(owner, repo, branchName, refs string) (string, error) {
	if branchName == refs {
		return branchName, nil
	}

	uri := fmt.Sprintf("/repos/%s/%s/branches", owner, repo)
	params := map[string]interface{}{
		"refs":        refs,
		"branch_name": branchName,
	}
	ctx := context.Background()

	type createBranchResult struct {
		Name string `json:"name"`
	}
	var createBranch createBranchResult

	result, err := utils.NewGiteeClient("POST", uri, utils.WithContext(ctx), utils.WithPayload(params)).Do()
	if err != nil {
		return "", err
	}

	body, err := result.GetRespBody()
	if err != nil {
		return "", err
	}
	json.Unmarshal(body, &createBranch)

	return createBranch.Name, nil
}

// createPr
//
//	@Description: 创建pr
//	@param repo
//	@param head
//	@param base
//	@return string
//	@return error
func createPr(owner, repo, head, base string) (string, error) {
	uri := fmt.Sprintf("/repos/%s/%s/pulls", owner, repo)
	params := map[string]interface{}{
		"head":  head,
		"base":  base,
		"title": fmt.Sprintf("创建pr:%s到%s", head, base),
	}
	ctx := context.Background()

	type createPrResult struct {
		HtmlUrl string `json:"html_url"`
	}
	var createPrRes createPrResult

	result, err := utils.NewGiteeClient("POST", uri, utils.WithContext(ctx), utils.WithPayload(params)).Do()
	if err != nil {
		return "", err
	}

	body, err := result.GetRespBody()
	if err != nil {
		return "", err
	}
	json.Unmarshal(body, &createPrRes)

	return createPrRes.HtmlUrl, nil
}

// prList
//
//	@Description: 获取pr列表
//	@param owner
//	@param repo
//	@param headBranch
//	@return []int
//	@return error
func prList(owner, repo, headBranch string) ([]int, error) {
	uri := fmt.Sprintf("/repos/%s/%s/pulls", owner, repo)
	params := map[string]interface{}{
		"state":    "open",
		"head":     headBranch,
		"page":     1,
		"per_page": 10,
	}
	ctx := context.Background()

	type prListResult struct {
		Id      int    `json:"number"`
		HtmlUrl string `json:"html_url"`
	}
	var prListRes []prListResult

	result, err := utils.NewGiteeClient("GET", uri, utils.WithContext(ctx), utils.WithQuery(params)).Do()
	if err != nil {
		return nil, err
	}

	body, err := result.GetRespBody()
	if err != nil {
		return nil, err
	}
	json.Unmarshal(body, &prListRes)

	prIds := make([]int, len(prListRes))
	for index, pr := range prListRes {
		prIds[index] = pr.Id
	}

	return prIds, nil
}

// reviewPr
//
//	@Description: pr审核or测试
//	@param owner
//	@param repo
//	@param action
//	@param prNum
//	@return error
func reviewPr(owner, repo, action string, prNum int) error {
	uri := fmt.Sprintf("/repos/%s/%s/pulls/%d/%s", owner, repo, prNum, action) // review
	params := map[string]interface{}{
		"force": true,
	}
	ctx := context.Background()

	_, err := utils.NewGiteeClient("POST", uri, utils.WithContext(ctx), utils.WithPayload(params)).Do()
	if err != nil {
		return err
	}

	return nil
}

// mergePR
//
//	@Description: 合并pr
//	@param owner
//	@param repo
//	@param prNum
//	@return error
func mergePR(owner, repo string, prNum int) error {
	uri := fmt.Sprintf("/repos/%s/%s/pulls/%d/merge", owner, repo, prNum)
	params := map[string]interface{}{}
	ctx := context.Background()

	_, err := utils.NewGiteeClient("PUT", uri, utils.WithContext(ctx), utils.WithPayload(params)).Do()
	if err != nil {
		return err
	}

	return nil
}

// createTag
//
//	@Description: 创建tag
//	@param owner
//	@param repo
//	@param tagName
//	@param refs
//	@return error
func createTag(owner, repo, tagName, refs string) error {
	uri := fmt.Sprintf("/repos/%s/%s/tags", owner, repo)
	params := map[string]interface{}{
		"tag_name": tagName,
		"refs":     refs,
	}
	ctx := context.Background()

	_, err := utils.NewGiteeClient("POST", uri, utils.WithContext(ctx), utils.WithPayload(params)).Do()
	if err != nil {
		return err
	}

	return nil
}
