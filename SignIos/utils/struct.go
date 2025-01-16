package utils

type XGorGonInfo struct {
	QueryMd5      string `json:"query_md5"`
	BodyMd5OrStub string `json:"body_md5_or_stub"`
	Unused        string `json:"unused"`
	SdkVersion    string `json:"'sdk_version'"`
	XKhronos      int64  `json:"x-khronos"`
}

// Host: api5-normal-c-hl.amemv.com
// Connection: keep-alive
// Content-Length: 74
// Cookie: odin_tt=3fdc0327a8a11fb3915a3687e1d760e660b5c864f832c289ecc4b338d4ccb38ddf18ca8afb85fe56013548bff204b6c0c3333cf0017aa95a5e5480b349ea531e; passport_csrf_token=91e49cff5af1888d4c45ed15774f8ff6; passport_csrf_token_default=91e49cff5af1888d4c45ed15774f8ff6; install_id=3965299072493930; ttreq=1$42572cfd3f56ac3c7fbbbeab8cf533572867484b
// x-vc-bdturing-sdk-version: 2.2.3
// Content-Type: application/x-www-form-urlencoded
// X-SS-Cookie: install_id=3965299072493930; ttreq=1$42572cfd3f56ac3c7fbbbeab8cf533572867484b; passport_csrf_token=91e49cff5af1888d4c45ed15774f8ff6; passport_csrf_token_default=91e49cff5af1888d4c45ed15774f8ff6; odin_tt=3fdc0327a8a11fb3915a3687e1d760e660b5c864f832c289ecc4b338d4ccb38ddf18ca8afb85fe56013548bff204b6c0c3333cf0017aa95a5e5480b349ea531e
// tt-request-time: 1722957691645
// User-Agent: AwemeDS 20.0.1 rv: 200013 (iPhone; iOS 13.5.1; zh_CN) Cronet
// x-tt-passport-csrf-token: 91e49cff5af1888d4c45ed15774f8ff6
// sdk-version: 2
// x-tt-dt: AAAZUOP5FRGTM436K454443ULEP4PV7KNLB7OWQJ7XOCR5F6NJH3VEVURFHB36OEAZUPW33P6ABLDE3DS75J6O7FSN2J6HSDN3AWDWVDXWUEP5W45ITJIFGCWYABOCJ4OFRJPWKTR3KZRENPZRKXBPA
// passport-sdk-version: 5.12.1
// X-SS-STUB: 7858C10C8B83A7CDE2E4963C6D991399
// X-SS-DP: 1349
// x-tt-trace-id: 00-2847fc930d4666b21f926d9ff6c30545-2847fc930d4666b2-01
// Accept-Encoding: gzip, deflate
// X-Argus: faHEg05wbE8JtVckdM2WChfjHJx9k7Xy3SXVDfF5X06kiVZSGW/cQxY1q78T6cUQ5Pva7TDuSr2sD4b22jD0u5PTqSbh0QxkiWUemhiS5u4kbETYHd/1F9uQINi9VFa2seCHcLudTUaBzOpWBG+X7rfy0QbPmu5gd+rSJjNaKq1NqsAEv60Oi2ZnQxKoeqqv/SqvsfZpXLa88ytHLUU5A9f0s2Sb7pP+rwRjOkW96lRi3d23OZQlHYf9yIrSd/OWBR8=
// X-Gorgon: 8404007920017a904a6fcf3fd2300be4d98277e430e5a2563d89
// X-Khronos: 1722957691
// X-Ladon: t0zKVsI2zdZVVcDekaTjwVR0P/HDDC/0x8olYW1gEaCr8E2f
