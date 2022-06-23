package user

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/DataWorkbench/account/config"
	"github.com/DataWorkbench/common/qerror"
	"github.com/go-ldap/ldap/v3"
	"io/ioutil"
)

var LdapProvider *LdapAuth

type LdapAuth struct {
	ldapConf *config.LdapConfig
}

func InitLdap(config *config.LdapConfig) {
	LdapProvider = &LdapAuth{
		ldapConf: config,
	}
}

func (l *LdapAuth) getConn() (*ldap.Conn, error) {
	if l.ldapConf.InsecureSkipVerify {
		conn, err := ldap.DialURL(l.ldapConf.Url, ldap.DialWithTLSConfig(&tls.Config{InsecureSkipVerify: l.ldapConf.InsecureSkipVerify}))
		if err != nil {
			return nil, err
		}
		return conn, nil
	}
	tlsConfig := tls.Config{}
	tlsConfig.RootCAs = x509.NewCertPool()
	var caCert []byte
	var err error
	// Load CA cert
	if l.ldapConf.RootCA != "" {
		if caCert, err = ioutil.ReadFile(l.ldapConf.RootCA); err != nil {
			return nil, err
		}
	}
	if l.ldapConf.RootCAData != "" {
		if caCert, err = base64.StdEncoding.DecodeString(l.ldapConf.RootCAData); err != nil {
			return nil, err
		}
	}
	if caCert != nil {
		tlsConfig.RootCAs.AppendCertsFromPEM(caCert)
	}
	return ldap.DialURL("tcp", ldap.DialWithTLSConfig(&tlsConfig))

}

func (l *LdapAuth) Authentication(username, password string) (map[string]interface{}, error) {

	m := make(map[string]interface{})

	conn, err := l.getConn()
	if err != nil {
		return m, err
	}
	verify := &tls.Config{InsecureSkipVerify: l.ldapConf.InsecureSkipVerify}
	// Reconnect with TLS
	if l.ldapConf.StartTLS {
		err := conn.StartTLS(verify)
		if err != nil {
			return m, err
		}
	}
	defer conn.Close()
	// First bind with a read only user
	err = conn.Bind(l.ldapConf.ManagerDN, l.ldapConf.ManagerPassword)
	if err != nil {
		return m, err
	}
	// Search for the given username
	searchRequest := ldap.NewSearchRequest(
		l.ldapConf.UserSearchBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=organizationalPerson)(%s=%s))", l.ldapConf.LoginAttribute, username),
		[]string{},
		nil,
	)
	sr, err := conn.Search(searchRequest)
	if err != nil {
		return m, err
	}
	if len(sr.Entries) != 1 {
		return m, qerror.UserNotExists.Format(username)
	}
	userdn := sr.Entries[0].DN
	mail := sr.Entries[0].GetAttributeValue(l.ldapConf.MailAttribute)
	m["mail"] = mail
	// Bind as the user to verify their password
	err = conn.Bind(userdn, password)
	if err != nil {
		return m, qerror.UserNameOrPasswordError
	}
	return m, nil
}
