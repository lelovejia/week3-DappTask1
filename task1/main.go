package main

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"log"
	"time"
)

func main() {
	// 设置配置文件名和路径
	viper.SetConfigName("config") // 配置文件名称 (不带扩展名)
	viper.SetConfigType("yaml")   // 如果配置文件名没有扩展名，需要设置此项
	viper.AddConfigPath(".")      // 查找配置文件所在的路径
	// 自动读取环境变量
	viper.AutomaticEnv()
	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
	KEY := viper.GetString("KEY")
	//使用 ethclient 连接到 Sepolia 测试网络。
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/9974e485af4f472697f238af9f10f98e" + KEY) //自己的API
	if err != nil {
		log.Fatal(err)
	}
	//实现查询指定区块号的区块信息，包括区块的哈希、时间戳、交易数量等。
	// 获取区块信息
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to retrieve block: %v", err)
	}

	// 4. 打印区块详细信息
	fmt.Println("\n========== 区块信息 ==========")
	fmt.Printf("区块号(Number): %d\n", block.Number().Int64())
	fmt.Printf("区块哈希(Hash): %s\n", block.Hash().Hex())
	fmt.Printf("父区块哈希(Parent Hash): %s\n", block.ParentHash().Hex())
	fmt.Printf("时间戳(Timestamp): %s (%d)\n",
		time.Unix(int64(block.Time()), 0).Format("2006-01-02 15:04:05"),
		block.Time())
	fmt.Printf("交易数量(Transactions): %d\n", len(block.Transactions()))
	fmt.Printf("矿工地址(Miner): %s\n", block.Coinbase().Hex())
	fmt.Printf("难度(Difficulty): %d\n", block.Difficulty().Int64())
	fmt.Printf("Gas限制(Gas Limit): %d\n", block.GasLimit())
	fmt.Printf("已用Gas(Gas Used): %d\n", block.GasUsed())
	fmt.Printf("区块大小(Size): %d bytes\n", block.Size())
	fmt.Printf("Nonce: %#x\n", block.Nonce())
	fmt.Printf("状态根(State Root): %s\n", block.Root().Hex())
	fmt.Printf("收据根(Receipt Root): %s\n", block.ReceiptHash().Hex())
	fmt.Println("=============================")
	//输出查询结果到控制台

	//account := common.HexToAddress("0x88aB3E10129Bf18922621e6C2bA36c3dFc95CE78")
	//balance, err := client.BalanceAt(context.Background(), account, nil)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(balance)
	//blockNumber := big.NewInt(5532993)
	//balanceAt, err := client.BalanceAt(context.Background(), account, blockNumber)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(balanceAt) // 25729324269165216042
	//fbalance := new(big.Float)
	//fbalance.SetString(balanceAt.String())
	//ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	//fmt.Println(ethValue) // 25.729324269165216041
	//pendingBalance, err := client.PendingBalanceAt(context.Background(), account)
	//fmt.Println(pendingBalance) // 25729324269165216042

}
