package vkapi

import "fmt"

type Error struct {
	ErrorCode  uint64 `json:"error_code"`
	ErrorMsg   string `json:"error_msg"`
	CaptchaSID string `json:"captcha_sid"`
	CaptchaImg string `json:"captcha_sid"`
}

func (err Error) Error() string {
	if err.CaptchaSID != "" {
		return fmt.Sprintf("vkapi: %s (%d); captcha_sid: %s, captcha_img: %s",
			err.ErrorMsg, err.ErrorCode, err.CaptchaSID, err.CaptchaImg)
	} else {
		return fmt.Sprintf("vkapi: %s (%d)", err.ErrorMsg, err.ErrorCode)
	}
}
