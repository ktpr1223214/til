package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-ldap/ldap"
	"github.com/pkg/errors"
)

var (
	// ErrUserNotFound ユーザーが見つからなかったエラー
	ErrUserNotFound = errors.New("user not found")
)

// Client はLDAPサーバーへの接続情報を持っている
type Client struct {
	host   string
	baseDN string
}

// SearchResult LDAP サーバーへの検索結果
type SearchResult struct {
	Name        string
	DN          string
	Description string
}

// BindUser Bind を行うユーザー情報
type BindUser struct {
	Name     string
	Password string
}

// BindConn LDAP Bind connection
type BindConn struct {
	client *Client
	conn   *ldap.Conn
}

// Close Close connection
func (b *BindConn) Close() {
	b.conn.Close()
}

// NewClient LDAP 設定作成
func NewClient(host string, baseDN string) *Client {
	return &Client{
		host:   host,
		baseDN: baseDN,
	}
}

// Bind Bind LDAP User
func (c *Client) Bind(user *BindUser) (*BindConn, error) {
	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", c.host, 389))
	if err != nil {
		return nil, errors.Wrap(err, "failed to dial LDAP server")
	}
	if err := conn.Bind(user.Name, user.Password); err != nil {
		return nil, errors.Wrapf(err, "failed to bind User %s", user.Name)
	}

	return &BindConn{
		client: c,
		conn:   conn,
	}, nil
}

// Search 検索
// TODO: 引数
func (b *BindConn) Search() ([]*SearchResult, error) {
	sr, err := b.conn.Search(ldap.NewSearchRequest(
		b.client.baseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		// cf. https://confluence.atlassian.com/kb/how-to-write-ldap-search-filters-792496933.html
		"(objectClass=simpleSecurityObject)",
		[]string{"dn", "cm", "description"},
		nil,
	))
	if err != nil {
		return nil, errors.Wrap(err, "failed to search LDAP")
	}
	if len(sr.Entries) == 0 {
		return nil, ErrUserNotFound
	}

	var searchResults []*SearchResult
	for _, entry := range sr.Entries {
		searchResults = append(searchResults, &SearchResult{
			Name:        entry.GetAttributeValue("displayName"),
			Description: entry.GetAttributeValue("description"),
			DN:          entry.GetAttributeValue("dn"),
		})
	}
	return searchResults, nil
}

func main() {
	client := NewClient("localhost", "cn=admin,dc=sample,dc=com")
	// client := NewClient("localhost", "ou=people,dc=sample,dc=com")
	bind, err := client.Bind(&BindUser{
		Name:     "cn=admin,dc=sample,dc=com",
		Password: os.Getenv("PASSWORD"),
	})
	if err != nil {
		log.Printf("Failed to bind user %s", err)
		os.Exit(1)
	}

	results, err := bind.Search()
	if err != nil {
		log.Printf("Failed to search user %s", err)
		os.Exit(1)
	}
	for _, res := range results {
		log.Println(res)
	}

	bind.Close()
}
