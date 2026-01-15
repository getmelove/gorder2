package tests

import (
	"context"
	"fmt"
	"log"
	"testing"

	sw "github.com/getmelove/gorder2/internal/common/client/order" // 引入生成的订单服务客户端代码
	_ "github.com/getmelove/gorder2/internal/common/config"        // 引入配置模块，确保 viper 配置被加载
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

var (
	ctx = context.Background()
	// server 变量用于构建目标服务器的地址。
	// 它从配置文件中读取 "order.http-addr" 配置项（例如 "localhost:8080"），
	// 并拼接成完整的 API 基础 URL。
	server = fmt.Sprintf("http://%s/api", viper.GetString("order.http-addr"))
)

// TestMain 是 Go 语言测试框架的入口函数。
// 在运行任何具体测试之前，它会先执行 setup 逻辑（这里是 before 函数），
// 然后调用 m.Run() 执行所有测试用例。
func TestMain(m *testing.M) {
	before()
	m.Run()
}

// before 函数用于执行测试前的初始化工作。
// 这里仅仅是打印服务器地址，但在更复杂的测试中，可以在这里初始化数据库连接或 mock 对象。
func before() {
	log.Printf("server=%s", server)
}

// TestCreateOrder_success 测试创建订单的成功场景。
func TestCreateOrder_success(t *testing.T) {
	// 调用 getResponse 辅助函数发送创建订单请求。
	// 这里模拟了一个合法的请求：
	// CustomerId 为 "123"
	// Items 包含一个商品，ID 为 "test-item-1"，数量为 1。
	response := getResponse(t, "123", sw.PostCustomerCustomerIdOrdersJSONRequestBody{
		CustomerId: "123",
		Items: []sw.ItemWithQuantity{
			{
				Id:       "prod_ThjsHE5qzUNVYT",
				Quantity: 1,
			},
			{
				Id:       "prod_TiCZTq9XoPLJ8V",
				Quantity: 1,
			},
		},
	})

	// 打印响应体，方便调试时查看。
	t.Logf("body=%s", string(response.Body))

	// 断言 HTTP 状态码为 200 (OK)。
	assert.Equal(t, 200, response.StatusCode())

	// 断言业务错误码 (Errno) 为 0。
	// 在该系统的设计中，Errno = 0 通常表示业务逻辑执行成功。
	assert.Equal(t, 0, response.JSON200.Errno)
}

// TestCreateOrder_invalidParams 测试参数无效时的场景。
func TestCreateOrder_invalidParams(t *testing.T) {
	// 构造一个无效的请求：Items 为 nil。
	response := getResponse(t, "123", sw.PostCustomerCustomerIdOrdersJSONRequestBody{
		CustomerId: "123",
		Items:      nil,
	})

	// 即使参数无效，HTTP 层面可能仍然返回 200 OK，
	// 具体的错误信息通过响应体中的业务字段返回。
	assert.Equal(t, 200, response.StatusCode())

	// 断言业务错误码 (Errno) 为 2。
	// Errno = 2 在这里预期对应于参数校验错误或其他特定的业务错误。
	assert.Equal(t, 2, response.JSON200.Errno)
}

// getResponse 是一个辅助函数，封装了发送 HTTP 请求的细节。
// 它接收测试对象 t、客户 ID 和请求体，返回解析后的响应对象。
func getResponse(t *testing.T, customerID string, body sw.PostCustomerCustomerIdOrdersJSONRequestBody) *sw.PostCustomerCustomerIdOrdersResponse {
	t.Helper() // 标记为辅助函数，这样测试失败时报错行号会指向调用者而不是这里。

	// 创建一个新的 API 客户端。
	// sw (client/order) 是根据 OpenAPI 规范生成的代码。
	client, err := sw.NewClientWithResponses(server)
	if err != nil {
		t.Fatal(err) // 如果创建客户端失败，终止测试。
	}

	// 调用生成的 PostCustomerCustomerIdOrdersWithResponse 方法发送 POST 请求。
	response, err := client.PostCustomerCustomerIdOrdersWithResponse(ctx, customerID, body)
	if err != nil {
		t.Fatal(err) // 如果网络请求失败，终止测试。
	}

	return response
}
