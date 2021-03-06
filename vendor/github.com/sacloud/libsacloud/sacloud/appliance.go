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

// Appliance アプライアンス基底クラス
type Appliance struct {
	*Resource               // ID
	propAvailability        // 有功状態
	propName                // 名称
	propDescription         // 説明
	propPlanID              // プラン
	propServiceClass        // サービスクラス
	propInstance            // インスタンス
	propSwitch              // スイッチ
	propInterfaces          // インターフェース配列
	propIcon                // アイコン
	propTags                // タグ
	propCreatedAt           // 作成日時
	Class            string `json:",omitempty"` // リソースクラス
	SettingHash      string `json:",omitempty"` // 設定ハッシュ値

}

//HACK Appliance:Zone.IDがRoute/LoadBalancerの場合でデータ型が異なるため
//それぞれのstruct定義でZoneだけ上書きした構造体を定義して使う

// ApplianceRemarkBase アプライアンス Remark 基底クラス
type ApplianceRemarkBase struct {
	Servers []interface{} // 配下のサーバー群

	Switch  *ApplianceRemarkSwitch  `json:",omitempty"` // 接続先スイッチ
	VRRP    *ApplianceRemarkVRRP    `json:",omitempty"` // VRRP
	Network *ApplianceRemarkNetwork `json:",omitempty"` // ネットワーク

	//Zone *Resource `json:",omitempty"`
	//Plan    *Resource
}

//type ApplianceServer struct {
//	IPAddress string `json:",omitempty"`
//}

// ApplianceRemarkSwitch スイッチ
type ApplianceRemarkSwitch struct {
	ID        ID `json:",omitempty"` // リソースID
	propScope    // スコープ
}

// ApplianceRemarkVRRP VRRP
type ApplianceRemarkVRRP struct {
	VRID int // VRID
}

// ApplianceRemarkNetwork ネットワーク
type ApplianceRemarkNetwork struct {
	NetworkMaskLen int    `json:",omitempty"` // ネットワークマスク長
	DefaultRoute   string // デフォルトルート

}
