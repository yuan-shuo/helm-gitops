#!/bin/bash

helm plugin install https://github.com/yuan-shuo/helm-gitops --verify=false
curl -sSL \
  https://github.com/yuan-shuo/helm-gitops/releases/download/v0.5.1/helm-gitops_0.5.1_linux_amd64.tar.gz \
  -o helm-gitops.tgz
tar -xzf helm-gitops.tgz

helm plugin install https://github.com/yuan-shuo/helm-gitops/releases/download/v0.5.1/helm-gitops_0.5.1_linux_amd64.tar.gz


# # ./scripts/install-local.sh

# # 构建gitops二进制文件
# go build -o bin/gitops .

# # 创建测试项目
# ./bin/gitops create test-nor
# ./bin/gitops create test-act --actions

# # 进入测试项目目录
# cd test-nor
# cd test-act

# # 链接远程测试仓库

# ## gitee/no-action
# git remote add origin https://gitee.com/yuan-shuo188/helm-test1.git

# ## github-with-action
# git remote add origin https://github.com/yuan-shuo/helmci-test1.git

# git push -u origin main

# # 测试checkout创建分支功能
# ../bin/gitops checkout feature/test1

# # 假设修改chart
# echo 'test file new1' > ./templates/new1.txt

# # 测试lint本地检查功能
# ../bin/gitops lint

# # 测试commit基础功能: 本地提交
# ../bin/gitops commit -m "fix:foo1"

# # 测试push基础功能: 推送本地分支到远程仓库
# ../bin/gitops push

# # 测试checkout同步主分支功能: 从主分支拉取最新代码 + 基于主分支创建新分支
# ../bin/gitops checkout feature/test2 -s

# # 假设修改chart
# echo '# test file new2' >> ./templates/service.yaml

# # 测试commit高级功能: 自动push + 创建PRcommit标记, gitee显示提交结果 -> fix:foo2 [create-pr]
# # ../bin/gitops commit -m "fix:foo2" --push --pr
# # 6. 版本管理
# ../bin/gitops version # 仅查询当前版本
# ../bin/gitops version -m pr -l patch # 传统 PR 模式（先开 release 分支 → 提 PR → CI 自动 tag）
# ../bin/gitops version -m main -l patch # 快捷主分支模式（直接 commit + tag + 同时推送）

# # 测试打印版本功能
# ../bin/gitops version

# # 测试版本升级功能: 升级patch|minor|major版本
# ../bin/gitops version --bump patch
# ../bin/gitops version --bump minor
# ../bin/gitops version --bump major

# # 把远端已经不存在的分支从本地 remote-tracking 里去掉(未测试)
# git fetch --prune # 明显缺少本地分支清理功能，遍历，然后一个一个询问y/N是否删除