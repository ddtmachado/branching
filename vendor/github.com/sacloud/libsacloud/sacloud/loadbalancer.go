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

import "strconv"

// LoadBalancer ロードバランサー
type LoadBalancer struct {
	*Appliance // アプライアンス共通属性

	Remark   *LoadBalancerRemark   `json:",omitempty"` // リマーク
	Settings *LoadBalancerSettings `json:",omitempty"` // ロードバランサー設定
}

// IsHA 冗長化されている場合にtrueを返す
func (l *LoadBalancer) IsHA() bool {
	isHA := false
	if len(l.Remark.Servers) > 1 {
		if v, ok := l.Remark.Servers[1].(map[string]string); ok {
			if _, ok := v["IPAddress"]; ok {
				isHA = true
			}
		}
	}
	return isHA
}

// IPAddress1 ロードバランサ本体のIPアドレス(1番目)を返す
func (l *LoadBalancer) IPAddress1() string {
	if len(l.Remark.Servers) > 0 {
		if v, ok := l.Remark.Servers[0].(map[string]string); ok {
			if v, ok := v["IPAddress"]; ok {
				return v
			}
		}
	}
	return ""
}

// IPAddress2 ロードバランサ本体のIPアドレス(2番目)を返す
func (l *LoadBalancer) IPAddress2() string {
	if len(l.Remark.Servers) > 1 {
		if v, ok := l.Remark.Servers[1].(map[string]string); ok {
			if v, ok := v["IPAddress"]; ok {
				return v
			}
		}
	}
	return ""
}

// LoadBalancerRemark リマーク
type LoadBalancerRemark struct {
	*ApplianceRemarkBase
	// TODO Zone
	//Zone *Resource
}

// LoadBalancerSettings ロードバランサー設定リスト
type LoadBalancerSettings struct {
	LoadBalancer []*LoadBalancerSetting // ロードバランサー設定リスト
}

// LoadBalancerSetting ロードバランサー仮想IP設定
type LoadBalancerSetting struct {
	VirtualIPAddress string                `json:",omitempty"` // 仮想IPアドレス
	Port             string                `json:",omitempty"` // ポート番号
	DelayLoop        string                `json:",omitempty"` // 監視間隔
	SorryServer      string                `json:",omitempty"` // ソーリーサーバー
	Description      string                `json:",omitempty"` // 説明
	Servers          []*LoadBalancerServer `json:",omitempty"` // 仮想IP配下の実サーバー
}

// LoadBalancerServer 仮想IP設定配下のサーバー
type LoadBalancerServer struct {
	IPAddress   string                   `json:",omitempty"` // IPアドレス
	Port        string                   `json:",omitempty"` // ポート番号
	HealthCheck *LoadBalancerHealthCheck `json:",omitempty"` // ヘルスチェック
	Enabled     string                   `json:",omitempty"` // 有効/無効
	Status      string                   `json:",omitempty"` // ステータス
	ActiveConn  string                   `json:",omitempty"` // アクティブなコネクション
}

// LoadBalancerHealthCheck ヘルスチェック
type LoadBalancerHealthCheck struct {
	Protocol string `json:",omitempty"` // プロトコル
	Path     string `json:",omitempty"` // HTTP/HTTPSの場合のリクエストパス
	Status   string `json:",omitempty"` // HTTP/HTTPSの場合の期待するレスポンスコード
}

// LoadBalancerPlan ロードバランサープラン
type LoadBalancerPlan int

var (
	// LoadBalancerPlanStandard スタンダードプラン
	LoadBalancerPlanStandard = LoadBalancerPlan(1)
	// LoadBalancerPlanPremium プレミアムプラン
	LoadBalancerPlanPremium = LoadBalancerPlan(2)
)

// CreateLoadBalancerValue ロードバランサー作成用パラメーター
type CreateLoadBalancerValue struct {
	SwitchID     ID               // 接続先スイッチID
	VRID         int              // VRID
	Plan         LoadBalancerPlan // プラン
	IPAddress1   string           // IPアドレス
	MaskLen      int              // ネットワークマスク長
	DefaultRoute string           // デフォルトルート
	Name         string           // 名称
	Description  string           // 説明
	Tags         []string         // タグ
	Icon         *Resource        // アイコン
}

// CreateDoubleLoadBalancerValue ロードバランサー(冗長化あり)作成用パラメーター
type CreateDoubleLoadBalancerValue struct {
	*CreateLoadBalancerValue
	IPAddress2 string // IPアドレス2
}

// AllowLoadBalancerHealthCheckProtocol ロードバランサーでのヘルスチェック対応プロトコルリスト
func AllowLoadBalancerHealthCheckProtocol() []string {
	return []string{"http", "https", "ping", "tcp"}
}

// CreateNewLoadBalancerSingle ロードバランサー作成(冗長化なし)
func CreateNewLoadBalancerSingle(values *CreateLoadBalancerValue, settings []*LoadBalancerSetting) (*LoadBalancer, error) {

	lb := &LoadBalancer{
		Appliance: &Appliance{
			Class:           "loadbalancer",
			propName:        propName{Name: values.Name},
			propDescription: propDescription{Description: values.Description},
			propTags:        propTags{Tags: values.Tags},
			propPlanID:      propPlanID{Plan: &Resource{ID: ID(values.Plan)}},
			propIcon: propIcon{
				&Icon{
					Resource: values.Icon,
				},
			},
		},
		Remark: &LoadBalancerRemark{
			ApplianceRemarkBase: &ApplianceRemarkBase{
				Switch: &ApplianceRemarkSwitch{
					ID: values.SwitchID,
				},
				VRRP: &ApplianceRemarkVRRP{
					VRID: values.VRID,
				},
				Network: &ApplianceRemarkNetwork{
					NetworkMaskLen: values.MaskLen,
					DefaultRoute:   values.DefaultRoute,
				},
				Servers: []interface{}{
					map[string]string{"IPAddress": values.IPAddress1},
				},
			},
		},
	}

	for _, s := range settings {
		lb.AddLoadBalancerSetting(s)
	}

	return lb, nil
}

// CreateNewLoadBalancerDouble ロードバランサー(冗長化あり)作成
func CreateNewLoadBalancerDouble(values *CreateDoubleLoadBalancerValue, settings []*LoadBalancerSetting) (*LoadBalancer, error) {
	lb, err := CreateNewLoadBalancerSingle(values.CreateLoadBalancerValue, settings)
	if err != nil {
		return nil, err
	}
	lb.Remark.Servers = append(lb.Remark.Servers, map[string]string{"IPAddress": values.IPAddress2})
	return lb, nil
}

// AddLoadBalancerSetting ロードバランサー仮想IP設定追加
//
// ロードバランサー設定は仮想IPアドレス単位で保持しています。
// 仮想IPを増やす場合にこのメソッドを利用します。
func (l *LoadBalancer) AddLoadBalancerSetting(setting *LoadBalancerSetting) {
	if l.Settings == nil {
		l.Settings = &LoadBalancerSettings{}
	}
	if l.Settings.LoadBalancer == nil {
		l.Settings.LoadBalancer = []*LoadBalancerSetting{}
	}
	l.Settings.LoadBalancer = append(l.Settings.LoadBalancer, setting)
}

// DeleteLoadBalancerSetting ロードバランサー仮想IP設定の削除
func (l *LoadBalancer) DeleteLoadBalancerSetting(vip string, port string) {
	res := []*LoadBalancerSetting{}
	for _, l := range l.Settings.LoadBalancer {
		if l.VirtualIPAddress != vip || l.Port != port {
			res = append(res, l)
		}
	}

	l.Settings.LoadBalancer = res
}

// AddServer 仮想IP設定配下へ実サーバーを追加
func (s *LoadBalancerSetting) AddServer(server *LoadBalancerServer) {
	if s.Servers == nil {
		s.Servers = []*LoadBalancerServer{}
	}
	s.Servers = append(s.Servers, server)
}

// DeleteServer 仮想IP設定配下の実サーバーを削除
func (s *LoadBalancerSetting) DeleteServer(ip string, port string) {
	res := []*LoadBalancerServer{}
	for _, server := range s.Servers {
		if server.IPAddress != ip || server.Port != port {
			res = append(res, server)
		}
	}

	s.Servers = res

}

// LoadBalancerStatusResult ロードバランサーのステータスAPI戻り値
type LoadBalancerStatusResult []*LoadBalancerStatus

// Get VIPに対応するステータスを取得
func (l *LoadBalancerStatusResult) Get(vip string) *LoadBalancerStatus {
	for _, v := range *l {
		if v.VirtualIPAddress == vip {
			return v
		}
	}
	return nil
}

// LoadBalancerStatus ロードバランサーのステータス
type LoadBalancerStatus struct {
	VirtualIPAddress string
	Port             string
	Servers          []*LoadBalancerServerStatus `json:",omitempty"`
	CPS              string
}

// Get IPアドレスに対応する実サーバのステータスを取得
func (l *LoadBalancerStatus) Get(ip string) *LoadBalancerServerStatus {
	for _, v := range l.Servers {
		if v.IPAddress == ip {
			return v
		}
	}
	return nil
}

// NumCPS CPSを数値にして返す
func (l *LoadBalancerStatus) NumCPS() int {
	v, _ := strconv.Atoi(l.CPS) // nolint - ignore error
	return v
}

// NumPort Portを数値にして返す
func (l *LoadBalancerStatus) NumPort() int {
	v, _ := strconv.Atoi(l.Port) // nolint - ignore error
	return v
}

// LoadBalancerServerStatus ロードバランサーのVIP配下の実サーバのステータス
type LoadBalancerServerStatus struct {
	ActiveConn string
	IPAddress  string
	Status     string
	Port       string
	CPS        string
}

// NumActiveConn ActiveConnを数値にして返す
func (l *LoadBalancerServerStatus) NumActiveConn() int {
	v, _ := strconv.Atoi(l.ActiveConn) // nolint - ignore error
	return v
}

// NumCPS CPSを数値にして返す
func (l *LoadBalancerServerStatus) NumCPS() int {
	v, _ := strconv.Atoi(l.CPS) // nolint - ignore error
	return v
}

// NumPort Portを数値にして返す
func (l *LoadBalancerServerStatus) NumPort() int {
	v, _ := strconv.Atoi(l.Port) // nolint - ignore error
	return v
}
