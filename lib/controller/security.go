/*
 * Copyright 2020 InfAI (CC SES)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/SENERGY-Platform/import-deploy/lib/auth"
	"github.com/SENERGY-Platform/import-deploy/lib/model"
	"net/http"
	"net/url"
	"runtime/debug"
)

func (this *Controller) checkBool(jwt auth.Token, kind string, id string, action model.AuthAction) (allowed bool, err error) {
	if IsAdmin(jwt) {
		return true, nil
	}
	req, err := http.NewRequest("GET", this.config.PermissionsUrl+"/jwt/check/"+url.QueryEscape(kind)+"/"+url.QueryEscape(id)+"/"+action.String()+"/bool", nil)
	if err != nil {
		debug.PrintStack()
		return false, err
	}
	req.Header.Set("Authorization", jwt.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		debug.PrintStack()
		return false, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		return false, errors.New(buf.String())
	}
	err = json.NewDecoder(resp.Body).Decode(&allowed)
	if err != nil {
		debug.PrintStack()
		return false, err
	}
	return allowed, nil
}

func IsAdmin(jwt auth.Token) bool {
	return contains(jwt.RealmAccess.Roles, "admin")
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
