package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	// 大概念，有密碼的話，一定是私要使用。公鑰不需要密碼！
	// 有密碼的話，都是針對 pem.Block 的 Bytes 資料進行加解密！
	// 公鑰是由私鑰產生，ssh公鑰是由公鑰產生
	password := []byte("pass")

	// 讀取現有私鑰
	prv, err := os.ReadFile("private_key.pem")
	// prv, err := os.ReadFile("id_rsa")
	must(err)

	p, _ := pem.Decode(prv)
	if password != nil { // 有密碼的話，要解析一下Bytes
		body, err := x509.DecryptPEMBlock(p, password)
		must(err)
		p.Bytes = body
	}

	// 自行產生公鑰
	// pv, err := x509.ParsePKCS1PrivateKey(p.Bytes) // 轉成私鑰物件
	// must(err)
	// pub := GenPubKeyBytes(&pv.PublicKey) // 取得公鑰 []byte

	// 讀取現有公鑰
	pub, err := os.ReadFile("public_key.pem")
	must(err)

	// 只用ssh套件去解析
	// rp, err := ssh.ParsePrivateKeyWithPassphrase(prv, password)
	// must(err)
	// rpb := rp.PublicKey().(ssh.CryptoPublicKey).CryptoPublicKey().(*rsa.PublicKey)
	// pub := GenPubKeyBytes(rpb)

	// 自行產生公私鑰
	// prv, pub, sshPubKey := GenRsaKey(password)
	// os.WriteFile("private_key.pem", prv, os.ModePerm)
	// os.WriteFile("public_key.pem", pub, os.ModePerm)
	// os.WriteFile("public_key.pub", sshPubKey, os.ModePerm)

	rawData := []byte("Hello TLS")

	// 加密驗證型
	signData := RsaSignWithSha256(rawData, prv, password)
	fmt.Println("簽署 (用私鑰加密) = ", hex.EncodeToString(signData))                  // 私鑰 -> pem.Decode -> x509.Parse -> Private KEY -> Sign Data Using Private KEY
	fmt.Println("驗證 (用公鑰確認) = ", RsaVerySignWithSha256(rawData, signData, pub)) // 公鑰 -> pem.Decode -> x509.Parse -> Public KEY -> Verify Private Key Using Public KEY

	// 非對稱加解密型
	data := RsaEncrypt(rawData, pub)
	fmt.Println("加密 = ", hex.EncodeToString(data))                // 公鑰 ->  加密 RawData -> Data
	fmt.Println("解密 = ", string(RsaDecrypt(data, prv, password))) // Data -> 私鑰解密 -> RawData
}

//签名
func RsaSignWithSha256(data, privateKeyBytes, password []byte) []byte {
	h := sha256.New()
	h.Write(data)
	hashed := h.Sum(nil)
	// 簽名重點 - 資料使用「演算法」SHA256 加密

	block, _ := pem.Decode(privateKeyBytes)
	if block == nil {
		panic(errors.New("private key error"))
	}

	if password != nil {
		by, err := x509.DecryptPEMBlock(block, password)
		must(err)
		block.Bytes = by
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("ParsePKCS8PrivateKey err", err)
		panic(err)
	}

	// 簽名重點 - 私鑰物件簽名
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		fmt.Printf("Error from signing: %s\n", err)
		panic(err)
	}

	return signature
}

//验证
func RsaVerySignWithSha256(data, signData, keyBytes []byte) bool {
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	// 驗證重點 - 公鑰物件驗證
	hashed := sha256.Sum256(data)
	err = rsa.VerifyPKCS1v15(pubKey.(*rsa.PublicKey), crypto.SHA256, hashed[:], signData)
	if err != nil {
		panic(err)
	}
	return true
}

// 公钥加密
func RsaEncrypt(data, keyBytes []byte) []byte {
	//解密pem格式的公钥
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("public key error"))
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)

	//公鑰加密重點
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
	if err != nil {
		panic(err)
	}
	return ciphertext
}

// 私钥解密
func RsaDecrypt(ciphertext, keyBytes, password []byte) []byte {
	if password != nil { // 解析有設定密碼的私鑰
		key, err := ssh.ParseRawPrivateKeyWithPassphrase(keyBytes, password)
		must(err)
		priv := key.(*rsa.PrivateKey)
		data, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
		must(err)
		return data
	}

	//获取私钥
	block, _ := pem.Decode(keyBytes)
	if block == nil {
		panic(errors.New("private key error!"))
	}

	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	// 私鑰解密重點
	data, err := rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
	if err != nil {
		panic(err)
	}
	return data
}

//RSA公钥私钥产生
func GenRsaKey(password []byte) (prvkey, pubkey, sshPubKey []byte) {
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		panic(err)
	}

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if password != nil {
		block, err = x509.EncryptPEMBlock(rand.Reader, block.Type, block.Bytes, password, x509.PEMCipherAES128)
		must(err)
	}

	prvkey = pem.EncodeToMemory(block) // 私鑰 []byte

	publicKey := &privateKey.PublicKey
	pubkey = GenPubKeyBytes(publicKey) // 公鑰 []byte

	publicRsaKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		panic(err)
	}

	sshPubKey = ssh.MarshalAuthorizedKey(publicRsaKey) // ssh公鑰 []byte
	return
}

func GenPubKeyBytes(publicKey *rsa.PublicKey) []byte {
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		panic(err)
	}
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}

	return pem.EncodeToMemory(block)
}
