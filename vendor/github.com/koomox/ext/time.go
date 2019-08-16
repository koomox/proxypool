package ext

import "time"

const (
	customTimeFormat = "2006-01-02 15:04:05"
)

/*
 * GetTimeNowUTC
 * @Return UTC时间的字符串时间
 */
func TimeNowUTC() (t string, err error) {
	tn := time.Now()
	utc, err := time.LoadLocation("UTC")
	if err != nil {
		return
	}
	return tn.In(utc).Format(customTimeFormat), nil
}

func TimeNowCST() (t string, err error) {
	tn := time.Now()
	cst, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return
	}
	return tn.In(cst).Format(customTimeFormat), nil
}

// Not Expires Return true or false
func TimeIsNotExpires(ct string) bool {
	t, _ := TimeNowUTC()
	if t > ct {
		return false // Expires
	}
	return true // Not Expires
}
