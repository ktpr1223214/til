---
title: Keycloak
---

## KeyCloak
### setup
* docker
    * 最低限動かすだけ

``` yaml
version: '3'

services:
  keycloak:
    image: jboss/keycloak:6.0.1
    ports:
      - 80:8080
    environment:
      KEYCLOAK_USER: admin
      KEYCLOAK_PASSWORD: password
```

## OpenID Connect/OAuth 2.0
### setup 手順
* oidc-sample などでクライアント追加
	* Client Protocol はもちろん openid-connect
* Settings
	* Access type: credential
* Credentials
	* Secret を控えておく
* Client Scopes から、適当な名前で新規追加
    * With recent keycloak version 4.6.0 the client id is apparently no longer automatically added to the audience field 'aud' of the access token. Therefore even though the login succeeds the client rejects the user. To fix this you need to configure the audience for your clients
        * cf. https://stackoverflow.com/questions/53550321/keycloak-gatekeeper-aud-claim-and-client-id-do-not-match
* Create Protocol Mapper '<client-id>-audience' とか        
    * Name: <name>
    * Choose Mapper type: Audience
    * Included Client Audience: <client-id>
    * Add to access token: on
* Configure client my-app in the "Clients" menu
    * Client Scopes tab in my-app settings
    * Add available client scopes "good-service" to assigned default client scopes    
* コード側を書く

## SAML の検証
### setup 手順
1. docker なりで動かしはじめる
2. clients を追加
    * ClientID: http://localhost:9000/saml/metadata など
    * 認証方式: saml
3. settings を追加
    * Root URL: http://localhost:9000
    * Valid Redirect URIs: http://localhost:9000/*
4. SAML Keys を export
    * Private key の前後には、-----BEGIN PRIVATE KEY-----/-----END PRIVATE KEY-----
    * Certificate の前後には、-----BEGIN CERTIFICATE-----/-----END CERTIFICATE-----
5. Mappers 追加
    * Mapper Type: User Property
    * Name/Property/SAML Attribute Name: username 
    * SAML Attribute NameFormat: Basic
6. User 追加
    * 名前: 適当
    * Credentials から password を設定


The location of the assertion consumer service MAY be determined using metadata (as in [SAMLMeta]).
The identity provider MUST have some means to establish that this location is in fact controlled by the
service provider. A service provider MAY indicate the SAML binding and the specific assertion consumer
service to use in its <AuthnRequest> and the identity provider MUST honor them if it can.

### 挙動
* Login console に遷移したタイミングで cookie に以下が入る
    * AUTH_SESSION_ID
    * KC_RESTART
* ログイン認証成功すると、cookie に token が入る

## エラー
* error=client_not_found
    * ClientID が間違っている可能性がある
* Invalid redirect uri
    * Clients/Settings から Valid Redirect URIs を修正
        * <your app url>/* 辺りで動く
        * https://stackoverflow.com/questions/45352880/keycloak-invalid-parameter-redirect-uri?rq=1        

Openconnect ID/Oatuth/saml 2.0
keycloak 

* tree 構造でユーザ管理
    * 人に id/password が管理できて
* AD にありイコール LDAP
    * AD の方が完成度が高い
* LDAP の BIND
    * authentication は1機能に過ぎない    
    * search も可能
        * 問い合わせができるということ

## SP metadata
* metadata go struct
``` go
type EntityDescriptor struct {
	XMLName                       xml.Name      `xml:"urn:oasis:names:tc:SAML:2.0:metadata EntityDescriptor"`
	EntityID                      string        `xml:"entityID,attr"`
	ID                            string        `xml:",attr,omitempty"`
	ValidUntil                    time.Time     `xml:"validUntil,attr,omitempty"`
	CacheDuration                 time.Duration `xml:"cacheDuration,attr,omitempty"`
	Signature                     *etree.Element
	RoleDescriptors               []RoleDescriptor               `xml:"RoleDescriptor"`
	IDPSSODescriptors             []IDPSSODescriptor             `xml:"IDPSSODescriptor"`
	SPSSODescriptors              []SPSSODescriptor              `xml:"SPSSODescriptor"`
	AuthnAuthorityDescriptors     []AuthnAuthorityDescriptor     `xml:"AuthnAuthorityDescriptor"`
	AttributeAuthorityDescriptors []AttributeAuthorityDescriptor `xml:"AttributeAuthorityDescriptor"`
	PDPDescriptors                []PDPDescriptor                `xml:"PDPDescriptor"`
	AffiliationDescriptor         *AffiliationDescriptor
	Organization                  *Organization
	ContactPerson                 *ContactPerson
	AdditionalMetadataLocations   []string `xml:"AdditionalMetadataLocation"`
}

// ...
type Alias EntityDescriptor
aux := &struct {
    ValidUntil    RelaxedTime `xml:"validUntil,attr,omitempty"`
    CacheDuration Duration    `xml:"cacheDuration,attr,omitempty"`
    *Alias
}{
    ValidUntil:    RelaxedTime(m.ValidUntil),
    CacheDuration: Duration(m.CacheDuration),
    Alias:         (*Alias)(&m),
}	
```

* metadata xml
``` xml
<EntityDescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" validUntil="2019-06-20T08:43:29.378Z" entityID="http://localhost:8000/saml/metadata">
  <SPSSODescriptor xmlns="urn:oasis:names:tc:SAML:2.0:metadata" validUntil="2019-06-20T08:43:29.378087Z" protocolSupportEnumeration="urn:oasis:names:tc:SAML:2.0:protocol" AuthnRequestsSigned="false" WantAssertionsSigned="true">
    <KeyDescriptor use="signing">
      <KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#">
        <X509Data>
          <X509Certificate>MIICuTCCAaECBgFrSjCkJTANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVodHRwOi8vbG9jYWxob3N0OjgwMDAwHhcNMTkwNjEyMDUzNTMwWhcNMjkwNjEyMDUzNzEwWjAgMR4wHAYDVQQDDBVodHRwOi8vbG9jYWxob3N0OjgwMDAwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCKfZ1ihJvJF535wRPub/iy5qYhD2vjuOirDB649b3etcjF8bFYfRJbxfvDs8hw730iKLFSzRXHX9sin2nd4szhgzNum79aA8JMlkMfWp7jPQmSS/u7EZS30xYLIUkEJeN9lYsEg9NipMqISVLBG7T3QsxzTave3eqBtsPcT4K+fVe9h2pDsvmv7uT3Tempy2SzmTkCaghJFtsHZbdEVH2j8GYIREoCZoCEmCILpcdd1ZkhvlWBw9AC5euInOzTt+l+3fVccA/6CArw/+P1Sqo9t+PUxpQ0EUJRPgU/JqIwi55ManFqJfNiKLf9zGLFQxnYaIbH5+7/vqXOf+ijkifBAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAEBHtvp6QIFY2uRxhZJQMzNVqAKyXJ/LgrIgPYe2oGo59LBtbCw/bbl5hyBHoHf8mz1jlp5bMwGSnO25SrvOEJwa6BoAo41rHzCqQtnI93zOv9EfVgkmwdCe7rNORiXWSFdvQi/b8yfDEtl9R2wde9PUnAJENEfFEvATMzocKDkfl9VtXfuYkSfRlh/LSzaRpGzpraNPHNZIC4YbBAJ5UNlkWmwh5+39ZEwnZaeFoFqRnarSYkD9GRnsdqU4t36tjSG+hEUqd3/J+H8xlYrntRo1aK5KsDPonOMwlgcSkLHLJiwAuASWlwuo/igc17iWIYGcm7xYW3l9FcvM2qE6Bp0=</X509Certificate>
        </X509Data>
      </KeyInfo>
    </KeyDescriptor>
    <KeyDescriptor use="encryption">
      <KeyInfo xmlns="http://www.w3.org/2000/09/xmldsig#">
        <X509Data>
          <X509Certificate>MIICuTCCAaECBgFrSjCkJTANBgkqhkiG9w0BAQsFADAgMR4wHAYDVQQDDBVodHRwOi8vbG9jYWxob3N0OjgwMDAwHhcNMTkwNjEyMDUzNTMwWhcNMjkwNjEyMDUzNzEwWjAgMR4wHAYDVQQDDBVodHRwOi8vbG9jYWxob3N0OjgwMDAwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQCKfZ1ihJvJF535wRPub/iy5qYhD2vjuOirDB649b3etcjF8bFYfRJbxfvDs8hw730iKLFSzRXHX9sin2nd4szhgzNum79aA8JMlkMfWp7jPQmSS/u7EZS30xYLIUkEJeN9lYsEg9NipMqISVLBG7T3QsxzTave3eqBtsPcT4K+fVe9h2pDsvmv7uT3Tempy2SzmTkCaghJFtsHZbdEVH2j8GYIREoCZoCEmCILpcdd1ZkhvlWBw9AC5euInOzTt+l+3fVccA/6CArw/+P1Sqo9t+PUxpQ0EUJRPgU/JqIwi55ManFqJfNiKLf9zGLFQxnYaIbH5+7/vqXOf+ijkifBAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAEBHtvp6QIFY2uRxhZJQMzNVqAKyXJ/LgrIgPYe2oGo59LBtbCw/bbl5hyBHoHf8mz1jlp5bMwGSnO25SrvOEJwa6BoAo41rHzCqQtnI93zOv9EfVgkmwdCe7rNORiXWSFdvQi/b8yfDEtl9R2wde9PUnAJENEfFEvATMzocKDkfl9VtXfuYkSfRlh/LSzaRpGzpraNPHNZIC4YbBAJ5UNlkWmwh5+39ZEwnZaeFoFqRnarSYkD9GRnsdqU4t36tjSG+hEUqd3/J+H8xlYrntRo1aK5KsDPonOMwlgcSkLHLJiwAuASWlwuo/igc17iWIYGcm7xYW3l9FcvM2qE6Bp0=</X509Certificate>
        </X509Data>
      </KeyInfo>
      <EncryptionMethod Algorithm="http://www.w3.org/2001/04/xmlenc#aes128-cbc"></EncryptionMethod>
      <EncryptionMethod Algorithm="http://www.w3.org/2001/04/xmlenc#aes192-cbc"></EncryptionMethod>
      <EncryptionMethod Algorithm="http://www.w3.org/2001/04/xmlenc#aes256-cbc"></EncryptionMethod>
      <EncryptionMethod Algorithm="http://www.w3.org/2001/04/xmlenc#rsa-oaep-mgf1p"></EncryptionMethod>
    </KeyDescriptor>
    <AssertionConsumerService Binding="urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST" Location="http://localhost:8000/saml/acs" index="1"></AssertionConsumerService>
  </SPSSODescriptor>
</EntityDescriptor>
```

## 
* 認証周りの話を進めていくというのが先で
    * 監視基盤もユーザ認証があったら便利なはずで     
* proxy の話
    * ある A のときは、SSO でするが、別の場合には Azure AD は使わないみたいな
        * Azure AD は R の社員のみなので
    * 純粋な SSO でなくても、認証 backend が共通で使えるのかどうか
    * SSO のプロトコルは SAML2.0/OAuth
* 今後アプリ側で実装していく場合に楽そうなもの        
    * web app の人たちは Oauth が楽かも？
        * facebook とかでの認証があるので
* そこ楽なやり方で
    * 統一でもなんでも
* アプリケーション登録も Azure AD ではなく中間で見ても良いかも
* Oauth のトークンは、アプリーProxy 間だけで
    * SSO サーバの代理
    * Azure AD が完全にインターフェースとして抽象化される感じになる
        * アプリケーション共通基盤的なものが作りやすくなるかも
* 利用者申請とか、権限も中間でやるようにすると、認証は Azure AD で、認可は中間でみたいな世界
    * 統括の雅子さんで対応できるとか、承認がウェブベースで飛んでプロダクトの管理者がオッケーって言えば使えるみたいな

* アプリA/B 作ってみてちゃんと動くか検証
* 付替え的なことができるのか
* 登録するところを抽象化するのかどうか
    * B の場合だと内部的なログインページにリダイレクトなりして、IdP として振る舞う
    * それが無理なケースに、Azure 側に行ってリダイレクト的な
* 真ん中にユーザ登録があれば、
* 一旦は proxy で受けて、さらにそれをアプリ側に返すということが出来たほうが自由度は上がる    

* SAML の方が良いかも？
    * OpenID だと
    
* プロトコルの決定は
    * 中間層的な扱いしやすさ
    * アプリ的な嬉しさ
        * これは二段階でも良いかも    

1. プロトコル決定
2. プロトコル読み替えが出来るのかどうか

## 個人的メモ
* azure ad に連携アプリ増やすみたいなことは簡単にできるんか？
* できれば複数アプリで検証(同じでなく)
    * 純粋な auth（認証）
    * a で auth で b 実行検証（認証）
    * アプリ内部での認可区分がうまくいくか(admin/user とか)
    * a から b の REST API をたたくとか
    * js のもの

## Reference
* [KeycloakでOpenID Connectを使ってみる（Spring Bootアプリケーション編）](https://qiita.com/tamura__246/items/4290a1035e1adcab733b)
* [Spring Securityを使って、KeycloakでOpenID Connect](https://kazuhira-r.hatenablog.com/entry/20180902/1535873954)
* [Spring Boot 1/2のアプリにKeycloakのOpenID Connectを使ってシングルサインオン](https://blog.ik.am/entries/445)