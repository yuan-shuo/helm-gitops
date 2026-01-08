# Helm-GitOps

一个helm的扩展，为helm提供gitops相关辅助功能

## 使用

相对于常规自行创建helm chart，随后在内部粘贴.gitignore等文件，自行修改创建分支、提交、版本号，手动编写argocd所需的yaml，此扩展提供了较为舒适的简化方案：

```bash
# 1. 创建：在原 helm chart 基础上添加 GitOps 骨架
helm gitops create test
helm gitops create demo --actions        # 同时生成 .github/workflows/ci-test.yaml


# 2. 切换分支：自动进入开发分支（可选同步主分支）
helm gitops checkout feature/foo
helm gitops checkout hotfix/bar -s       # 先 pull origin/main 再创建分支


# 3. 提交代码：add + commit + 可选 push & 自动 PR
helm gitops commit -m "fix: foo"                    # 本地提交
helm gitops commit -m "feat: xxx" --push            # 提交并推送
helm gitops commit -m "ci: update" --pr --push      # 提交 + 推送 + 自动提 PR（含 [create-pr] commit 标记）


# 4. 本地检查：helm lint + helm unittest
helm gitops lint


# 5. 推送代码：lint → push（保护分支拦截）
helm gitops push                                      # 推送到 origin/当前分支
helm gitops push --remote ci                          # 推送到名为 ci 的远程仓库


# 6. 版本管理：读版本号 & 一键毕业发布
helm gitops version                                   # 打印当前 Chart 版本
helm gitops version --bump patch|minor|major                      # 一键毕业：创建 release/vx.y.z 分支 → 改版本 → commit → push → PR
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
