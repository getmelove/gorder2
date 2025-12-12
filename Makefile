# 声明伪目标，避免同名文件与时间戳干扰
.PHONY: gen
# 组合目标：一次生成 proto 和 openapi 代码/文档
gen: genproto genopenapi

.PHONY: genproto
# 生成 Protobuf 相关代码
genproto:
	@./scripts/genproto.sh

.PHONY: genopenapi
# 生成 OpenAPI 相关产物（如客户端或文档）
genopenapi:
	@./scripts/genopenapi.sh
