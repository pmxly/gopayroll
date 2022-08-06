package common

const TimeLayout = "2006-01-02 15:04:05"
const StdLocation = "Asia/Chongqing"
const LocalLocation = "Asia/Chongqing"
const LocalEscLoc = "Asia%2FChongqing"

//max number of concurrent payroll being processed by this worker at a time
//单个worker一次处理的任务最大并发数，如果prefetch_count小于该值，则以prefetch_count为准
const Concurrency = 10