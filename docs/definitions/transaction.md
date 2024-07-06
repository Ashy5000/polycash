### Transaction (semi-formal definition)

$$
tx := \{
  sender(pubKey)\
  recipient(pubKey)\
  amount(float)\
  sig(dilithium3Signature)\
  timestamp(time)\
  contracts(array(contract))\
  fromSmartContract(bool)\
  body(rawBytes)\
  bodySignatures(array(signature))
\}
$$

$$
txHashDelim := ':'
$$

$$
tx::string() \to string \{
  concat(sender::string(), txHashDelim, recipient::string(), txHashDelim, amount::string(), timestamp::string())
\}
$$

$$
tx::hash () \to hash256 \{
  sha256(tx::string())
\}
$$
