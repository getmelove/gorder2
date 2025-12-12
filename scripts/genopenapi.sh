#!/usr/bin/env bash

# 严格模式：命令出错、未定义变量、管道出错即退出
set -euo pipefail

# 允许 ** 递归匹配
shopt -s globstar

# 约束：必须在仓库根目录执行脚本
if ! [[ "$0" =~ scripts/genopenapi.sh ]]; then
  echo "must be run from repository root"
  exit 255
fi

# 引入通用日志/辅助函数
source ./scripts/lib.sh

# OpenAPI 文件所在根目录
OPENAPI_ROOT="./api/openapi"

# 选用的 server 代码生成器（数组仅允许一个）
GEN_SERVER=(
#  "chi-server"
#  "echo-server"
  "gin-server"
)

# 确保只启用了一个 server 生成器，避免混用
if [ "${#GEN_SERVER[@]}" -ne 1 ]; then
  log_error "GEN_SERVER enables more than 1 server, please check."
  exit 255
fi

log_callout "Using ${GEN_SERVER[0]}"

function openapi_files {
  # 列出 OpenAPI 文件名（当前未使用，仅保留参考）
  openapi_files=$(ls ${OPENAPI_ROOT})
  echo "${openapi_files[@]}"
}

# 按参数生成代码：output_dir 输出目录；package 包名；service OpenAPI 文件前缀
function gen() {
  local output_dir=$1
  local package=$2
  local service=$3

  # 准备输出目录，清理旧的生成文件
  run mkdir -p "$output_dir"
  run find "$output_dir" -type f -name "*.gen.go" -delete

  # 准备客户端输出目录（工具函数在 lib.sh 中定义）
  prepare_dir "internal/common/client/$service"

  # 生成仅包含类型定义的代码（服务端包）
  run oapi-codegen -generate types -o "$output_dir/openapi_types.gen.go" -package "$package" "api/openapi/$service.yml"
  # 生成 server 框架代码（根据选择的 gin-server）
  run oapi-codegen -generate "$GEN_SERVER" -o "$output_dir/openapi_api.gen.go" -package "$package" "api/openapi/$service.yml"

  # 生成客户端调用代码（client + types），放在 client 目录
  run oapi-codegen -generate client -o "internal/common/client/$service/openapi_client.gen.go" -package "$service" "api/openapi/$service.yml"
  run oapi-codegen -generate types -o "internal/common/client/$service/openapi_types.gen.go" -package "$service" "api/openapi/$service.yml"
}

gen internal/order/ports ports order

log_success "openapi generate success!"
