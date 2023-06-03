/*
xpan

xpanapi

API version: 0.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// Quotaresponse struct for Quotaresponse
type Quotaresponse struct {
	Errno     *int32 `json:"errno,omitempty"`
	Total     *int64 `json:"total,omitempty"`
	Free      *int64 `json:"free,omitempty"`
	RequestId *int64 `json:"request_id,omitempty"`
	Expire    *bool  `json:"expire,omitempty"`
	Used      *int64 `json:"used,omitempty"`
}

// NewQuotaresponse instantiates a new Quotaresponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewQuotaresponse() *Quotaresponse {
	this := Quotaresponse{}
	return &this
}

// NewQuotaresponseWithDefaults instantiates a new Quotaresponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewQuotaresponseWithDefaults() *Quotaresponse {
	this := Quotaresponse{}
	return &this
}

// GetErrno returns the Errno field value if set, zero value otherwise.
func (o *Quotaresponse) GetErrno() int32 {
	if o == nil || o.Errno == nil {
		var ret int32
		return ret
	}
	return *o.Errno
}

// GetErrnoOk returns a tuple with the Errno field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Quotaresponse) GetErrnoOk() (*int32, bool) {
	if o == nil || o.Errno == nil {
		return nil, false
	}
	return o.Errno, true
}

// HasErrno returns a boolean if a field has been set.
func (o *Quotaresponse) HasErrno() bool {
	if o != nil && o.Errno != nil {
		return true
	}

	return false
}

// SetErrno gets a reference to the given int32 and assigns it to the Errno field.
func (o *Quotaresponse) SetErrno(v int32) {
	o.Errno = &v
}

// GetTotal returns the Total field value if set, zero value otherwise.
func (o *Quotaresponse) GetTotal() int64 {
	if o == nil || o.Total == nil {
		var ret int64
		return ret
	}
	return *o.Total
}

// GetTotalOk returns a tuple with the Total field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Quotaresponse) GetTotalOk() (*int64, bool) {
	if o == nil || o.Total == nil {
		return nil, false
	}
	return o.Total, true
}

// HasTotal returns a boolean if a field has been set.
func (o *Quotaresponse) HasTotal() bool {
	if o != nil && o.Total != nil {
		return true
	}

	return false
}

// SetTotal gets a reference to the given int64 and assigns it to the Total field.
func (o *Quotaresponse) SetTotal(v int64) {
	o.Total = &v
}

// GetFree returns the Free field value if set, zero value otherwise.
func (o *Quotaresponse) GetFree() int64 {
	if o == nil || o.Free == nil {
		var ret int64
		return ret
	}
	return *o.Free
}

// GetFreeOk returns a tuple with the Free field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Quotaresponse) GetFreeOk() (*int64, bool) {
	if o == nil || o.Free == nil {
		return nil, false
	}
	return o.Free, true
}

// HasFree returns a boolean if a field has been set.
func (o *Quotaresponse) HasFree() bool {
	if o != nil && o.Free != nil {
		return true
	}

	return false
}

// SetFree gets a reference to the given int64 and assigns it to the Free field.
func (o *Quotaresponse) SetFree(v int64) {
	o.Free = &v
}

// GetRequestId returns the RequestId field value if set, zero value otherwise.
func (o *Quotaresponse) GetRequestId() int64 {
	if o == nil || o.RequestId == nil {
		var ret int64
		return ret
	}
	return *o.RequestId
}

// GetRequestIdOk returns a tuple with the RequestId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Quotaresponse) GetRequestIdOk() (*int64, bool) {
	if o == nil || o.RequestId == nil {
		return nil, false
	}
	return o.RequestId, true
}

// HasRequestId returns a boolean if a field has been set.
func (o *Quotaresponse) HasRequestId() bool {
	if o != nil && o.RequestId != nil {
		return true
	}

	return false
}

// SetRequestId gets a reference to the given int64 and assigns it to the RequestId field.
func (o *Quotaresponse) SetRequestId(v int64) {
	o.RequestId = &v
}

// GetExpire returns the Expire field value if set, zero value otherwise.
func (o *Quotaresponse) GetExpire() bool {
	if o == nil || o.Expire == nil {
		var ret bool
		return ret
	}
	return *o.Expire
}

// GetExpireOk returns a tuple with the Expire field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Quotaresponse) GetExpireOk() (*bool, bool) {
	if o == nil || o.Expire == nil {
		return nil, false
	}
	return o.Expire, true
}

// HasExpire returns a boolean if a field has been set.
func (o *Quotaresponse) HasExpire() bool {
	if o != nil && o.Expire != nil {
		return true
	}

	return false
}

// SetExpire gets a reference to the given bool and assigns it to the Expire field.
func (o *Quotaresponse) SetExpire(v bool) {
	o.Expire = &v
}

// GetUsed returns the Used field value if set, zero value otherwise.
func (o *Quotaresponse) GetUsed() int64 {
	if o == nil || o.Used == nil {
		var ret int64
		return ret
	}
	return *o.Used
}

// GetUsedOk returns a tuple with the Used field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Quotaresponse) GetUsedOk() (*int64, bool) {
	if o == nil || o.Used == nil {
		return nil, false
	}
	return o.Used, true
}

// HasUsed returns a boolean if a field has been set.
func (o *Quotaresponse) HasUsed() bool {
	if o != nil && o.Used != nil {
		return true
	}

	return false
}

// SetUsed gets a reference to the given int64 and assigns it to the Used field.
func (o *Quotaresponse) SetUsed(v int64) {
	o.Used = &v
}

func (o Quotaresponse) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Errno != nil {
		toSerialize["errno"] = o.Errno
	}
	if o.Total != nil {
		toSerialize["total"] = o.Total
	}
	if o.Free != nil {
		toSerialize["free"] = o.Free
	}
	if o.RequestId != nil {
		toSerialize["request_id"] = o.RequestId
	}
	if o.Expire != nil {
		toSerialize["expire"] = o.Expire
	}
	if o.Used != nil {
		toSerialize["used"] = o.Used
	}
	return json.Marshal(toSerialize)
}

type NullableQuotaresponse struct {
	value *Quotaresponse
	isSet bool
}

func (v NullableQuotaresponse) Get() *Quotaresponse {
	return v.value
}

func (v *NullableQuotaresponse) Set(val *Quotaresponse) {
	v.value = val
	v.isSet = true
}

func (v NullableQuotaresponse) IsSet() bool {
	return v.isSet
}

func (v *NullableQuotaresponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableQuotaresponse(val *Quotaresponse) *NullableQuotaresponse {
	return &NullableQuotaresponse{value: val, isSet: true}
}

func (v NullableQuotaresponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableQuotaresponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
