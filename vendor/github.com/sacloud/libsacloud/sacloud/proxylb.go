// Copyright 2016-2020 The Libsacloud Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sacloud

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ProxyLB ProxyLB(CommonServiceItem)
type ProxyLB struct {
	*Resource        // ID
	propName         // 名称
	propDescription  // 説明
	propServiceClass // サービスクラス
	propIcon         // アイコン
	propTags         // タグ
	propCreatedAt    // 作成日時
	propModifiedAt   // 変更日時
	propAvailability // 有効状態

	Status   *ProxyLBStatus  `json:",omitempty"` // ステータス
	Provider ProxyLBProvider `json:",omitempty"` // プロバイダ
	Settings ProxyLBSettings `json:",omitempty"` // ProxyLB設定

}

// ProxyLBSettings ProxyLB設定
type ProxyLBSettings struct {
	ProxyLB ProxyLBSetting `json:",omitempty"` // ProxyLB ProxyLBエントリー
}

// ProxyLBStatus ProxyLBステータス
type ProxyLBStatus struct {
	FQDN             string   `json:",omitempty"` // 割り当てられたFQDN(site-*******.proxylb?.sakura.ne.jp) UseVIPFailoverがtrueの場合のみ有効
	VirtualIPAddress string   `json:",omitempty"` // 割り当てられたVIP UseVIPFailoverがfalseの場合のみ有効
	ProxyNetworks    []string `json:",omitempty"` // プロキシ元ネットワークアドレス(CIDR)
	UseVIPFailover   bool     // VIPフェイルオーバ
}

// ProxyLBProvider プロバイダ
type ProxyLBProvider struct {
	Class string `json:",omitempty"` // クラス
}

// CreateNewProxyLB ProxyLB作成
func CreateNewProxyLB(name string) *ProxyLB {
	return &ProxyLB{
		Resource: &Resource{},
		propName: propName{Name: name},
		Provider: ProxyLBProvider{
			Class: "proxylb",
		},
		Settings: ProxyLBSettings{
			ProxyLB: ProxyLBSetting{
				HealthCheck: defaultProxyLBHealthCheck,
				SorryServer: ProxyLBSorryServer{},
				Servers:     []ProxyLBServer{},
				//LetsEncrypt:   ProxyLBACMESetting{},
				StickySession: ProxyLBSessionSetting{},
			},
		},
	}
}

// ProxyLBPlan ProxyLBプラン
type ProxyLBPlan int

var (
	// ProxyLBPlan100 100cpsプラン
	ProxyLBPlan100 = ProxyLBPlan(100)
	// ProxyLBPlan500 500cpsプラン
	ProxyLBPlan500 = ProxyLBPlan(500)
	// ProxyLBPlan1000 1,000cpsプラン
	ProxyLBPlan1000 = ProxyLBPlan(1000)
	// ProxyLBPlan5000 5,000cpsプラン
	ProxyLBPlan5000 = ProxyLBPlan(5000)
	// ProxyLBPlan10000 10,000cpsプラン
	ProxyLBPlan10000 = ProxyLBPlan(10000)
	// ProxyLBPlan50000 50,000cpsプラン
	ProxyLBPlan50000 = ProxyLBPlan(50000)
	// ProxyLBPlan100000 100,000cpsプラン
	ProxyLBPlan100000 = ProxyLBPlan(100000)
	// ProxyLBPlan400000 400,000cpsプラン
	ProxyLBPlan400000 = ProxyLBPlan(400000)
)

// AllowProxyLBPlans 有効なプランIDリスト
var AllowProxyLBPlans = []int{
	int(ProxyLBPlan100),
	int(ProxyLBPlan500),
	int(ProxyLBPlan1000),
	int(ProxyLBPlan5000),
	int(ProxyLBPlan10000),
	int(ProxyLBPlan50000),
	int(ProxyLBPlan100000),
	int(ProxyLBPlan400000),
}

// GetPlan プラン取得(デフォルト: 1000cps)
func (p *ProxyLB) GetPlan() ProxyLBPlan {
	classes := strings.Split(p.ServiceClass, "/")
	class, err := strconv.Atoi(classes[len(classes)-1])
	if err != nil {
		return ProxyLBPlan1000
	}
	return ProxyLBPlan(class)
}

// SetPlan プラン指定
func (p *ProxyLB) SetPlan(plan ProxyLBPlan) {
	p.ServiceClass = fmt.Sprintf("cloud/proxylb/plain/%d", plan)
}

// SetHTTPHealthCheck HTTPヘルスチェック 設定
func (p *ProxyLB) SetHTTPHealthCheck(hostHeader, path string, delayLoop int) {
	if delayLoop <= 0 {
		delayLoop = 10
	}

	p.Settings.ProxyLB.HealthCheck.Protocol = "http"
	p.Settings.ProxyLB.HealthCheck.Host = hostHeader
	p.Settings.ProxyLB.HealthCheck.Path = path
	p.Settings.ProxyLB.HealthCheck.DelayLoop = delayLoop
}

// SetTCPHealthCheck TCPヘルスチェック 設定
func (p *ProxyLB) SetTCPHealthCheck(delayLoop int) {
	if delayLoop <= 0 {
		delayLoop = 10
	}

	p.Settings.ProxyLB.HealthCheck.Protocol = "tcp"
	p.Settings.ProxyLB.HealthCheck.Host = ""
	p.Settings.ProxyLB.HealthCheck.Path = ""
	p.Settings.ProxyLB.HealthCheck.DelayLoop = delayLoop
}

// SetSorryServer ソーリーサーバ 設定
func (p *ProxyLB) SetSorryServer(ipaddress string, port int) {
	var pt *int
	if port > 0 {
		pt = &port
	}
	p.Settings.ProxyLB.SorryServer = ProxyLBSorryServer{
		IPAddress: ipaddress,
		Port:      pt,
	}
}

// ClearSorryServer ソーリーサーバ クリア
func (p *ProxyLB) ClearSorryServer() {
	p.SetSorryServer("", 0)
}

// HasProxyLBServer ProxyLB配下にサーバーを保持しているか判定
func (p *ProxyLB) HasProxyLBServer() bool {
	return len(p.Settings.ProxyLB.Servers) > 0
}

// ClearProxyLBServer ProxyLB配下のサーバーをクリア
func (p *ProxyLB) ClearProxyLBServer() {
	p.Settings.ProxyLB.Servers = []ProxyLBServer{}
}

// AddBindPort バインドポート追加
func (p *ProxyLB) AddBindPort(mode string, port int, redirectToHTTPS, supportHTTP2 bool, addResponseHeader []*ProxyLBResponseHeader) {
	p.Settings.ProxyLB.AddBindPort(mode, port, redirectToHTTPS, supportHTTP2, addResponseHeader)
}

// DeleteBindPort バインドポート削除
func (p *ProxyLB) DeleteBindPort(mode string, port int) {
	p.Settings.ProxyLB.DeleteBindPort(mode, port)
}

// ClearBindPorts バインドポート クリア
func (p *ProxyLB) ClearBindPorts() {
	p.Settings.ProxyLB.BindPorts = []*ProxyLBBindPorts{}
}

// AddServer ProxyLB配下のサーバーを追加
func (p *ProxyLB) AddServer(ip string, port int, enabled bool, serverGroup string) {
	p.Settings.ProxyLB.AddServer(ip, port, enabled, serverGroup)
}

// DeleteServer ProxyLB配下のサーバーを削除
func (p *ProxyLB) DeleteServer(ip string, port int) {
	p.Settings.ProxyLB.DeleteServer(ip, port)
}

// ProxyLBSetting ProxyLBセッティング
type ProxyLBSetting struct {
	HealthCheck   ProxyLBHealthCheck  // ヘルスチェック
	SorryServer   ProxyLBSorryServer  // ソーリーサーバー
	BindPorts     []*ProxyLBBindPorts // プロキシ方式(プロトコル&ポート)
	Servers       []ProxyLBServer     // サーバー
	Rules         []ProxyLBRule       // 振り分けルール
	LetsEncrypt   *ProxyLBACMESetting `json:",omitempty"` // Let's encryptでの証明書取得設定
	StickySession ProxyLBSessionSetting
	Timeout       *ProxyLBTimeout `json:",omitempty"` // タイムアウト
}

// ProxyLBSorryServer ソーリーサーバ
type ProxyLBSorryServer struct {
	IPAddress string // IPアドレス
	Port      *int   // ポート
}

// AddBindPort バインドポート追加
func (s *ProxyLBSetting) AddBindPort(mode string, port int, redirectToHTTPS, supportHTTP2 bool, addResponseHeader []*ProxyLBResponseHeader) {
	var isExist bool
	for i := range s.BindPorts {
		if s.BindPorts[i].ProxyMode == mode && s.BindPorts[i].Port == port {
			isExist = true
		}
	}

	if !isExist {
		s.BindPorts = append(s.BindPorts, &ProxyLBBindPorts{
			ProxyMode:         mode,
			Port:              port,
			RedirectToHTTPS:   redirectToHTTPS,
			SupportHTTP2:      supportHTTP2,
			AddResponseHeader: addResponseHeader,
		})
	}
}

// DeleteBindPort バインドポート削除
func (s *ProxyLBSetting) DeleteBindPort(mode string, port int) {
	var res []*ProxyLBBindPorts
	for i := range s.BindPorts {
		if s.BindPorts[i].ProxyMode != mode || s.BindPorts[i].Port != port {
			res = append(res, s.BindPorts[i])
		}
	}
	s.BindPorts = res
}

// AddServer ProxyLB配下のサーバーを追加
func (s *ProxyLBSetting) AddServer(ip string, port int, enabled bool, serverGroup string) {
	var record ProxyLBServer
	var isExist = false
	for i := range s.Servers {
		if s.Servers[i].IPAddress == ip && s.Servers[i].Port == port {
			isExist = true
			s.Servers[i].Enabled = enabled
		}
	}

	if !isExist {
		record = ProxyLBServer{
			IPAddress:   ip,
			Port:        port,
			Enabled:     enabled,
			ServerGroup: serverGroup,
		}
		s.Servers = append(s.Servers, record)
	}
}

// DeleteServer ProxyLB配下のサーバーを削除
func (s *ProxyLBSetting) DeleteServer(ip string, port int) {
	var res []ProxyLBServer
	for i := range s.Servers {
		if s.Servers[i].IPAddress != ip || s.Servers[i].Port != port {
			res = append(res, s.Servers[i])
		}
	}

	s.Servers = res
}

// AllowProxyLBBindModes プロキシ方式
var AllowProxyLBBindModes = []string{"http", "https", "tcp"}

// ProxyLBBindPorts プロキシ方式
type ProxyLBBindPorts struct {
	ProxyMode         string                   `json:",omitempty"`      // モード(プロトコル)
	Port              int                      `json:",omitempty"`      // ポート
	RedirectToHTTPS   bool                     `json:"RedirectToHttps"` // HTTPSへのリダイレクト(モードがhttpの場合のみ)
	SupportHTTP2      bool                     `json:"SupportHttp2"`    // HTTP/2のサポート(モードがhttpsの場合のみ)
	AddResponseHeader []*ProxyLBResponseHeader `json:",omitempty"`      // レスポンスヘッダ
}

// ProxyLBResponseHeader ポートごとの追加レスポンスヘッダ
type ProxyLBResponseHeader struct {
	Header string // ヘッダ名称(英字, 数字, ハイフン)
	Value  string // 値(英字, 数字, 半角スペース, 一部記号（!#$%&'()*+,-./:;<=>?@[]^_`{|}~）)
}

// ProxyLBServer ProxyLB配下のサーバー
type ProxyLBServer struct {
	IPAddress   string `json:",omitempty"` // IPアドレス
	Port        int    `json:",omitempty"` // ポート
	ServerGroup string // サーバグループ
	Enabled     bool   // 有効/無効
}

// NewProxyLBServer ProxyLB配下のサーバ作成
func NewProxyLBServer(ipaddress string, port int) *ProxyLBServer {
	return &ProxyLBServer{
		IPAddress: ipaddress,
		Port:      port,
		Enabled:   true,
	}
}

// ProxyLBRule ProxyLBの振り分けルール
type ProxyLBRule struct {
	Host        string `json:",omitempty"` // ホストヘッダのパターン(ワイルドカードとして?と*が利用可能)
	Path        string `json:",omitempty"` // パス
	ServerGroup string
}

// ProxyLBACMESetting Let's Encryptでの証明書取得設定
type ProxyLBACMESetting struct {
	Enabled    bool
	CommonName string `json:",omitempty"`
}

// ProxyLBSessionSetting セッション維持機能設定
type ProxyLBSessionSetting struct {
	Enabled bool
	Method  string `json:",omitempty"`
}

// ProxyLBTimeout 実サーバの通信タイムアウト
type ProxyLBTimeout struct {
	InactiveSec int `json:",omitempty"` // 10から600まで1秒刻みで設定可
}

// ProxyLBStickySessionDefaultMethod セッション維持のデフォルトメソッド(クッキー)
const ProxyLBStickySessionDefaultMethod = "cookie"

// AllowProxyLBHealthCheckProtocols プロキシLBで利用できるヘルスチェックプロトコル
var AllowProxyLBHealthCheckProtocols = []string{"http", "tcp"}

// ProxyLBHealthCheck ヘルスチェック
type ProxyLBHealthCheck struct {
	Protocol  string `json:",omitempty"` // プロトコル
	Host      string `json:",omitempty"` // 対象ホスト
	Path      string `json:",omitempty"` // HTTPの場合のリクエストパス
	DelayLoop int    `json:",omitempty"` // 監視間隔

}

var defaultProxyLBHealthCheck = ProxyLBHealthCheck{
	Protocol:  "http",
	Host:      "",
	Path:      "/",
	DelayLoop: 10,
}

// ProxyLBAdditionalCerts additional certificates
type ProxyLBAdditionalCerts []*ProxyLBCertificate

// ProxyLBCertificates ProxyLBのSSL証明書
type ProxyLBCertificates struct {
	PrimaryCert     *ProxyLBCertificate
	AdditionalCerts ProxyLBAdditionalCerts
}

// UnmarshalJSON UnmarshalJSON(AdditionalCertsが空の場合に空文字を返す問題への対応)
func (p *ProxyLBAdditionalCerts) UnmarshalJSON(data []byte) error {
	targetData := strings.Replace(strings.Replace(strings.Replace(string(data), " ", "", -1), "\n", "", -1), `""`, ``, -1)
	if targetData == `` {
		return nil
	}

	var certs []*ProxyLBCertificate
	if err := json.Unmarshal(data, &certs); err != nil {
		return err
	}

	*p = certs
	return nil
}

// SetPrimaryCert PrimaryCertを設定
func (p *ProxyLBCertificates) SetPrimaryCert(cert *ProxyLBCertificate) {
	p.PrimaryCert = cert
}

// SetPrimaryCertValue PrimaryCertを設定
func (p *ProxyLBCertificates) SetPrimaryCertValue(serverCert, intermediateCert, privateKey string) {
	p.PrimaryCert = &ProxyLBCertificate{
		ServerCertificate:       serverCert,
		IntermediateCertificate: intermediateCert,
		PrivateKey:              privateKey,
	}
}

// AddAdditionalCert AdditionalCertを追加
func (p *ProxyLBCertificates) AddAdditionalCert(serverCert, intermediateCert, privateKey string) {
	p.AdditionalCerts = append(p.AdditionalCerts, &ProxyLBCertificate{
		ServerCertificate:       serverCert,
		IntermediateCertificate: intermediateCert,
		PrivateKey:              privateKey,
	})
}

// RemoveAdditionalCertAt 指定のインデックスを持つAdditionalCertを削除
func (p *ProxyLBCertificates) RemoveAdditionalCertAt(index int) {
	var certs []*ProxyLBCertificate
	for i, cert := range p.AdditionalCerts {
		if i != index {
			certs = append(certs, cert)
		}
	}
	p.AdditionalCerts = certs
}

// RemoveAdditionalCert 指定の内容を持つAdditionalCertを削除
func (p *ProxyLBCertificates) RemoveAdditionalCert(serverCert, intermediateCert, privateKey string) {
	var certs []*ProxyLBCertificate
	for _, cert := range p.AdditionalCerts {
		if !(cert.ServerCertificate == serverCert && cert.IntermediateCertificate == intermediateCert && cert.PrivateKey == privateKey) {
			certs = append(certs, cert)
		}
	}
	p.AdditionalCerts = certs
}

// RemoveAdditionalCerts AdditionalCertsを全て削除
func (p *ProxyLBCertificates) RemoveAdditionalCerts() {
	p.AdditionalCerts = []*ProxyLBCertificate{}
}

// ParseServerCertificate サーバ証明書のパース
func (p *ProxyLBCertificates) ParseServerCertificate() (*x509.Certificate, error) {
	cert, e := p.parseCertificate(p.PrimaryCert.ServerCertificate)
	if e != nil {
		return nil, e
	}
	return cert, nil
}

// ParseIntermediateCertificate 中間証明書のパース
func (p *ProxyLBCertificates) ParseIntermediateCertificate() (*x509.Certificate, error) {
	cert, e := p.parseCertificate(p.PrimaryCert.IntermediateCertificate)
	if e != nil {
		return nil, e
	}
	return cert, nil
}

func (p *ProxyLBCertificates) parseCertificate(certPEM string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block != nil {
		return x509.ParseCertificate(block.Bytes)
	}
	return nil, fmt.Errorf("can't decode certificate")
}

// ProxyLBCertificate ProxyLBのSSL証明書詳細
type ProxyLBCertificate struct {
	ServerCertificate       string    // サーバ証明書
	IntermediateCertificate string    // 中間証明書
	PrivateKey              string    // 秘密鍵
	CertificateEndDate      time.Time `json:",omitempty"` // 有効期限
	CertificateCommonName   string    `json:",omitempty"` // CommonName
}

// UnmarshalJSON UnmarshalJSON(CertificateEndDateのtime.TimeへのUnmarshal対応)
func (p *ProxyLBCertificate) UnmarshalJSON(data []byte) error {
	var tmp map[string]interface{}
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	p.ServerCertificate = tmp["ServerCertificate"].(string)
	p.IntermediateCertificate = tmp["IntermediateCertificate"].(string)
	p.PrivateKey = tmp["PrivateKey"].(string)
	p.CertificateCommonName = tmp["CertificateCommonName"].(string)
	endDate := tmp["CertificateEndDate"].(string)
	if endDate != "" {
		date, err := time.Parse("Jan _2 15:04:05 2006 MST", endDate)
		if err != nil {
			return err
		}
		p.CertificateEndDate = date
	}

	return nil
}

// ParseServerCertificate サーバ証明書のパース
func (p *ProxyLBCertificate) ParseServerCertificate() (*x509.Certificate, error) {
	cert, e := p.parseCertificate(p.ServerCertificate)
	if e != nil {
		return nil, e
	}
	return cert, nil
}

// ParseIntermediateCertificate 中間証明書のパース
func (p *ProxyLBCertificate) ParseIntermediateCertificate() (*x509.Certificate, error) {
	cert, e := p.parseCertificate(p.IntermediateCertificate)
	if e != nil {
		return nil, e
	}
	return cert, nil
}

func (p *ProxyLBCertificate) parseCertificate(certPEM string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(certPEM))
	if block != nil {
		return x509.ParseCertificate(block.Bytes)
	}
	return nil, fmt.Errorf("can't decode certificate")
}

// ProxyLBHealth ProxyLBのヘルスチェック戻り値
type ProxyLBHealth struct {
	ActiveConn int                    // アクティブなコネクション数
	CPS        int                    // 秒あたりコネクション数
	Servers    []*ProxyLBHealthServer // 実サーバのステータス
	CurrentVIP string                 // 現在のVIP
}

// ProxyLBHealthServer ProxyLBの実サーバのステータス
type ProxyLBHealthServer struct {
	ActiveConn int    // アクティブなコネクション数
	Status     string // ステータス(UP or DOWN)
	IPAddress  string // IPアドレス
	Port       string // ポート
	CPS        int    // 秒あたりコネクション数
}
