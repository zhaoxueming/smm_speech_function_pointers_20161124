package share

import (
    "log"
)

const (

    ERROR_SYSTEM                            = "系统错误！"
    ERROR_PARAM_FORMAT                      = "参数格式错误！"
    ERROR_DATE_FORMAT                       = "时间格式错误！"
    ERROR_GOREQUEST                         = "服务访问错误！"
    ERROR_SERVICE_REQUEST_CODE              = "服务访问失败！"
    ERROR_SERVICE_REQUEST_DATA              = "服务访问数据格式错误！"

    ERROR_FILE_UPLOAD_NAME                  = "获取文件名称失败！"
    ERROR_FILE_UPLOAD_BYTES                 = "获取文件内容数据失败！"
    ERROR_FILE_UPLOAD_TO_QINIU              = "上传文件失败！"
)

func init(){
    log.SetFlags(log.Lshortfile | log.LstdFlags)
}
