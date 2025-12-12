#!/usr/bin/env bash

# 严格模式：命令出错、未定义变量、管道出错即退出
set -euo pipefail

# 允许 ** 递归匹配
shopt -s globstar

# 约束：必须在仓库根目录执行脚本
if ! [[ "$0" =~ scripts/genproto.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

# 引入通用日志/辅助函数
source ./scripts/lib.sh

# proto 文件所在根目录
API_ROOT="./api"

function dirs {
  # 收集所有 .proto 所在目录的末级目录名（去重），如 orderpb
  dirs=()
  while IFS= read -r dir; do
    dirs+=("$dir")
  done < <(find . -type f -name "*.proto" -exec dirname {} \; | xargs -n1 basename | sort -u)
  echo "${dirs[@]}"
}

function pb_files {
  # 列出仓库中的所有 .proto 文件路径
  pb_files=$(find . -type f -name '*.proto')
  echo "${pb_files[@]}"
}

function gen_for_modules() {
  # protoc 输出目录
  local go_out="internal/common/genproto"
  if [ -d "$go_out" ]; then
    # 若已有生成结果，先清理整个输出目录
    log_warning "found existing $go_out, cleaning all files under it"
    run rm -rf $go_out
  fi

  for dir in $(dirs); do
    # 目录名如 orderpb；去掉末尾 "pb" 得到 service 前缀（order）
    echo "dir=$dir"
    local service="${dir:0:${#dir}-2}"
    local pb_file="${service}.proto"

    # 确保输出子目录存在且是空的
    if [ -d "$go_out/$dir" ]; then
        log_warning "found existing $go_out/$dir, cleaning all files under it"
        run rm -rf "$go_out"/$dir/*
    else
      run mkdir -p "$go_out/$dir"
    fi
    log_info "generating code for $service to $go_out/$dir"

    # 生成 Go 和 gRPC 代码，路径保持 source_relative 以便 go module 引用
    run protoc \
      -I="/usr/local/include/" \
      -I="${API_ROOT}" \
      "--go_out=${go_out}" --go_opt=paths=source_relative \
      --go-grpc_opt=require_unimplemented_servers=false \
      "--go-grpc_out=${go_out}" --go-grpc_opt=paths=source_relative \
      "${API_ROOT}/${dir}/$pb_file"
  done
  log_success "protoc gen done!"
}

echo "directories containing protos to be built: $(dirs)"
echo "found pb_files: $(pb_files)"
gen_for_modules
