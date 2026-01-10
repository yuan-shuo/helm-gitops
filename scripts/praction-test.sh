go build -o bin/gitops .
./bin/gitops create test-act --actions
cd test-act
git remote add origin https://github.com/yuan-shuo/helmci-test1.git
git push -u origin main
../bin/gitops version -m pr -l patch # 传统 PR 模式（先开 release 分支 → 提 PR → CI 自动 tag）