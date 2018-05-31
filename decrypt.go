package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

func init() {
	parser.AddCommand("decrypt",
		"Decrypts encrypted text, returning the plaintext data",
		"Decrypts the encrypted DEK via KMS, decrypts the data with the DEK, "+
			"outputs to file",
		&decryptCommand)
}

type DecryptCommand struct {
	Filepath string `short:"f" long:"filepath" description:"Path of file to get encrypted string from" default:"./cipher.txt"`
	Nonce    string `short:"n" long:"nonce" description:"Nonce for encryption" required:"true"`
}

var decryptCommand DecryptCommand

func (x *DecryptCommand) Execute(args []string) error {
	fmt.Println("Decrypting...")
	dat, err := ioutil.ReadFile(x.Filepath)
	check(err)
	cipherBase64 := string(dat)
	cipherBytes, errb := base64.StdEncoding.DecodeString(cipherBase64)
	check(errb)
	encryptedDek := cipherBytes[len(cipherBytes)-113 : len(cipherBytes)]
	decryptedDek := googleKMSCrypto(encryptedDek, defaultOptions.ProjectID,
		defaultOptions.LocationID, defaultOptions.KeyRingID,
		defaultOptions.CryptoKeyID, false)
	fmt.Printf("%x\n", decryptedDek)
	plainText := cipherText(cipherBytes[0:len(cipherBytes)-113], cipherblock(decryptedDek), []byte(x.Nonce), false)
	fmt.Printf("plaintext: %s\n", plainText)
	return nil
}
