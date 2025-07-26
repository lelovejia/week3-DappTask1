package main

import (
	"DAppGo/counter"
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"log"
	"math/big"
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
	privatekey := viper.GetString("privateKey")

	// 连接到Sepolia测试网络
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + KEY)
	if err != nil {
		log.Fatal(err)
	}

	// 替换为您的私钥（仅用于测试，不要在主网使用真实私钥）
	privateKey, err := crypto.HexToECDSA(privatekey)
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // 0 ether
	auth.GasLimit = uint64(300000) // 设置合理的gas limit
	auth.GasPrice = gasPrice

	//// 部署合约（可选，如果合约已部署可以跳过）
	//address, tx1, _, err := counter.DeployCounter(auth, client)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println("Contract deployed at:", address.Hex())
	//fmt.Println("Transaction hash:", tx1.Hash().Hex())
	instance, err := counter.NewCounter(common.HexToAddress("0xd2c7Ad2a6707D26F6751a8Fb383F9C63590D3ce5"), client)
	if err != nil {
		log.Fatal(err)
	}
	// 调用getCount方法
	count, err := instance.GetCount(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current count:", count)

	// 调用increment方法
	tx, err := instance.Increment(auth)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Increment transaction hash:", tx.Hash().Hex())

	// 等待交易确认
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Fatal(err)
	}
	if receipt.Status != 1 {
		log.Fatal("Transaction failed")
	}

	// 再次获取计数
	count, err = instance.GetCount(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("New count:", count)
}
