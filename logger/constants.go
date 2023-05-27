package logger

import "net/http"

const (
	KEY_PAYLOAD        string = "payload"
	KEY_ERROR          string = "error"
	KEY_MESSAGE        string = "message"
	KEY_FUNCTION       string = "function"
	KEY_PATH           string = "path"
	KEY_KEY            string = "key"
	KEY_TIME_TAKEN     string = "time_taken"
	KEY_S3_FILE        string = "s3_file"
	KEY_ITEMS_COUNT    string = "items_count"
	KEY_SQS_QUEUE_NAME string = "sqs_queue"
	KEY_EVENT          string = "event"
)

type ErrResp struct {
	Error_Msg string
	Code      string
}

//codes
const (
	INVALID_REQ = http.StatusBadRequest
	SERVER_ERR  = http.StatusInternalServerError
	SUCCESS     = http.StatusOK
	AUTH_ERR    = http.StatusUnauthorized
)

var (
	SUCCESS_RESP = ErrResp{
		Error_Msg: "SUCCESS",
		Code:      "P10",
	}

	UNAUTHORIZED_RESP = ErrResp{
		Error_Msg: "UNAUTHORIZED",
		Code:      "P20",
	}

	INVALID_REQUEST_PAYLOAD_RESP = ErrResp{
		Error_Msg: "INVALID_REQUEST",
		Code:      "P30",
	}

	SERVER_ERROR_RESP = ErrResp{
		Error_Msg: "SERVER_ERROR",
		Code:      "P40",
	}

	UNSUPPORTED_MESSAGE_RESP = ErrResp{
		Error_Msg: "UNSUPPORTED_MESSAGE",
		Code:      "P000",
	}
)

const (
	RedisError          = "REDIS_ERROR"
	StructToMap         = "GOLANG:STRUCT_TO_MAP"
	BindJSONtoStruct    = "GOLANG:BINDING_JSON_TO_STRUCT_FAILED"
	JSONParseError      = "JSON_PARSE_ERROR"
	JSONUnmarshallError = "JSON_UNMARSHALL_ERROR"
	NILValueError       = "NIL_VALUE_ERROR"
)
