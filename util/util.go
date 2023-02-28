package util

import (
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/auth/credentials"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"math/rand"
	"time"
)

var HTML = "<div style=\"line-height:1.7;color:#000000;font-size:14px;font-family:Arial\"><div><div style=\"width: 100%;line-height:1.7;color:#000000;font-size:14px;font-family:Arial\">\n      <div style=\"width: 100%;height: 80px;display: flex;justify-content: flex-startalign-content: center;;align-items: center;background: linear-gradient(45deg, #E942B4 0%, #1E46D4 100%);\" class=\"color-header\">\n        < img src=\"https://array-fi.oss-ap-southeast-1.aliyuncs.com/common/logo.png\" alt=\"logo\" style=\"height: 40px;cursor: pointer;margin-left: 35px;\" />\n      </div>\n\n      <div class=\"content\" style=\"font-size: 20px;text-align: left; line-height: 26px;padding: 30px 35px;\">\n        <p>Verify Email</p >\n        <p>You are verifying your email address with the verification code:</p >\n        <p style=\"font-size:20px; font-weight:600\">968455</p >\n        <p>\n          The verification code is valid for 10 minutes, please do not disclose it. If this is not\n          done by you, please contact the customer service team immediately.\n        </p >\n        <p>Array团队</p >\n      </div>\n\n      <div class=\"footer\" style=\"line-height: 16px;text-align: center;margin-top: 100px;\">\n        <p>Array is committed to protecting the security of your account and transactions</p >\n        <p>- If you suspect that you have received a fraudulent message, please contact Customer Service immediately - </p >\n        <p>If you are in doubt about the authenticity of the message, do not hesitate to verify the message through official channels</p >\n        <p>- Do not share your verification code with anyone, including official customer service and staff</p >\n      </div>\n    </div>\n\n\n\n\n\n\n<br /></div><br /><br /></div>"

func EmailSign(email string, code string) {
	config := sdk.NewConfig()

	credential := credentials.NewAccessKeyCredential("LTAI5tPQHGSx62BmtQtXo5q4", "9Z3Bc6tfyCQ9S6ACABJjyI0PtpKE9t")
	/* use STS Token
	credential := credentials.NewStsTokenCredential("<your-access-key-id>", "<your-access-key-secret>", "<your-sts-token>")
	*/
	client, err := sdk.NewClientWithOptions("cn-hangzhou", config, credential)
	if err != nil {
		panic(err)
	}

	request := requests.NewCommonRequest()
	//发送邮箱
	request = SendMessage(email, code, request)

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())
}

func SendMessage(mail string, code string, request *requests.CommonRequest) *requests.CommonRequest {
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dm.aliyuncs.com"
	request.Version = "2015-11-23"
	request.ApiName = "SingleSendMail"
	request.QueryParams["AddressType"] = "1"
	request.QueryParams["AccountName"] = "admin@arrayemail.com"
	request.QueryParams["TagName"] = "1111"
	request.QueryParams["ToAddress"] = mail
	request.QueryParams["Subject"] = "Array: Verify Your Email Address"
	//request.QueryParams["TextBody"] = "验证码测试"
	request.QueryParams["FromAlias"] = "ArrayService"
	request.QueryParams["ReplyToAddress"] = "true"
	//fmt.Printf(fmt.Sprintf(HTML,code))
	request.QueryParams["HtmlBody"] = "<div>\n  <div\n    style=\"width: 100%;line-height: 1.7; " +
		"color: #000000; font-size: 14px; font-family: Arial; display: flex; " +
		"justify-content: center; align-items: center;\"\n  >\n    <div style=\"width:660px; " +
		"background-color: #181818;border-radius:6px;\">\n     " +
		" <div\n        style=\"width: 100%; line-height: 1.7; color: #FFF; " +
		"font-size: 14px; font-family: Arial\"\n      >\n        " +
		"<div\n          style=\"width: 100%;height: 80px;display: flex;justify-content: " +
		"flex-startalign-content: center;;align-items: center;background: linear-gradient(45deg, #E942B4 0%, #1E46D4 100%);" +
		"border-radius:6px 6px 0 0; \"\n        >\n          " +
		"<img\n            src=\"https://array-fi.oss-ap-southeast-1.aliyuncs.com/common/logo.png\"\n            " +
		"alt=\"logo\"\n            style=\"height: 40px; cursor: pointer; margin-left: 35px\"\n          />\n       " +
		" </div>\n\n        <div\n          class=\"content\"\n          style=\"font-size: 20px; text-align: " +
		"left; line-height: 26px; padding: 30px 35px\"\n        >\n          <p>Verify Email</p>\n          " +
		"<p>You are verifying your email address with the verification code:</p>          " +
		"<p style=\"font-size: 20px; font-weight: 600\">" + code + "</p>\n=          <p>\n            " +
		"The verification code is valid for 10 minutes, please do not disclose it. If this is not\n           " +
		" done by you, please contact the customer service team immediately.\n          </p>\n          " +
		"<p>Array Team</p>\n        </div>\n\n        <div class=\"footer\" style=\"line-height: 16px;" +
		" text-align: center; margin-top: 100px; padding-bottom: 30px;\">\n         " +
		" <p>Array is committed to protecting the security of your account and transactions</p>\n          " +
		"<p>\n            - If you suspect that you have received a fraudulent message, please contact Customer\n            " +
		"Service immediately -\n          </p>\n          " +
		"<p>\n            If you are in doubt about the authenticity of the message, do not hesitate to verify the\n          " +
		"  message through official channels\n          </p>\n          <p>\n            " +
		"- Do not share your verification code with anyone, including official customer service\n           " +
		" and staff\n          </p>\n        </div>\n      </div>\n    </div>\n\n    <br />\n\n   " +
		" <br /><br />\n  </div>\n  <br />\n</div>\n"

	return request
}

// RandomString 生成随机数验证码
/**
* size 随机码的位数
* kind 0    // 纯数字
       1    // 小写字母
       2    // 大写字母
       3    // 数字、大小写字母
*/
func RandomString(size int, kind int) string {
	ikind, kinds, rsbytes := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	isAll := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if isAll { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		rsbytes[i] = uint8(base + rand.Intn(scope))
	}
	return string(rsbytes)
}
