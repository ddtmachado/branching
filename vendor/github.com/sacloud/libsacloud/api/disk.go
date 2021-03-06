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

package api

import (
	"fmt"
	"time"

	"github.com/sacloud/libsacloud/sacloud"
)

var (
	// allowDiskEditTags ディスクの編集可否判定に用いるタグ
	allowDiskEditTags = []string{
		"os-unix",
		"os-linux",
	}

	// bundleInfoWindowsHostClass ディスクの編集可否判定に用いる、BundleInfoでのWindows判定文字列
	bundleInfoWindowsHostClass = "ms_windows"
)

// DiskAPI ディスクAPI
type DiskAPI struct {
	*baseAPI
}

// NewDiskAPI ディスクAPI作成
func NewDiskAPI(client *Client) *DiskAPI {
	return &DiskAPI{
		&baseAPI{
			client: client,
			// FuncGetResourceURL
			FuncGetResourceURL: func() string {
				return "disk"
			},
		},
	}
}

// SortByConnectionOrder 接続順でのソート
func (api *DiskAPI) SortByConnectionOrder(reverse bool) *DiskAPI {
	api.sortBy("ConnectionOrder", reverse)
	return api
}

// WithServerID サーバーID条件
func (api *DiskAPI) WithServerID(id sacloud.ID) *DiskAPI {
	api.FilterBy("Server.ID", id)
	return api
}

// Create 新規作成
func (api *DiskAPI) Create(value *sacloud.Disk) (*sacloud.Disk, error) {
	//HACK: さくらのAPI側仕様: 戻り値:Successがbool値へ変換できないため文字列で受ける
	type diskResponse struct {
		*sacloud.Response
		// Success
		Success string `json:",omitempty"`
	}
	res := &diskResponse{}

	rawBody := &sacloud.Request{}
	rawBody.Disk = value
	if len(value.DistantFrom) > 0 {
		rawBody.DistantFrom = value.DistantFrom
		value.DistantFrom = []sacloud.ID{}
	}

	err := api.create(rawBody, res)
	if err != nil {
		return nil, err
	}
	return res.Disk, nil
}

// CreateWithConfig ディスク作成とディスクの修正、サーバ起動(指定されていれば)を１回のAPI呼び出しで実行
func (api *DiskAPI) CreateWithConfig(value *sacloud.Disk, config *sacloud.DiskEditValue, bootAtAvailable bool) (*sacloud.Disk, error) {
	//HACK: さくらのAPI側仕様: 戻り値:Successがbool値へ変換できないため文字列で受ける("Accepted"などが返る)
	type diskResponse struct {
		*sacloud.Response
		// Success
		Success string `json:",omitempty"`
	}
	res := &diskResponse{}

	type diskRequest struct {
		*sacloud.Request
		Config          *sacloud.DiskEditValue `json:",omitempty"`
		BootAtAvailable bool                   `json:",omitempty"`
	}

	rawBody := &diskRequest{
		Request:         &sacloud.Request{},
		BootAtAvailable: bootAtAvailable,
	}
	rawBody.Disk = value
	rawBody.Config = config

	if len(value.DistantFrom) > 0 {
		rawBody.DistantFrom = value.DistantFrom
		value.DistantFrom = []sacloud.ID{}
	}

	err := api.create(rawBody, res)
	if err != nil {
		return nil, err
	}
	return res.Disk, nil
}

// NewCondig ディスクの修正用パラメーター作成
func (api *DiskAPI) NewCondig() *sacloud.DiskEditValue {
	return &sacloud.DiskEditValue{}
}

// Config ディスクの修正
func (api *DiskAPI) Config(id sacloud.ID, disk *sacloud.DiskEditValue) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/config", api.getResourceURL(), id)
	)

	// HACK APIからの戻り値"Success"がboolではなく文字列となっているためapi.modifyを使わずに実装する。
	res := &struct {
		IsOk bool `json:"is_ok,omitempty"` // is_ok項目
	}{}
	err := api.baseAPI.request(method, uri, disk, res)
	if err != nil {
		return false, err
	}
	return res.IsOk, nil

}

func (api *DiskAPI) install(id sacloud.ID, body *sacloud.Disk) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/install", api.getResourceURL(), id)
	)
	//HACK: さくらのAPI側仕様: 戻り値:Successがbool値へ変換できないため文字列で受ける
	type diskResponse struct {
		*sacloud.ResultFlagValue
		// Success
		Success string `json:",omitempty"`
	}
	res := &diskResponse{}
	rawBody := &sacloud.Request{}
	rawBody.Disk = body
	if len(body.DistantFrom) > 0 {
		rawBody.DistantFrom = body.DistantFrom
		body.DistantFrom = []sacloud.ID{}
	}

	err := api.baseAPI.request(method, uri, rawBody, res)
	if err != nil {
		return false, err
	}
	return res.IsOk, nil
}

// ReinstallFromBlank ブランクディスクから再インストール
func (api *DiskAPI) ReinstallFromBlank(id sacloud.ID, sizeMB int) (bool, error) {
	var body = &sacloud.Disk{}
	body.SetSizeMB(sizeMB)

	return api.install(id, body)
}

// ReinstallFromArchive アーカイブからの再インストール
func (api *DiskAPI) ReinstallFromArchive(id sacloud.ID, archiveID sacloud.ID, distantFrom ...sacloud.ID) (bool, error) {
	var body = &sacloud.Disk{}
	body.SetSourceArchive(archiveID)
	if len(distantFrom) > 0 {
		body.SetDistantFrom(distantFrom)
	}
	return api.install(id, body)
}

// ReinstallFromDisk ディスクからの再インストール
func (api *DiskAPI) ReinstallFromDisk(id sacloud.ID, diskID sacloud.ID, distantFrom ...sacloud.ID) (bool, error) {
	var body = &sacloud.Disk{}
	body.SetSourceDisk(diskID)
	if len(distantFrom) > 0 {
		body.SetDistantFrom(distantFrom)
	}
	return api.install(id, body)
}

// ToBlank ディスクを空にする
func (api *DiskAPI) ToBlank(diskID sacloud.ID) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/to/blank", api.getResourceURL(), diskID)
	)
	return api.modify(method, uri, nil)
}

// ResizePartition パーティションのリサイズ
func (api *DiskAPI) ResizePartition(diskID sacloud.ID) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/resize-partition", api.getResourceURL(), diskID)
	)
	// HACK APIからの戻り値"Success"がboolではなく文字列となっているためapi.modifyを使わずに実装する。
	res := &struct {
		IsOk bool `json:"is_ok,omitempty"` // is_ok項目
	}{}
	err := api.baseAPI.request(method, uri, nil, res)
	if err != nil {
		return false, err
	}
	return res.IsOk, nil
}

// ResizePartitionBackground パーティションのリサイズ(非同期)
func (api *DiskAPI) ResizePartitionBackground(diskID sacloud.ID) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/resize-partition", api.getResourceURL(), diskID)
		body   = map[string]interface{}{
			"Background": true,
		}
	)

	// HACK APIからの戻り値"Success"がboolではなく文字列となっているためapi.modifyを使わずに実装する。
	res := &struct {
		IsOk bool `json:"is_ok,omitempty"` // is_ok項目
	}{}
	err := api.baseAPI.request(method, uri, body, res)
	if err != nil {
		return false, err
	}
	return res.IsOk, nil
}

// DisconnectFromServer サーバーとの接続解除
func (api *DiskAPI) DisconnectFromServer(diskID sacloud.ID) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%d/to/server", api.getResourceURL(), diskID)
	)
	return api.modify(method, uri, nil)
}

// ConnectToServer サーバーとの接続
func (api *DiskAPI) ConnectToServer(diskID sacloud.ID, serverID sacloud.ID) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/to/server/%d", api.getResourceURL(), diskID, serverID)
	)
	return api.modify(method, uri, nil)
}

// State ディスクの状態を取得し有効な状態か判定
func (api *DiskAPI) State(diskID sacloud.ID) (bool, error) {
	disk, err := api.Read(diskID)
	if err != nil {
		return false, err
	}
	return disk.IsAvailable(), nil
}

// SleepWhileCopying コピー終了まで待機
func (api *DiskAPI) SleepWhileCopying(id sacloud.ID, timeout time.Duration) error {
	handler := waitingForAvailableFunc(func() (hasAvailable, error) {
		return api.Read(id)
	}, 0)
	return blockingPoll(handler, timeout)
}

// AsyncSleepWhileCopying コピー終了まで待機(非同期)
func (api *DiskAPI) AsyncSleepWhileCopying(id sacloud.ID, timeout time.Duration) (chan (interface{}), chan (interface{}), chan (error)) {
	handler := waitingForAvailableFunc(func() (hasAvailable, error) {
		return api.Read(id)
	}, 0)
	return poll(handler, timeout)
}

// Monitor アクティビティーモニター取得
func (api *DiskAPI) Monitor(id sacloud.ID, body *sacloud.ResourceMonitorRequest) (*sacloud.MonitorValues, error) {
	return api.baseAPI.monitor(id, body)
}

// CanEditDisk ディスクの修正が可能か判定
func (api *DiskAPI) CanEditDisk(id sacloud.ID) (bool, error) {

	disk, err := api.Read(id)
	if err != nil {
		return false, err
	}

	if disk == nil {
		return false, nil
	}

	// BundleInfoがあれば編集不可
	if disk.BundleInfo != nil && disk.BundleInfo.HostClass == bundleInfoWindowsHostClass {
		// Windows
		return false, nil
	}

	// SophosUTMであれば編集不可
	if disk.HasTag("pkg-sophosutm") || disk.IsSophosUTM() {
		return false, nil
	}
	// OPNsenseであれば編集不可
	if disk.HasTag("distro-opnsense") {
		return false, nil
	}
	// Netwiser VEであれば編集不可
	if disk.HasTag("pkg-netwiserve") {
		return false, nil
	}

	// ソースアーカイブ/ソースディスクともに持っていない場合
	if disk.SourceArchive == nil && disk.SourceDisk == nil {
		//ブランクディスクがソース
		return false, nil
	}

	for _, t := range allowDiskEditTags {
		if disk.HasTag(t) {
			// 対応OSインストール済みディスク
			return true, nil
		}
	}

	// ここまできても判定できないならソースに投げる
	if disk.SourceDisk != nil && disk.SourceDisk.Availability != "discontinued" {
		return api.client.Disk.CanEditDisk(disk.SourceDisk.ID)
	}
	if disk.SourceArchive != nil && disk.SourceArchive.Availability != "discontinued" {
		return api.client.Archive.CanEditDisk(disk.SourceArchive.ID)
	}

	return false, nil

}

// GetPublicArchiveIDFromAncestors 祖先の中からパブリックアーカイブのIDを検索
func (api *DiskAPI) GetPublicArchiveIDFromAncestors(id sacloud.ID) (sacloud.ID, bool) {

	emptyID := sacloud.EmptyID

	disk, err := api.Read(id)
	if err != nil {
		return emptyID, false
	}

	if disk == nil {
		return emptyID, false
	}

	// BundleInfoがあれば編集不可
	if disk.BundleInfo != nil && disk.BundleInfo.HostClass == bundleInfoWindowsHostClass {
		// Windows
		return emptyID, false
	}

	// SophosUTMであれば編集不可
	if disk.HasTag("pkg-sophosutm") || disk.IsSophosUTM() {
		return emptyID, false
	}
	// OPNsenseであれば編集不可
	if disk.HasTag("distro-opnsense") {
		return emptyID, false
	}
	// Netwiser VEであれば編集不可
	if disk.HasTag("pkg-netwiserve") {
		return emptyID, false
	}

	for _, t := range allowDiskEditTags {
		if disk.HasTag(t) {
			// 対応OSインストール済みディスク
			return disk.ID, true
		}
	}

	// ここまできても判定できないならソースに投げる
	if disk.SourceDisk != nil && disk.SourceDisk.Availability != "discontinued" {
		return api.client.Disk.GetPublicArchiveIDFromAncestors(disk.SourceDisk.ID)
	}
	if disk.SourceArchive != nil && disk.SourceArchive.Availability != "discontinued" {
		return api.client.Archive.GetPublicArchiveIDFromAncestors(disk.SourceArchive.ID)
	}
	return emptyID, false

}
