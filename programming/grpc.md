---
title: gRPC
---

## gRPC

### 概要
* gRPC Remote Procedure Call
* HTTP/2 を用いた RPC フレームワーク

* メリット
    * Binary protocol (HTTP/2)
    * Multiplexing many requests on one connection (HTTP/2)
    * Header compression (HTTP/2)
    * Strongly typed service and message definition (Protobuf)
    * Idiomatic client/server library implementations in many languages

### 使い方手順
* サーバ
    1. proto 定義を書く
    2. proto からコードを生成
    3. サービスを実装
    4. gRPC サーバを起動
        * 実装したサービスを登録して起動
* クライアント
    1. proto からコードを生成
    2. クライアント実装
    
* grpc-gateway
    * grpc-gateway は gRPC で書かれた API を古典的な JSON over HTTP の API に変換・提供するためのミドルウェア
    * このツールはコード生成器として、ある種のリバースプロキシサーバーを生成

### Load balancing の方法
* Proxy か クライアントで対応するかのいずれか

* Proxy
    * Pros
        * No client-side awareness of backend 
        * Works with untrusted clients
    * Cons
        * 
    
* Client
    * Pros
        * High performance because elimination of extra hop
    * Cons
        *         
#### 

    
## Reference
* [gRPC and REST with gRPC in practice](https://speakerdeck.com/kazegusuri/grpc-and-rest-with-grpc-in-practice)
* [gRPC Load Balancing](https://grpc.io/blog/loadbalancing/)
    * よく纏まっている