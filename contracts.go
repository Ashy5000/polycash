package main

type ContractParty struct {
	Signature Signature
	PublicKey PublicKey
}

type Contract struct {
	Contents string
	Parties  []ContractParty
}
