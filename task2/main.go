package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

	//使用 ethclient 连接到 Sepolia 测试网络。
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/" + KEY)
	if err != nil {
		log.Fatal(err)
	}
	//构造一笔简单的以太币转账交易，指定发送方、接收方和转账金额。
	//对交易进行签名，并将签名后的交易发送到网络。
	//输出交易的哈希值

	//1.加载私钥
	privateKey, err := crypto.HexToECDSA(privatekey)
	if err != nil {
		log.Fatal(err)
	}

	//2.获取发送的帐户的公共地址
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	//3.读取我们应该用于帐户交易的随机数
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}

	//4.设置将要转移的 ETH 数量
	value := big.NewInt(1000000000000000) // in wei (0.001 eth)

	//5.ETH 转账的燃气应设上限为“21000”单位
	gasLimit := uint64(21000)

	//6.平均燃气价格
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	//7.将 ETH 发送给谁
	toAddress := common.HexToAddress("0xefE6Caccd0140869810a10f10E801fB9F8890f60")

	//8.生成我们的未签名以太坊事务
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	//9.获取链ID（用于EIP-155签名）
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	//10.使用发件人的私钥对事务进行签名
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	//11.将已签名的事务广播到整个网络
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx sent: %s", signedTx.Hash().Hex())
}
