package errors

const (
	ENone                      = 0x00010000
	ERecordCreateFailed        = 0x00011001 /*mysql*/
	ERecordNotFound            = 0x00011002
	ERecordUpdateFailed        = 0x00011003
	ERecordFindFailed          = 0x00011004
	EDbInitFailed              = 0x00011005
	ERecordSaveFailed          = 0x00011006
	ERecordDeleteFailed        = 0x00011007
	EJsonpbMarshalFailed       = 0x00012001 /*json*/
	EJsonpbUnmarshalFailed     = 0x00012002
	ESimpleJsonMarshalFailed   = 0x00012003
	ESimpleJsonUnmarshalFailed = 0x00012004
	ESimpleJsonAssertFailed    = 0x00012005
	EJsonMarshalFailed         = 0x00012006
	EJsonUnmarshalFailed       = 0x00012007
	EHttpGetFailed             = 0x00013001 /*http*/
	EHttpPostFailed            = 0x00013002
	EHttpRequestEmpty          = 0x00013003
	EHttpRequestTimeOut        = 0x00013004
	EIoutilReadAllFailed       = 0x00014001 /*io*/
	EParamsNotSetError         = 0x00015001 /*params*/
	EParamsValueError          = 0x00015002
	EStrconvFailed             = 0x00016001 /*strconv*/
	EPublishFailed             = 0x00017001 /*broker*/
	EBase64DecodeFailed        = 0x00018001 /*encode decode*/
	EESSearchFailed            = 0x00019001 /*es*/
	EUtilPb2Map                = 0x00020001 /*util*/
)
