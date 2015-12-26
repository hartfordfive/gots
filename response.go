package main

type SecretResponse struct {
	CustId             string `json:"custid,omitempty"`
	Value              string `json:"value,omitempty"`
	MetadataKey        string `json:"metadata_key,omitempty"`
	SecretKey          string `json:"secret_key,omitempty"`
	Ttl                int    `json:"ttl,omitempty"`
	MetadataTtl        int    `json:"metadata_ttl,omitempty"`
	SecretTtl          int    `json:"secret_ttl,omitempty"`
	Recipient          string `json:"recipient,omitempty"`
	Created            int    `json:"created,omitempty"`
	Updated            int    `json:"updated,omitempty"`
	PassphraseRequired bool   `json:"passphrase_required,omitempty"`
	ApiStatus          int    `json:"status,omitempty"`
	HttpRespCode       int    `json:"",omitempty"`
}
