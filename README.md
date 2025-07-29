# gitee cmd

## 安装

```
go install github.com/lizuoqiang/gitee-cmd@latest
```

```
export GITEE_ACCESS_TOKEN=xxx
```

```
gitee-cmd --help
```

## 功能

### 创建pr

`
1.搜索有develop-14.38.0分支的仓库

2.遍历以上仓库创建release-14.38.0

3.从develop-14.38.0创建pr到release-14.38.0
`

```
gitee-cmd -action create_pr -branch develop-14.38.0 -type release
```

### 合并pr

`
1.搜索有release-14.38.0分支的仓库

2.遍历以上仓库的pr列表，查询源分支为release-14.38.0的pr，合并pr

3.从master创建tag(如有)
`

```
gitee-cmd -action merge_pr -branch release-14.38.0 -tag v14.38.0
```