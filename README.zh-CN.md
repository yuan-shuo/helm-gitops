# Helm-GitOps

一个helm的扩展，为helm提供gitops相关辅助功能

## 使用

相对于常规自行创建helm chart，随后在内部粘贴.gitignore等文件，自行修改创建分支、提交、版本号，手动编写argocd所需的yaml，此扩展提供了较为舒适的简化方案：

```bash
# 1.创建: 在原 helm chart 基础上添加了.gitignore, 已经初始化后的.git仓库, .github/workflow/ci-test.yaml等文件
helm gitops create test

# 2.切换分支: 自动执行checkout进入一个开发分支
helm gitops checkout feature/foo

# ...常规编写chart

# 3.提交代码: 自动执行版本号增加(可选)并commit, 以及推送(可选)
helm gitops commit -m "fix: foo"

# 4.本地检查: 自动执行 helm lint + unittest
helm gitops lint

# 5.打标签(git tag): 自动基于Chart的版本号对当前内容打标签, 以及推送(可选)
helm gitops tag
```

## 安装







../bin/helm-gitops checkout feature/foo
../bin/helm-gitops commit -m "fix: foo"

./bin/helm-gitops create test

go build -o bin/helm-gitops .

- 创建
  create
- 进入开发分支
  checkout
- // 编写代码 ...
- 本地检查
  lint
- bump 提交
  commit
- 自动tag
  tag