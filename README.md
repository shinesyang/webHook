
`webhook用于grafana的飞书与微信企业应用号告警`

1. webhook需要配合grafana配置才能使用，grafana的配置方法参考: https://www.jianshu.com/p/cb2f2f4861d5?v=1701244595432
2. webhook使用方法如下:
>1. 修改配置config/config
>2. webhook三个接口可自行选择
3. webhook接口
>1. /webhook/receive       告警飞书/微信
>2. /webhook/feishu        告警飞书
>3. /webhook/wechat        告警微信
4. webhook默认验证账号: `monitor_webhook/gx52cGjRVsck`
