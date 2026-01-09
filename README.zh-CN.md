# Helm-GitOps

一个helm的扩展，为helm提供gitops相关辅助功能

## 使用

### git相关功能

相对于常规自行创建helm chart，随后在内部粘贴.gitignore等文件，自行修改创建分支、提交、版本号等一系列繁琐操作，此扩展提供了较为舒适的简化方案：

```bash
# 1. 创建：在原 helm chart 基础上添加 GitOps 骨架
helm gitops create test
helm gitops create test --actions        # 同时生成 .github/workflows/<>.yaml


# 2. 切换分支：自动进入开发分支（可选同步主分支）
helm gitops checkout feature/foo
helm gitops checkout feature/foo -s       # 先 pull origin/main 再创建分支


# 3. 提交代码：add + commit + 可选 push & 自动 PR
helm gitops commit -m "fix: foo"                    # 本地提交
helm gitops commit -m "feat: xxx" --push            # 提交并推送
helm gitops commit -m "ci: update" --pr --push      # 提交 + 推送 + 自动提 PR（含 [create-pr] commit 标记）


# 4. 本地检查：helm lint + helm unittest
helm gitops lint


# 5. 推送代码：lint → push（保护分支拦截）
helm gitops push                                      # 推送到 origin/当前分支


# 6. 版本管理
helm gitops version # 仅查询当前版本
helm gitops version -m pr -l patch # 传统 PR 模式（先开 release 分支 → 提 PR → CI 自动 tag）
helm gitops version -m main -l patch # 快捷主分支模式（直接 commit + tag + 同时推送）
# --mode=main|pr
# --level=patch|minor|major
```

### 环境仓库功能

待开发

### argocd功能

待开发

## 安装

### 使用 helm plugin install

```bash
helm plugin install https://github.com/yuan-shuo/helm-gitops
```

### 使用二进制文件

- 前往：[Releases · yuan-shuo/helm-gitops](https://github.com/yuan-shuo/helm-gitops/releases)下载对应操作系统的二进制文件

- 随后将解压得到的gitops二进制文件放在`$HELM_PLUGIN_DIR/bin/` 目录下

- 给予gitops二进制文件执行权限

```bash
chmod +x $HELM_PLUGIN_DIR/bin/gitops
```

## 环境需求

* git（version>=2.23）
* helm
* helm-unitest (可选，可以通过执行 **`helm plugin install https://github.com/helm-unittest/helm-unittest --verify=false`** 来安装)
