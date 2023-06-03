/*
xpan

xpanapi

API version: 0.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"bytes"
	_context "context"
	_ioutil "io/ioutil"
	_nethttp "net/http"
	_neturl "net/url"
)

// Linger please
var (
	_ _context.Context
)

// MultimediafileApiService MultimediafileApi service
type MultimediafileApiService service

type ApiXpanfilelistallRequest struct {
	ctx         _context.Context
	ApiService  *MultimediafileApiService
	accessToken *string
	path        *string
	recursion   *int32
	web         *string
	start       *int32
	limit       *int32
	order       *string
	desc        *int32
}

func (r ApiXpanfilelistallRequest) AccessToken(accessToken string) ApiXpanfilelistallRequest {
	r.accessToken = &accessToken
	return r
}
func (r ApiXpanfilelistallRequest) Path(path string) ApiXpanfilelistallRequest {
	r.path = &path
	return r
}
func (r ApiXpanfilelistallRequest) Recursion(recursion int32) ApiXpanfilelistallRequest {
	r.recursion = &recursion
	return r
}
func (r ApiXpanfilelistallRequest) Web(web string) ApiXpanfilelistallRequest {
	r.web = &web
	return r
}
func (r ApiXpanfilelistallRequest) Start(start int32) ApiXpanfilelistallRequest {
	r.start = &start
	return r
}
func (r ApiXpanfilelistallRequest) Limit(limit int32) ApiXpanfilelistallRequest {
	r.limit = &limit
	return r
}
func (r ApiXpanfilelistallRequest) Order(order string) ApiXpanfilelistallRequest {
	r.order = &order
	return r
}
func (r ApiXpanfilelistallRequest) Desc(desc int32) ApiXpanfilelistallRequest {
	r.desc = &desc
	return r
}

func (r ApiXpanfilelistallRequest) Execute() (string, *_nethttp.Response, error) {
	return r.ApiService.XpanfilelistallExecute(r)
}

/*
Xpanfilelistall Method for Xpanfilelistall

listall

	@param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiXpanfilelistallRequest
*/
func (a *MultimediafileApiService) Xpanfilelistall(ctx _context.Context) ApiXpanfilelistallRequest {
	return ApiXpanfilelistallRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
//
//	@return string
func (a *MultimediafileApiService) XpanfilelistallExecute(r ApiXpanfilelistallRequest) (string, *_nethttp.Response, error) {
	var (
		localVarHTTPMethod  = _nethttp.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue string
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MultimediafileApiService.Xpanfilelistall")
	if err != nil {
		return localVarReturnValue, nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/rest/2.0/xpan/multimedia?method=listall&openapi=xpansdk"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}
	if r.accessToken == nil {
		return localVarReturnValue, nil, reportError("accessToken is required and must be specified")
	}
	if r.path == nil {
		return localVarReturnValue, nil, reportError("path is required and must be specified")
	}
	if r.recursion == nil {
		return localVarReturnValue, nil, reportError("recursion is required and must be specified")
	}

	localVarQueryParams.Add("access_token", parameterToString(*r.accessToken, ""))
	localVarQueryParams.Add("path", parameterToString(*r.path, ""))
	localVarQueryParams.Add("recursion", parameterToString(*r.recursion, ""))
	if r.web != nil {
		localVarQueryParams.Add("web", parameterToString(*r.web, ""))
	}
	if r.start != nil {
		localVarQueryParams.Add("start", parameterToString(*r.start, ""))
	}
	if r.limit != nil {
		localVarQueryParams.Add("limit", parameterToString(*r.limit, ""))
	}
	if r.order != nil {
		localVarQueryParams.Add("order", parameterToString(*r.order, ""))
	}
	if r.desc != nil {
		localVarQueryParams.Add("desc", parameterToString(*r.desc, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json; charset=UTF-8"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = _ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiXpanmultimediafilemetasRequest struct {
	ctx         _context.Context
	ApiService  *MultimediafileApiService
	accessToken *string
	fsids       *string
	thumb       *string
	extra       *string
	dlink       *string
	path        *string
	needmedia   *int32
}

func (r ApiXpanmultimediafilemetasRequest) AccessToken(accessToken string) ApiXpanmultimediafilemetasRequest {
	r.accessToken = &accessToken
	return r
}
func (r ApiXpanmultimediafilemetasRequest) Fsids(fsids string) ApiXpanmultimediafilemetasRequest {
	r.fsids = &fsids
	return r
}
func (r ApiXpanmultimediafilemetasRequest) Thumb(thumb string) ApiXpanmultimediafilemetasRequest {
	r.thumb = &thumb
	return r
}
func (r ApiXpanmultimediafilemetasRequest) Extra(extra string) ApiXpanmultimediafilemetasRequest {
	r.extra = &extra
	return r
}

// 下载地址。下载接口需要用到dlink。
func (r ApiXpanmultimediafilemetasRequest) Dlink(dlink string) ApiXpanmultimediafilemetasRequest {
	r.dlink = &dlink
	return r
}

// 查询共享目录或专属空间内文件时需要。共享目录格式： /uk-fsid（其中uk为共享目录创建者id， fsid对应共享目录的fsid）。专属空间格式：/_pcs_.appdata/xpan/。
func (r ApiXpanmultimediafilemetasRequest) Path(path string) ApiXpanmultimediafilemetasRequest {
	r.path = &path
	return r
}
func (r ApiXpanmultimediafilemetasRequest) Needmedia(needmedia int32) ApiXpanmultimediafilemetasRequest {
	r.needmedia = &needmedia
	return r
}

func (r ApiXpanmultimediafilemetasRequest) Execute() (string, *_nethttp.Response, error) {
	return r.ApiService.XpanmultimediafilemetasExecute(r)
}

/*
Xpanmultimediafilemetas Method for Xpanmultimediafilemetas

multimedia filemetas

	@param ctx _context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiXpanmultimediafilemetasRequest
*/
func (a *MultimediafileApiService) Xpanmultimediafilemetas(ctx _context.Context) ApiXpanmultimediafilemetasRequest {
	return ApiXpanmultimediafilemetasRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
//
//	@return string
func (a *MultimediafileApiService) XpanmultimediafilemetasExecute(r ApiXpanmultimediafilemetasRequest) (string, *_nethttp.Response, error) {
	var (
		localVarHTTPMethod  = _nethttp.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue string
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "MultimediafileApiService.Xpanmultimediafilemetas")
	if err != nil {
		return localVarReturnValue, nil, GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/rest/2.0/xpan/multimedia?method=filemetas&openapi=xpansdk"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := _neturl.Values{}
	localVarFormParams := _neturl.Values{}
	if r.accessToken == nil {
		return localVarReturnValue, nil, reportError("accessToken is required and must be specified")
	}
	if r.fsids == nil {
		return localVarReturnValue, nil, reportError("fsids is required and must be specified")
	}

	localVarQueryParams.Add("access_token", parameterToString(*r.accessToken, ""))
	localVarQueryParams.Add("fsids", parameterToString(*r.fsids, ""))
	if r.thumb != nil {
		localVarQueryParams.Add("thumb", parameterToString(*r.thumb, ""))
	}
	if r.extra != nil {
		localVarQueryParams.Add("extra", parameterToString(*r.extra, ""))
	}
	if r.dlink != nil {
		localVarQueryParams.Add("dlink", parameterToString(*r.dlink, ""))
	}
	if r.path != nil {
		localVarQueryParams.Add("path", parameterToString(*r.path, ""))
	}
	if r.needmedia != nil {
		localVarQueryParams.Add("needmedia", parameterToString(*r.needmedia, ""))
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json; UTF-8"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := _ioutil.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = _ioutil.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}