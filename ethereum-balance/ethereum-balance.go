package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"go-blockchain-expert/src"
	"math/big"
	"strconv"
)

func main() {
	c, err := NewClient("https://bsc-dataseed1.binance.org/")
	if err != nil {
		panic(err)
	}

	var ps = struct {
		token string
		as    []struct {
			alias   string
			address string
		}
	}{
		token: "0x55d398326f99059ff775485246999027b3197955",
		as: []struct {
			alias   string
			address string
		}{
			{
				alias:   "鲁任秀身份证",
				address: "0xA1dB86087c3BEAE6f242a6a4C06373907e5664dF",
			}, {
				alias:   "李寿松身份证",
				address: "0x788776498A796e7901EAeB288651a264166C076f",
			}, {
				alias:   "段丹丹身份证",
				address: "0x62Ccb8B26EABA824fB987Cb243a35Dd0795859BE",
			}, {
				alias:   "段盼盼身份证",
				address: "0x4BA964F9216AaA30761492F3459A70a6f263Bda8",
			}, {
				alias:   "李金梅身份证",
				address: "0xF47890c55e668794547E4aC295d9BD986C190fBd",
			}, {
				alias:   "段长平身份证",
				address: "0x422FdaBcfC0a6f3eeA8420bA8044401E1fE15C8e",
			}, {
				alias:   "谭春花身份证",
				address: "0x9273350B09A11dEa6b7a3A53CD4e2950605293D1",
			}, {
				alias:   "李强身份证",
				address: "0x05A53b212d9538B4C50fE1014204A34f3953f1b9",
			}, {
				alias:   "李强护照",
				address: "0x86356438e4CaF573637C3Ce20696b1758c7018Bc",
			},
		},
	}

	for _, a := range ps.as {
		b, err := getERC20Balance(c, ps.token, a.address)
		if err != nil {
			panic(err)
		}

		f, err := strconv.ParseFloat(b, 64)
		if err != nil {
			panic(err)
		}

		// 格式化为6位小数，自动四舍五入并补零
		fmt.Println(a.alias, fmt.Sprintf("%.6f", f))
	}

	fmt.Println("成功")
}

func NewClient(host string) (*ethclient.Client, error) {
	rpcDial, err := rpc.Dial(host)
	if err != nil {
		return nil, err
	}

	return ethclient.NewClient(rpcDial), nil
}

func getERC20Balance(c *ethclient.Client, token, address string) (string, error) {
	e, err := src.NewIERC20(common.HexToAddress(token), c)
	if err != nil {
		return "", err
	}

	d, err := e.Decimals(nil)
	if err != nil {
		return "", err
	}

	b, err := e.BalanceOf(nil, common.HexToAddress(address))
	if err != nil {
		return "", err
	}

	return formatUnits(b, d), nil
}

// 精度转换函数（核心计算）
func formatUnits(amount *big.Int, decimals uint8) string {
	// 创建精度除数：10^decimals
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil)

	// 将数值转换为浮点数
	amountFloat := new(big.Float).SetInt(amount)
	divisorFloat := new(big.Float).SetInt(divisor)

	// 执行除法运算
	result := new(big.Float).Quo(amountFloat, divisorFloat)

	// 格式化为字符串（保留所有小数位）
	return result.Text('f', int(decimals))
}
