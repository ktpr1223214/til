---
title: AD
---

## Active Directory
* アカウントの管理を行うディレクトリ・サービス・システム

### 基礎概念
* 認証と認可
    * 認証: ユーザ本人を確認すること
    * 認可: ユーザにリソースへのアクセスを許可するかどうか判断すること
* 認証と認可でやりたいこと
    * 適切な人に適切な権限を与える
* 新たなセキュリティ境界は ID(Identity) である
    * 企業ネットワークの「内側」と「外側」という区別なく、その組織のリソースにアクセスしてくるすべてのユーザ、そのユーザが使用しているデバイスを、
    強力に管理・保護・認証し、アクセス制御を行うことが求められる

### 関連用語
* ワークグループ
* Kerberos 認証
* 

## Azure Ad
* 主な機能は4つ
    * ディレクトリサービス機能
    * ID 管理機能
    * アプリケーションアクセス管理機能
    * オンプレミス Active Directory ドメインとの統合機能

### AD と Azure AD の相違点
* 使用目的の相違
    * オンプレミスのサーバに対するシングルサインオンを実現するのが Active Directory
        * オンプレミスのアカウントやリソースに対する認証基盤 
    * クラウドサービスに対するシングルサインオンを実現するのが Azure AD
        * クラウドのアカウントやリソースに対する認証基盤
    * つまり、Active Directory と Azure AD ではシングルサインオンでアクセスできる範囲が異なる
* プロトコルの相違
    * AD
        * Kerberos プロトコルで認証と認可
        * ディレクトリのアクセスには、LDAP
    * Azure AD
        * SAML/WS-Federation/OpenID Connect/OAuth などのプロトコルで認証と認可
        * ディレクトリのアクセスには、REST ベースの Azure Graph API
            * どちらも HTTP/HTTPS ベース
* ドメイン構成の相違

## Reference
* [認証にまつわるセキュリティの新常識](https://speakerdeck.com/kthrtty/ren-zheng-nimatuwarusekiyuriteifalsexin-chang-shi)
