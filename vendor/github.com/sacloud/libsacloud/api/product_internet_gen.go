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

/************************************************
  generated by IDE. for [ProductInternetAPI]
************************************************/

import (
	"github.com/sacloud/libsacloud/sacloud"
)

/************************************************
   To support Setxxx interface for Find()
************************************************/

// SetEmpty 検索条件のリセット
func (api *ProductInternetAPI) SetEmpty() {
	api.reset()
}

// SetOffset オフセット
func (api *ProductInternetAPI) SetOffset(offset int) {
	api.offset(offset)
}

// SetLimit リミット
func (api *ProductInternetAPI) SetLimit(limit int) {
	api.limit(limit)
}

// SetInclude 取得する項目
func (api *ProductInternetAPI) SetInclude(key string) {
	api.include(key)
}

// SetExclude 除外する項目
func (api *ProductInternetAPI) SetExclude(key string) {
	api.exclude(key)
}

// SetFilterBy 指定キーでのフィルター
func (api *ProductInternetAPI) SetFilterBy(key string, value interface{}) {
	api.filterBy(key, value, false)
}

// SetFilterMultiBy 任意項目でのフィルタ(完全一致 OR条件)
func (api *ProductInternetAPI) SetFilterMultiBy(key string, value interface{}) {
	api.filterBy(key, value, true)
}

// SetNameLike 名称条件
func (api *ProductInternetAPI) SetNameLike(name string) {
	api.FilterBy("Name", name)
}

// SetTag タグ条件
func (api *ProductInternetAPI) SetTag(tag string) {
	api.FilterBy("Tags.Name", tag)
}

// SetTags タグ(複数)条件
func (api *ProductInternetAPI) SetTags(tags []string) {
	api.FilterBy("Tags.Name", []interface{}{tags})
}

// func (api *ProductInternetAPI) SetSizeGib(size int) {
// 	api.FilterBy("SizeMB", size*1024)
// }

// func (api *ProductInternetAPI) SetSharedScope() {
// 	api.FilterBy("Scope", "shared")
// }

// func (api *ProductInternetAPI) SetUserScope() {
// 	api.FilterBy("Scope", "user")
// }

// SetSortBy 指定キーでのソート
func (api *ProductInternetAPI) SetSortBy(key string, reverse bool) {
	api.sortBy(key, reverse)
}

// SetSortByName 名称でのソート
func (api *ProductInternetAPI) SetSortByName(reverse bool) {
	api.sortByName(reverse)
}

// func (api *ProductInternetAPI) SetSortBySize(reverse bool) {
// 	api.sortBy("SizeMB", reverse)
// }

/************************************************
   To support fluent interface for Find()
************************************************/

// Reset 検索条件のリセット
func (api *ProductInternetAPI) Reset() *ProductInternetAPI {
	api.reset()
	return api
}

// Offset オフセット
func (api *ProductInternetAPI) Offset(offset int) *ProductInternetAPI {
	api.offset(offset)
	return api
}

// Limit リミット
func (api *ProductInternetAPI) Limit(limit int) *ProductInternetAPI {
	api.limit(limit)
	return api
}

// Include 取得する項目
func (api *ProductInternetAPI) Include(key string) *ProductInternetAPI {
	api.include(key)
	return api
}

// Exclude 除外する項目
func (api *ProductInternetAPI) Exclude(key string) *ProductInternetAPI {
	api.exclude(key)
	return api
}

// FilterBy 指定キーでのフィルター
func (api *ProductInternetAPI) FilterBy(key string, value interface{}) *ProductInternetAPI {
	api.filterBy(key, value, false)
	return api
}

// FilterMultiBy 任意項目でのフィルタ(完全一致 OR条件)
func (api *ProductInternetAPI) FilterMultiBy(key string, value interface{}) *ProductInternetAPI {
	api.filterBy(key, value, true)
	return api
}

// WithNameLike 名称条件
func (api *ProductInternetAPI) WithNameLike(name string) *ProductInternetAPI {
	return api.FilterBy("Name", name)
}

// WithTag タグ条件
func (api *ProductInternetAPI) WithTag(tag string) *ProductInternetAPI {
	return api.FilterBy("Tags.Name", tag)
}

// WithTags タグ(複数)条件
func (api *ProductInternetAPI) WithTags(tags []string) *ProductInternetAPI {
	return api.FilterBy("Tags.Name", []interface{}{tags})
}

// func (api *ProductInternetAPI) WithSizeGib(size int) *ProductInternetAPI {
// 	api.FilterBy("SizeMB", size*1024)
// 	return api
// }

// func (api *ProductInternetAPI) WithSharedScope() *ProductInternetAPI {
// 	api.FilterBy("Scope", "shared")
// 	return api
// }

// func (api *ProductInternetAPI) WithUserScope() *ProductInternetAPI {
// 	api.FilterBy("Scope", "user")
// 	return api
// }

// SortBy 指定キーでのソート
func (api *ProductInternetAPI) SortBy(key string, reverse bool) *ProductInternetAPI {
	api.sortBy(key, reverse)
	return api
}

// SortByName 名称でのソート
func (api *ProductInternetAPI) SortByName(reverse bool) *ProductInternetAPI {
	api.sortByName(reverse)
	return api
}

// func (api *ProductInternetAPI) SortBySize(reverse bool) *ProductInternetAPI {
// 	api.sortBy("SizeMB", reverse)
// 	return api
// }

/************************************************
  To support CRUD(Create/Read/Update/Delete)
************************************************/

//func (api *ProductInternetAPI) New() *sacloud.ProductInternet {
// 	return &sacloud.ProductInternet{}
//}

// func (api *ProductInternetAPI) Create(value *sacloud.ProductInternet) (*sacloud.ProductInternet, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.create(api.createRequest(value), res)
// 	})
// }

// Read 読み取り
func (api *ProductInternetAPI) Read(id sacloud.ID) (*sacloud.ProductInternet, error) {
	return api.request(func(res *sacloud.Response) error {
		return api.read(id, nil, res)
	})
}

// func (api *ProductInternetAPI) Update(id sacloud.ID, value *sacloud.ProductInternet) (*sacloud.ProductInternet, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.update(id, api.createRequest(value), res)
// 	})
// }

// func (api *ProductInternetAPI) Delete(id sacloud.ID) (*sacloud.ProductInternet, error) {
// 	return api.request(func(res *sacloud.Response) error {
// 		return api.delete(id, nil, res)
// 	})
// }

/************************************************
  Inner functions
************************************************/

func (api *ProductInternetAPI) setStateValue(setFunc func(*sacloud.Request)) *ProductInternetAPI {
	api.baseAPI.setStateValue(setFunc)
	return api
}

func (api *ProductInternetAPI) request(f func(*sacloud.Response) error) (*sacloud.ProductInternet, error) {
	res := &sacloud.Response{}
	err := f(res)
	if err != nil {
		return nil, err
	}
	return res.InternetPlan, nil
}

func (api *ProductInternetAPI) createRequest(value *sacloud.ProductInternet) *sacloud.Request {
	req := &sacloud.Request{}
	req.InternetPlan = value
	return req
}
