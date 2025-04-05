#!/bin/bash

godoc -http=:6060 & # バックグラウンド実行
RUNNING_PID=$!      # godocのPIDを取得
sleep 3            # godocの実行を3秒待機
wget -np -k -p -q -r -E http://localhost:6060/pkg/nand2tetris-go/?m=all
kill ${RUNNING_PID} # godocの実行を終了
find localhost:6060 -name "index.html" -delete
rm -rf docs/lib/*
rm -rf docs/pkg/*
mv -f localhost:6060/lib docs
mv -f localhost:6060/pkg docs
# rename docs/pkg/nand2tetris-go/index.html?m=all.html to docs/pkg/nand2tetris-go/index.html
mv -f docs/pkg/nand2tetris-go/index.html?m=all.html docs/pkg/nand2tetris-go/index.html
rm -rf localhost:6060
