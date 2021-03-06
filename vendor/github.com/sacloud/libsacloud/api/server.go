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

// ServerAPI サーバーAPI
type ServerAPI struct {
	*baseAPI
}

// NewServerAPI サーバーAPI作成
func NewServerAPI(client *Client) *ServerAPI {
	return &ServerAPI{
		&baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "server"
			},
		},
	}
}

// WithPlan サーバープラン条件
func (api *ServerAPI) WithPlan(planID string) *ServerAPI {
	return api.FilterBy("ServerPlan.ID", planID)
}

// WithStatus インスタンスステータス条件
func (api *ServerAPI) WithStatus(status string) *ServerAPI {
	return api.FilterBy("Instance.Status", status)
}

// WithStatusUp 起動状態条件
func (api *ServerAPI) WithStatusUp() *ServerAPI {
	return api.WithStatus("up")
}

// WithStatusDown ダウン状態条件
func (api *ServerAPI) WithStatusDown() *ServerAPI {
	return api.WithStatus("down")
}

// WithISOImage ISOイメージ条件
func (api *ServerAPI) WithISOImage(imageID sacloud.ID) *ServerAPI {
	return api.FilterBy("Instance.CDROM.ID", imageID)
}

// SortByCPU CPUコア数でのソート
func (api *ServerAPI) SortByCPU(reverse bool) *ServerAPI {
	api.sortBy("ServerPlan.CPU", reverse)
	return api
}

// SortByMemory メモリサイズでのソート
func (api *ServerAPI) SortByMemory(reverse bool) *ServerAPI {
	api.sortBy("ServerPlan.MemoryMB", reverse)
	return api
}

// DeleteWithDisk 指定のディスクと共に削除する
func (api *ServerAPI) DeleteWithDisk(id sacloud.ID, disks []sacloud.ID) (*sacloud.Server, error) {
	return api.request(func(res *sacloud.Response) error {
		return api.delete(id, map[string]interface{}{"WithDisk": disks}, res)
	})
}

// State ステータス(Availability)取得
func (api *ServerAPI) State(id sacloud.ID) (string, error) {
	server, err := api.Read(id)
	if err != nil {
		return "", err
	}
	return string(server.Availability), nil
}

// IsUp 起動しているか判定
func (api *ServerAPI) IsUp(id sacloud.ID) (bool, error) {
	server, err := api.Read(id)
	if err != nil {
		return false, err
	}
	return server.Instance.IsUp(), nil
}

// IsDown ダウンしているか判定
func (api *ServerAPI) IsDown(id sacloud.ID) (bool, error) {
	server, err := api.Read(id)
	if err != nil {
		return false, err
	}
	return server.Instance.IsDown(), nil
}

// Boot 起動
func (api *ServerAPI) Boot(id sacloud.ID) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/power", api.getResourceURL(), id)
	)
	return api.modify(method, uri, nil)
}

// Shutdown シャットダウン(graceful)
func (api *ServerAPI) Shutdown(id sacloud.ID) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%d/power", api.getResourceURL(), id)
	)

	return api.modify(method, uri, nil)
}

// Stop シャットダウン(force)
func (api *ServerAPI) Stop(id sacloud.ID) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%d/power", api.getResourceURL(), id)
	)

	return api.modify(method, uri, map[string]bool{"Force": true})
}

// RebootForce 再起動
func (api *ServerAPI) RebootForce(id sacloud.ID) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/reset", api.getResourceURL(), id)
	)

	return api.modify(method, uri, nil)
}

// SleepUntilUp 起動するまで待機
func (api *ServerAPI) SleepUntilUp(id sacloud.ID, timeout time.Duration) error {
	handler := waitingForUpFunc(func() (hasUpDown, error) {
		return api.Read(id)
	}, 0)
	return blockingPoll(handler, timeout)
}

// SleepUntilDown ダウンするまで待機
func (api *ServerAPI) SleepUntilDown(id sacloud.ID, timeout time.Duration) error {
	handler := waitingForDownFunc(func() (hasUpDown, error) {
		return api.Read(id)
	}, 0)
	return blockingPoll(handler, timeout)
}

// ChangePlan サーバープラン変更(サーバーIDが変更となるため注意)
func (api *ServerAPI) ChangePlan(serverID sacloud.ID, plan *sacloud.ProductServer) (*sacloud.Server, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/plan", api.getResourceURL(), serverID)
		body   = &sacloud.ProductServer{}
	)
	body.CPU = plan.CPU
	body.MemoryMB = plan.MemoryMB
	body.Generation = plan.Generation
	body.Commitment = plan.Commitment

	return api.request(func(res *sacloud.Response) error {
		return api.baseAPI.request(method, uri, body, res)
	})
}

// FindDisk 指定サーバーに接続されているディスク一覧を取得
func (api *ServerAPI) FindDisk(serverID sacloud.ID) ([]sacloud.Disk, error) {
	server, err := api.Read(serverID)
	if err != nil {
		return nil, err
	}
	return server.Disks, nil
}

// InsertCDROM ISOイメージを挿入
func (api *ServerAPI) InsertCDROM(serverID sacloud.ID, cdromID sacloud.ID) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/cdrom", api.getResourceURL(), serverID)
	)

	req := &sacloud.Request{
		SakuraCloudResources: sacloud.SakuraCloudResources{
			CDROM: &sacloud.CDROM{Resource: &sacloud.Resource{ID: cdromID}},
		},
	}

	return api.modify(method, uri, req)
}

// EjectCDROM ISOイメージを取り出し
func (api *ServerAPI) EjectCDROM(serverID sacloud.ID, cdromID sacloud.ID) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%d/cdrom", api.getResourceURL(), serverID)
	)

	req := &sacloud.Request{
		SakuraCloudResources: sacloud.SakuraCloudResources{
			CDROM: &sacloud.CDROM{Resource: &sacloud.Resource{ID: cdromID}},
		},
	}

	return api.modify(method, uri, req)
}

// NewKeyboardRequest キーボード入力リクエストパラメーター作成
func (api *ServerAPI) NewKeyboardRequest() *sacloud.KeyboardRequest {
	return &sacloud.KeyboardRequest{}
}

// SendKey キーボード入力送信
func (api *ServerAPI) SendKey(serverID sacloud.ID, body *sacloud.KeyboardRequest) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/keyboard", api.getResourceURL(), serverID)
	)

	return api.modify(method, uri, body)
}

// NewMouseRequest マウス入力リクエストパラメーター作成
func (api *ServerAPI) NewMouseRequest() *sacloud.MouseRequest {
	return &sacloud.MouseRequest{
		Buttons: &sacloud.MouseRequestButtons{},
	}
}

// SendMouse マウス入力送信
func (api *ServerAPI) SendMouse(serverID sacloud.ID, mouseIndex string, body *sacloud.MouseRequest) (bool, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/mouse/%s", api.getResourceURL(), serverID, mouseIndex)
	)

	return api.modify(method, uri, body)
}

// NewVNCSnapshotRequest VNCスナップショット取得リクエストパラメーター作成
func (api *ServerAPI) NewVNCSnapshotRequest() *sacloud.VNCSnapshotRequest {
	return &sacloud.VNCSnapshotRequest{}
}

// GetVNCProxy VNCプロキシ情報取得
func (api *ServerAPI) GetVNCProxy(serverID sacloud.ID) (*sacloud.VNCProxyResponse, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("%s/%d/vnc/proxy", api.getResourceURL(), serverID)
		res    = &sacloud.VNCProxyResponse{}
	)
	err := api.baseAPI.request(method, uri, nil, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetVNCSize VNC画面サイズ取得
func (api *ServerAPI) GetVNCSize(serverID sacloud.ID) (*sacloud.VNCSizeResponse, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("%s/%d/vnc/size", api.getResourceURL(), serverID)
		res    = &sacloud.VNCSizeResponse{}
	)
	err := api.baseAPI.request(method, uri, nil, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetVNCSnapshot VNCスナップショット取得
func (api *ServerAPI) GetVNCSnapshot(serverID sacloud.ID, body *sacloud.VNCSnapshotRequest) (*sacloud.VNCSnapshotResponse, error) {
	var (
		method = "GET"
		uri    = fmt.Sprintf("%s/%d/vnc/snapshot", api.getResourceURL(), serverID)
		res    = &sacloud.VNCSnapshotResponse{}
	)
	err := api.baseAPI.request(method, uri, body, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Monitor アクティビティーモニター(CPU-TIME)取得
func (api *ServerAPI) Monitor(id sacloud.ID, body *sacloud.ResourceMonitorRequest) (*sacloud.MonitorValues, error) {
	return api.baseAPI.monitor(id, body)
}
