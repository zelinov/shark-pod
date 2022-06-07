# Sharkpod auth service

## TODO

- Subscribe to active user's transactions  
if user sold his sharks token, and he no longer has the required sharks amount -> logout
- 2022-04-28T12:09:31.727Z	ERROR	repository/tokenMongoRepository.go:83	bulk write exception: write errors: [E11000 duplicate key error collection: sharkdb.tokenMetadata index: public_key_1 dup key: { public_key: "8Zh9g8SUpdhvyjLQADkJHSv2j5SYsULVt52uYPM5DXzB" }]
  sharks/adapters/outbound/repository.(*TokenMongoRepository).SaveAll
  /var/www/sharkpod/back/sharkpod-backend/adapters/outbound/repository/tokenMongoRepository.go:83

## API

<details>
<summary><span style="font-size:1.3em;"><span style="color:MediumSeaGreen"><strong>POST</strong></span> /api/v1/login</span></summary>

**Request body**

```json
{
  "nonce": "${encrypted nonce by Solana}",
  "publicKey": "${base64 user public key}"
}
```

**Response 200**

```json
{
  "accessToken": "String",
  "refreshToken": "String",
  "publicKey": "String",
  "expiredAt": "Date in format 2022-04-07T00:48:48.073887+04:00"
}
```

**Response 401**

Header  
>Error-Code: INSUFFICIENT_NFT  

Body  
>less than 1 sharks tokens
</details>

<details>
<summary><span style="font-size:1.3em;"><span style="color:MediumSeaGreen"><strong>POST</strong></span> /api/v1/refresh</span></summary>

**Request body**

```json
{
  "refresh": "${Your refresh token from login response}"
}
```

**Response 200**

```json
{
  "accessToken": "String",
  "refreshToken": "String",
  "publicKey": "String",
  "expiredAt": "Date in format 2022-04-07T00:48:48.073887+04:00"
}
```

**Response 401**

Header
>Error-Code: INVALID_TOKEN

Body
>Token is invalid
</details>

<details>
<summary><span style="font-size:1.3em;"><span style="color:DodgerBlue"><strong>GET</strong></span> /api/v1/nonce/${publicKey}</span></summary>

**Request**

publicKey - Base64 user public key

**Response 200**

```json
{
  "nonce": "UUID string"
}
```
</details>

<details>
<summary><span style="font-size:1.3em;"><span style="color:MediumSeaGreen"><strong>POST</strong></span> /api/v1/check</span></summary>

**Request**

Header
>"Authorization": "Bearer ${accessToken}"

**Response 200**

```json
{
  "publicKey": "String"
}
```

**Response 401**

Header
>Error-Code: INVALID_TOKEN

Body
>Token is invalid
</details>

<details>
<summary><span style="font-size:1.3em;"><span style="color:MediumSeaGreen"><strong>POST</strong></span> /api/v1/logout</span></summary>

**Request**

Header
>"Authorization": "Bearer ${accessToken}"

**Response 200**

Empty body

**Response 401**

Header
>Error-Code: INVALID_TOKEN

Body
>Token is invalid
</details>

<details>
<summary><span style="font-size:1.3em;"><span style="color:DodgerBlue"><strong>GET</strong></span> /api/v1/tokens</span></summary>

**Request**

Header
>"Authorization": "Bearer ${accessToken}"

**Response 200**

```json
{
  "images": [
    {
      "url": "String"
    },
    {
      "url": "String"
    }
  ]
}
```
</details>

### Error codes

Header **Error-Code**

| Code              | Description                                                   |
|-------------------|---------------------------------------------------------------|
| INSUFFICIENT_NFT  | There are not enough non-fungible tokens in the user's wallet |  
| INVALID_TOKEN     | The jwt token is invalid                                      |  
| INVALID_SIGNATURE | The solana signature is invalid                               |  
| UNKNOWN           | Another errors                                                |  

