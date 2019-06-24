---
title: Authentication
---

## Authentication

### 歴史
* ユーザ認証
    * システムごとに個別に認証してサービスを利用
* アイデンティティ連携/Web SSO
    * 同一ドメインで Cookie ベースのシンプルな SSO

### 用語
* Identity Provider(IdP)
    * 認証情報を提供する側
* Service Provider(SP)
    * 認証情報を利用する側
* Pull Model, Push Model, Third-Party Security Service
    * SSO の3つのモデル
* claim(クレーム)
    * A piece of information asserted about a subject. 
        * ある主体に関して表明されるひとまとまりの情報
    * A claim is represented as a name/value pair consisting of a Claim Name and a Claim Value.

## 認証・認可のプロトコル
* SAML
    * 認証・認可         
* OAuth
    * 認可
* OpenID Connect
    * 認証

## OAuth 2.0
### RFC6490
* p.6  The interaction between the authorization server and resource server is beyond the scope of this specification.
    * らしい
* p.10 The access token provides an abstraction layer, replacing different authorization constructs (e.g., username and password) with a single
 token understood by the resource server. This abstraction enables issuing access tokens more restrictive than the authorization grant
 used to obtain them, as well as removing the resource server’s need to understand a wide range of authentication methods.    
* p.13 This framework was designed with the clear expectation that future work will define prescriptive profiles and extensions necessary to achieve full web-scale interoperability.
    * これもらしい
* p.16. The authorization server MAY establish a client authentication method with public clients. However, the authorization server MUST NOT rely
 on public client authentication for the purpose of identifying the client.
    * どゆこと？ 
    
## SAML
Defining and maintaining a standard, XML-based framework for creating and exchanging security information between online partners.

### SAML metadata
* 署名に関連する公開鍵の共有とか

### trouble shooting
* [トラブルシューティングのためにブラウザで SAML レスポンスを表示する方法](https://docs.aws.amazon.com/ja_jp/IAM/latest/UserGuide/troubleshoot_saml_view-saml-response.html)

## OAuth
* OAuth のユースケースの多様性は、クライアントの多様性から多くが発生
    * [](https://www.buildinsider.net/enterprise/openid/oauth20)
* クライアントユースケース一覧
    * Web アプリ
        * secret は秘匿可能

## JWT/JWS/JWE/JWA
### JWS
* 2つの serialization
    * JWS Compact Serialization or JWS JSON Serialization
* JOSE Header の members は、 JWS Protected Header/JWS Unprotected Header の members の Union
* JWS の input(The input to the digital signature or MAC computation.)
    * ASCII(BASE64URL(UTF8(JWS Protected Header)) || ’.’ || BASE64URL(JWS Payload))
        * ASCIII は全体にかかっていることに注意
    * 前段と後段の微妙な違いはなぜか？
        * JWS Payload は「任意のバイト列」である一方で、Header は Signature に関連するパラメータの記述であり、Unicode 文字列が想定されている（はず）
            * Payload は JSON も仮定されていない
        * なので、それをまず UTF-8 に変換するということではないか

## Reference
### RFC
* [RFC6749: The OAuth 2.0 Authorization Framework](https://tools.ietf.org/html/rfc6749)
* [EFC6750: The OAuth 2.0 Authorization Framework: Bearer Token Usage](https://tools.ietf.org/html/rfc6750.pdf)
* [RFC7636: Proof Key for Code Exchange by OAuth Public Clients](https://tools.ietf.org/html/rfc7636)
* [RFC7515: JSON Web Signature (JWS)](https://tools.ietf.org/html/rfc7515)
* [RFC7516: JSON Web Encryption (JWE)](https://tools.ietf.org/html/rfc7516)
* [RFC7529: JSON Web Token (JWT)](https://tools.ietf.org/html/rfc7519.pdf)
* [OpenID Connect Core 1.0](https://openid.net/specs/openid-connect-core-1_0.html)

### XML Signature
* [XML Signature Syntax and Processing Version 2.0](https://www.w3.org/TR/xmldsig-core2/)
* [XML Signature Best Practices](https://www.w3.org/TR/2013/NOTE-xmldsig-bestpractices-20130411/)

## JWT/JWS/JWE/JWA
* [Introduction to JSON Web Tokens](https://jwt.io/introduction/)
* [JWT token for multiple websites](https://stackoverflow.com/questions/35423800/jwt-token-for-multiple-websites)
* [JWT Debugger](https://jwt.io/#debugger)
    * トークンの中身が簡単に見れて便利

### SAML
* [SAML技術解説](http://xmlconsortium.org/websv/kaisetsu/C10/content.html)
* [SAML認証ができるまで](https://blog.cybozu.io/entry/4224)

### SSO
* [クラウド時代のシングルサインオン](https://www.osstech.co.jp/_media/techinfo/seminar/hbstudy-20110416-sso.pdf)

### API Security
* [API Security vs. Web Application Security Part 1: A Brief History of Web Application Architecture](https://www.levvel.io/blog/api-security-vs-web-application-security)
* [API Security vs. Web Application Security: Part 2](https://www.levvel.io/blog/api-security-vs-web-application-security-part-2)
    * Web app のいくつかの形態の紹介など

https://medium.com/mixi-developers/reducing-ip-address-restrictions-3793890639e0

https://medium.com/@robert.broeckelmann/saml-v2-0-vs-jwt-series-550551f4eb0d

https://www.oasis-open.org/committees/download.php/22553/sstc-saml-tech-overview-2%200-draft-13.pdf