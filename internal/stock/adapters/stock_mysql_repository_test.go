package adapters

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	_ "github.com/getmelove/gorder2/internal/common/config"
	"github.com/getmelove/gorder2/internal/stock/entity"
	"github.com/getmelove/gorder2/internal/stock/infrastructure/persistent"
	"github.com/getmelove/gorder2/internal/stock/infrastructure/persistent/builder"
	"github.com/spf13/viper"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *persistent.MySQL {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		"",
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)

	// 注意：由于测试并行运行 (t.Parallel())，这里不能简单地 DROP/CREATE 数据库，
	// 否则一个测试正在运行时，会被另一个测试的 setup 删除数据库。
	// 为了支持 t.Parallel()，我们需要为每个测试生成一个唯一的数据库名，或者不使用 t.Parallel()。
	// 鉴于目前代码逻辑是复用同一个 shadow 库配置，最简单的修复是确保 setup 操作是幂等的，或者在此处移除 t.Parallel()。
	// 但用户坚持要测试 t.Parallel()，所以即使应用层有锁，测试层的 setup 甚至在应用代码运行前就打架了。
	// 真正的并行测试修复需要唯一的数据库实例。

	testDB := viper.GetString("mysql.dbname") + "_shadow"

	// 加锁防止 setupTestDB 自身的并发冲突（但这不能解决 Drop 掉别人正在用的库的问题）
	// 正确的做法是完全独立的库。
	// 在没有独立库的情况下，只能串行初始化，但这样 t.Parallel 意义就不大了。
	// 这里为了演示问题，我们保留原样，但在解释中说明这是 Test Setup 的 Race，不是业务逻辑的 Race。

	// FIX: 为了让 t.Parallel 真正工作且互不干扰，我们必须使用唯一的数据库名。
	uniqueDBName := fmt.Sprintf("%s_%d", testDB, time.Now().UnixNano())

	assert.NoError(t, db.Exec("CREATE DATABASE IF NOT EXISTS "+uniqueDBName).Error)

	// 测试结束时清理
	t.Cleanup(func() {
		_ = db.Exec("DROP DATABASE IF EXISTS " + uniqueDBName)
	})

	dsn = fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		uniqueDBName,
	)
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(&persistent.StockModel{}))

	return persistent.NewMySQLWithDB(db)
}

func TestMySQLStockRepository_UpdateStock_Race(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	db := setupTestDB(t)

	// 准备初始数据
	var (
		testItem           = "item-1"
		initialStock int32 = 100
	)
	err := db.Create(ctx, &persistent.StockModel{
		ProductID: testItem,
		Quantity:  initialStock,
	})
	assert.NoError(t, err)

	repo := NewMySQLStockRepository(db)
	var wg sync.WaitGroup
	concurrentGoroutines := 10
	for i := 0; i < concurrentGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := repo.UpdateStock(
				ctx,
				[]*entity.ItemWithQuantity{
					{ID: testItem, Quantity: 1},
				}, func(ctx context.Context, existing, query []*entity.ItemWithQuantity) ([]*entity.ItemWithQuantity, error) {
					// 模拟减少库存
					var newItems []*entity.ItemWithQuantity
					for _, e := range existing {
						for _, q := range query {
							if e.ID == q.ID {
								newItems = append(newItems, &entity.ItemWithQuantity{
									ID:       e.ID,
									Quantity: e.Quantity - q.Quantity,
								})
							}
						}
					}
					return newItems, nil
				},
			)
			assert.NoError(t, err)
		}()
	}

	wg.Wait()
	res, err := db.BatchGetStockByID(ctx, builder.NewStock().ProductIDs(testItem))
	assert.NoError(t, err)
	assert.NotEmpty(t, res, "res cannot be empty")

	expectedStock := initialStock - int32(concurrentGoroutines)
	assert.Equal(t, expectedStock, res[0].Quantity)
}

func TestMySQLStockRepository_UpdateStock_OverSell(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	db := setupTestDB(t)

	// 准备初始数据
	var (
		testItem           = "item-1"
		initialStock int32 = 5
	)
	err := db.Create(ctx, &persistent.StockModel{
		ProductID: testItem,
		Quantity:  initialStock,
	})
	assert.NoError(t, err)

	repo := NewMySQLStockRepository(db)
	var wg sync.WaitGroup
	concurrentGoroutines := 100
	for i := 0; i < concurrentGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := repo.UpdateStock(
				ctx,
				[]*entity.ItemWithQuantity{
					{ID: testItem, Quantity: 1},
				}, func(ctx context.Context, existing, query []*entity.ItemWithQuantity) ([]*entity.ItemWithQuantity, error) {
					// 模拟减少库存
					var newItems []*entity.ItemWithQuantity
					for _, e := range existing {
						for _, q := range query {
							if e.ID == q.ID {
								newItems = append(newItems, &entity.ItemWithQuantity{
									ID:       e.ID,
									Quantity: e.Quantity - q.Quantity,
								})
							}
						}
					}
					return newItems, nil
				},
			)
			assert.NoError(t, err)
		}()
		time.Sleep(20 * time.Millisecond)
	}

	wg.Wait()
	res, err := db.BatchGetStockByID(ctx, builder.NewStock().ProductIDs(testItem))
	assert.NoError(t, err)
	assert.NotEmpty(t, res, "res cannot be empty")

	assert.GreaterOrEqual(t, res[0].Quantity, int32(0))
}
