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
	"strings"
	"time"

	"github.com/sacloud/libsacloud/sacloud"
	"github.com/sacloud/libsacloud/sacloud/ostype"
)

// ArchiveAPI アーカイブAPI
type ArchiveAPI struct {
	*baseAPI
	findFuncMapPerOSType map[ostype.ArchiveOSTypes]func() (*sacloud.Archive, error)
}

var (
	archiveLatestStableCentOSTags                          = []string{"current-stable", "distro-centos"}
	archiveLatestStableCentOS8Tags                         = []string{"centos-8-latest"}
	archiveLatestStableCentOS7Tags                         = []string{"centos-7-latest"}
	archiveLatestStableCentOS6Tags                         = []string{"centos-6-latest"}
	archiveLatestStableUbuntuTags                          = []string{"current-stable", "distro-ubuntu"}
	archiveLatestStableUbuntu1804Tags                      = []string{"ubuntu-18.04-latest"}
	archiveLatestStableUbuntu1604Tags                      = []string{"ubuntu-16.04-latest"}
	archiveLatestStableDebianTags                          = []string{"current-stable", "distro-debian"}
	archiveLatestStableDebian10Tags                        = []string{"debian-10-latest"}
	archiveLatestStableDebian9Tags                         = []string{"debian-9-latest"}
	archiveLatestStableCoreOSTags                          = []string{"current-stable", "distro-coreos"}
	archiveLatestStableRancherOSTags                       = []string{"current-stable", "distro-rancheros"}
	archiveLatestStableK3OSTags                            = []string{"current-stable", "distro-k3os"}
	archiveLatestStableKusanagiTags                        = []string{"current-stable", "pkg-kusanagi"}
	archiveLatestStableSophosUTMTags                       = []string{"current-stable", "pkg-sophosutm"}
	archiveLatestStableFreeBSDTags                         = []string{"current-stable", "distro-freebsd"}
	archiveLatestStableNetwiserTags                        = []string{"current-stable", "pkg-netwiserve"}
	archiveLatestStableOPNsenseTags                        = []string{"current-stable", "distro-opnsense"}
	archiveLatestStableWindows2016Tags                     = []string{"os-windows", "distro-ver-2016"}
	archiveLatestStableWindows2016RDSTags                  = []string{"os-windows", "distro-ver-2016", "windows-rds"}
	archiveLatestStableWindows2016RDSOfficeTags            = []string{"os-windows", "distro-ver-2016", "windows-rds", "with-office"}
	archiveLatestStableWindows2016SQLServerWeb             = []string{"os-windows", "distro-ver-2016", "windows-sqlserver", "sqlserver-2016", "edition-web"}
	archiveLatestStableWindows2016SQLServerStandard        = []string{"os-windows", "distro-ver-2016", "windows-sqlserver", "sqlserver-2016", "edition-standard"}
	archiveLatestStableWindows2016SQLServer2017Standard    = []string{"os-windows", "distro-ver-2016", "windows-sqlserver", "sqlserver-2017", "edition-standard"}
	archiveLatestStableWindows2016SQLServerStandardAll     = []string{"os-windows", "distro-ver-2016", "windows-sqlserver", "sqlserver-2016", "edition-standard", "windows-rds", "with-office"}
	archiveLatestStableWindows2016SQLServer2017StandardAll = []string{"os-windows", "distro-ver-2016", "windows-sqlserver", "sqlserver-2017", "edition-standard", "windows-rds", "with-office"}
	archiveLatestStableWindows2019Tags                     = []string{"os-windows", "distro-ver-2019"}
)

// NewArchiveAPI アーカイブAPI作成
func NewArchiveAPI(client *Client) *ArchiveAPI {
	api := &ArchiveAPI{
		baseAPI: &baseAPI{
			client: client,
			FuncGetResourceURL: func() string {
				return "archive"
			},
		},
	}

	api.findFuncMapPerOSType = map[ostype.ArchiveOSTypes]func() (*sacloud.Archive, error){
		ostype.CentOS:                              api.FindLatestStableCentOS,
		ostype.CentOS8:                             api.FindLatestStableCentOS8,
		ostype.CentOS7:                             api.FindLatestStableCentOS7,
		ostype.CentOS6:                             api.FindLatestStableCentOS6,
		ostype.Ubuntu:                              api.FindLatestStableUbuntu,
		ostype.Ubuntu1804:                          api.FindLatestStableUbuntu1804,
		ostype.Ubuntu1604:                          api.FindLatestStableUbuntu1604,
		ostype.Debian:                              api.FindLatestStableDebian,
		ostype.Debian10:                            api.FindLatestStableDebian10,
		ostype.Debian9:                             api.FindLatestStableDebian9,
		ostype.CoreOS:                              api.FindLatestStableCoreOS,
		ostype.RancherOS:                           api.FindLatestStableRancherOS,
		ostype.K3OS:                                api.FindLatestStableK3OS,
		ostype.Kusanagi:                            api.FindLatestStableKusanagi,
		ostype.SophosUTM:                           api.FindLatestStableSophosUTM,
		ostype.FreeBSD:                             api.FindLatestStableFreeBSD,
		ostype.Netwiser:                            api.FindLatestStableNetwiser,
		ostype.OPNsense:                            api.FindLatestStableOPNsense,
		ostype.Windows2016:                         api.FindLatestStableWindows2016,
		ostype.Windows2016RDS:                      api.FindLatestStableWindows2016RDS,
		ostype.Windows2016RDSOffice:                api.FindLatestStableWindows2016RDSOffice,
		ostype.Windows2016SQLServerWeb:             api.FindLatestStableWindows2016SQLServerWeb,
		ostype.Windows2016SQLServerStandard:        api.FindLatestStableWindows2016SQLServerStandard,
		ostype.Windows2016SQLServer2017Standard:    api.FindLatestStableWindows2016SQLServer2017Standard,
		ostype.Windows2016SQLServerStandardAll:     api.FindLatestStableWindows2016SQLServerStandardAll,
		ostype.Windows2016SQLServer2017StandardAll: api.FindLatestStableWindows2016SQLServer2017StandardAll,
	}

	return api
}

// OpenFTP FTP接続開始
func (api *ArchiveAPI) OpenFTP(id sacloud.ID) (*sacloud.FTPServer, error) {
	var (
		method = "PUT"
		uri    = fmt.Sprintf("%s/%d/ftp", api.getResourceURL(), id)
		//body   = map[string]bool{"ChangePassword": reset}
		res = &sacloud.Response{}
	)

	result, err := api.action(method, uri, nil, res)
	if !result || err != nil {
		return nil, err
	}

	return res.FTPServer, nil
}

// CloseFTP FTP接続終了
func (api *ArchiveAPI) CloseFTP(id sacloud.ID) (bool, error) {
	var (
		method = "DELETE"
		uri    = fmt.Sprintf("%s/%d/ftp", api.getResourceURL(), id)
	)
	return api.modify(method, uri, nil)

}

// SleepWhileCopying コピー終了まで待機
func (api *ArchiveAPI) SleepWhileCopying(id sacloud.ID, timeout time.Duration) error {
	handler := waitingForAvailableFunc(func() (hasAvailable, error) {
		return api.Read(id)
	}, 0)
	return blockingPoll(handler, timeout)
}

// AsyncSleepWhileCopying コピー終了まで待機(非同期)
func (api *ArchiveAPI) AsyncSleepWhileCopying(id sacloud.ID, timeout time.Duration) (chan (interface{}), chan (interface{}), chan (error)) {
	handler := waitingForAvailableFunc(func() (hasAvailable, error) {
		return api.Read(id)
	}, 0)
	return poll(handler, timeout)
}

// CanEditDisk ディスクの修正が可能か判定
func (api *ArchiveAPI) CanEditDisk(id sacloud.ID) (bool, error) {

	archive, err := api.Read(id)
	if err != nil {
		return false, err
	}

	if archive == nil {
		return false, nil
	}

	// BundleInfoがあれば編集不可
	if archive.BundleInfo != nil && archive.BundleInfo.HostClass == bundleInfoWindowsHostClass {
		// Windows
		return false, nil
	}

	// SophosUTMであれば編集不可
	if archive.HasTag("pkg-sophosutm") || archive.IsSophosUTM() {
		return false, nil
	}
	// OPNsenseであれば編集不可
	if archive.HasTag("distro-opnsense") {
		return false, nil
	}
	// Netwiser VEであれば編集不可
	if archive.HasTag("pkg-netwiserve") {
		return false, nil
	}
	// Juniper vSRXであれば編集不可
	if archive.HasTag("pkg-vsrx") {
		return false, nil
	}

	for _, t := range allowDiskEditTags {
		if archive.HasTag(t) {
			// 対応OSインストール済みディスク
			return true, nil
		}
	}

	// ここまできても判定できないならソースに投げる
	if archive.SourceDisk != nil && archive.SourceDisk.Availability != "discontinued" {
		return api.client.Disk.CanEditDisk(archive.SourceDisk.ID)
	}
	if archive.SourceArchive != nil && archive.SourceArchive.Availability != "discontinued" {
		return api.client.Archive.CanEditDisk(archive.SourceArchive.ID)
	}
	return false, nil

}

// GetPublicArchiveIDFromAncestors 祖先の中からパブリックアーカイブのIDを検索
func (api *ArchiveAPI) GetPublicArchiveIDFromAncestors(id sacloud.ID) (sacloud.ID, bool) {

	emptyID := sacloud.EmptyID

	archive, err := api.Read(id)
	if err != nil {
		return emptyID, false
	}

	if archive == nil {
		return emptyID, false
	}

	// BundleInfoがあれば編集不可
	if archive.BundleInfo != nil && archive.BundleInfo.HostClass == bundleInfoWindowsHostClass {
		// Windows
		return emptyID, false
	}

	// SophosUTMであれば編集不可
	if archive.HasTag("pkg-sophosutm") || archive.IsSophosUTM() {
		return emptyID, false
	}
	// OPNsenseであれば編集不可
	if archive.HasTag("distro-opnsense") {
		return emptyID, false
	}
	// Netwiser VEであれば編集不可
	if archive.HasTag("pkg-netwiserve") {
		return emptyID, false
	}
	// Juniper vSRXであれば編集不可
	if archive.HasTag("pkg-vsrx") {
		return emptyID, false
	}

	for _, t := range allowDiskEditTags {
		if archive.HasTag(t) {
			// 対応OSインストール済みディスク
			return archive.ID, true
		}
	}

	// ここまできても判定できないならソースに投げる
	if archive.SourceDisk != nil && archive.SourceDisk.Availability != "discontinued" {
		return api.client.Disk.GetPublicArchiveIDFromAncestors(archive.SourceDisk.ID)
	}
	if archive.SourceArchive != nil && archive.SourceArchive.Availability != "discontinued" {
		return api.client.Archive.GetPublicArchiveIDFromAncestors(archive.SourceArchive.ID)
	}
	return emptyID, false

}

// FindLatestStableCentOS 安定版最新のCentOSパブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableCentOS() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableCentOSTags)
}

// FindLatestStableCentOS8 安定版最新のCentOS8パブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableCentOS8() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableCentOS8Tags)
}

// FindLatestStableCentOS7 安定版最新のCentOS7パブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableCentOS7() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableCentOS7Tags)
}

// FindLatestStableCentOS6 安定版最新のCentOS6パブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableCentOS6() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableCentOS6Tags)
}

// FindLatestStableDebian 安定版最新のDebianパブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableDebian() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableDebianTags)
}

// FindLatestStableDebian10 安定版最新のDebian10パブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableDebian10() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableDebian10Tags)
}

// FindLatestStableDebian9 安定版最新のDebian9パブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableDebian9() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableDebian9Tags)
}

// FindLatestStableUbuntu 安定版最新のUbuntuパブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableUbuntu() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableUbuntuTags)
}

// FindLatestStableUbuntu1804 安定版最新のUbuntu1804パブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableUbuntu1804() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableUbuntu1804Tags)
}

// FindLatestStableUbuntu1604 安定版最新のUbuntu1604パブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableUbuntu1604() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableUbuntu1604Tags)
}

// FindLatestStableCoreOS 安定版最新のCoreOSパブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableCoreOS() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableCoreOSTags)
}

// FindLatestStableRancherOS 安定版最新のRancherOSパブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableRancherOS() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableRancherOSTags)
}

// FindLatestStableK3OS 安定版最新のk3OSパブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableK3OS() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableK3OSTags)
}

// FindLatestStableKusanagi 安定版最新のKusanagiパブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableKusanagi() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableKusanagiTags)
}

// FindLatestStableSophosUTM 安定板最新のSophosUTMパブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableSophosUTM() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableSophosUTMTags)
}

// FindLatestStableFreeBSD 安定版最新のFreeBSDパブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableFreeBSD() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableFreeBSDTags)
}

// FindLatestStableNetwiser 安定版最新のNetwiserパブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableNetwiser() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableNetwiserTags)
}

// FindLatestStableOPNsense 安定版最新のOPNsenseパブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableOPNsense() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableOPNsenseTags)
}

// FindLatestStableWindows2016 安定版最新のWindows2016パブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableWindows2016() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableWindows2016Tags, map[string]interface{}{
		"Name": "Windows Server 2016 Datacenter Edition",
	})
}

// FindLatestStableWindows2016RDS 安定版最新のWindows2016RDSパブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableWindows2016RDS() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableWindows2016RDSTags, map[string]interface{}{
		"Name": "Windows Server 2016 for RDS",
	})
}

// FindLatestStableWindows2016RDSOffice 安定版最新のWindows2016RDS(Office)パブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableWindows2016RDSOffice() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableWindows2016RDSOfficeTags, map[string]interface{}{
		"Name": "Windows Server 2016 for RDS(MS Office付)",
	})
}

// FindLatestStableWindows2016SQLServerWeb 安定版最新のWindows2016 SQLServer(Web) パブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableWindows2016SQLServerWeb() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableWindows2016SQLServerWeb, map[string]interface{}{
		"Name": "Windows Server 2016 for MS SQL 2016(Web)",
	})
}

// FindLatestStableWindows2016SQLServerStandard 安定版最新のWindows2016 SQLServer(Standard) パブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableWindows2016SQLServerStandard() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableWindows2016SQLServerStandard, map[string]interface{}{
		"Name": "Windows Server 2016 for MS SQL 2016(Standard)",
	})
}

// FindLatestStableWindows2016SQLServer2017Standard 安定版最新のWindows2016 SQLServer2017(Standard) パブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableWindows2016SQLServer2017Standard() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableWindows2016SQLServer2017Standard, map[string]interface{}{
		"Name": "Windows Server 2016 for MS SQL 2017(Standard)",
	})
}

// FindLatestStableWindows2016SQLServerStandardAll 安定版最新のWindows2016 SQLServer(RDS+Office) パブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableWindows2016SQLServerStandardAll() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableWindows2016SQLServerStandardAll, map[string]interface{}{
		"Name": "Windows Server 2016 for MS SQL 2016(Std) with RDS / MS Office",
	})
}

// FindLatestStableWindows2016SQLServer2017StandardAll 安定版最新のWindows2016 SQLServer2017(RDS+Office) パブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableWindows2016SQLServer2017StandardAll() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableWindows2016SQLServer2017StandardAll, map[string]interface{}{
		"Name": "Windows Server 2016 for MS SQL 2017(Std) with RDS / MS Office",
	})
}

// FindLatestStableWindows2019 安定版最新のWindows2019パブリックアーカイブを取得
func (api *ArchiveAPI) FindLatestStableWindows2019() (*sacloud.Archive, error) {
	return api.findByOSTags(archiveLatestStableWindows2019Tags, map[string]interface{}{
		"Name": "Windows Server 2019 Datacenter Edition",
	})
}

// FindByOSType 指定のOS種別の安定版最新のパブリックアーカイブを取得
func (api *ArchiveAPI) FindByOSType(os ostype.ArchiveOSTypes) (*sacloud.Archive, error) {
	if f, ok := api.findFuncMapPerOSType[os]; ok {
		return f()
	}

	return nil, fmt.Errorf("OSType [%s] is invalid", os)
}

func (api *ArchiveAPI) findByOSTags(tags []string, filterMap ...map[string]interface{}) (*sacloud.Archive, error) {

	api.Reset().WithTags(tags)

	for _, filters := range filterMap {
		for key, filter := range filters {
			api.FilterMultiBy(key, filter)
		}
	}
	res, err := api.Find()
	if err != nil {
		return nil, fmt.Errorf("Archive [%s] error : %s", strings.Join(tags, ","), err)
	}

	if len(res.Archives) == 0 {
		return nil, fmt.Errorf("Archive [%s] Not Found", strings.Join(tags, ","))
	}

	return &res.Archives[0], nil

}
