package vkapi

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Error", func() {
	var _ error = Error{}

	It("should return error description", func() {
		Expect(Error{
			ErrorCode: 42,
			ErrorMsg:  "Test error",
		}).To(MatchError("vkapi: Test error (42)"))
	})

	It("should return captcha error description", func() {
		Expect(Error{
			ErrorCode:  14,
			ErrorMsg:   "Captcha needed",
			CaptchaSID: "123456789012",
			CaptchaImg: "http://api.vk.com/captcha.php?sid=123456789012&s=1",
		}).To(MatchError(
			"vkapi: Captcha needed (14); captcha_sid: 123456789012, captcha_img: http://api.vk.com/captcha.php?sid=123456789012&s=1",
		))
	})

})
