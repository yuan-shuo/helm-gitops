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

# 3.提交代码: 自动执commit, 以及推送(可选)
helm gitops commit -m "fix: foo"

# 4.本地检查: 自动执行 helm lint + unittest
helm gitops lint

# 5.推送代码: 先确认当前分支不是主分支, 随后lint检查, 最后执行push
helm gitops push
```

## 安装








go build -o bin/gitops .
./bin/gitops create test
cd test
../bin/gitops checkout feature/foo
../bin/gitops commit -m "fix: foo"
git remote add origin https://gitee.com/yuan-shuo188/helm-test1.git && git push -u origin main
../bin/gitops push
进入 feature/foo 分支gitee页面, 提交PR, 审核+测试, 合并到main分支
