# 羽毛球场抢票脚本简单教程

## 概览

![image-20231226183355147](https://raw.githubusercontent.com/drinkwateronly/Image-Host/main/image/image-20231226183355147.png)

本羽毛球场抢票脚本具备如下特点：

- 抢票成功率非常高，但不保证绝对成功。
- 预约场地不可指定，随机预约。
- 不需要校园卡密码即可抢票，无任何非法获取、存储信息的行为。
- 本脚本免费，请勿用于盈利或大规模传播。
- 为保证脚本的正确使用，请详细阅读下面的步骤和注意事项。





## 抢票步骤

1. 按照窗口提示填写如下信息
   - 预约人学号、姓名、手机号。
   - 选择场馆，填写预约年月日，格式为`YYYY-MM-DD`，例如`2024-01-01`。选择预约时间段。
   - 获取并填写cookie，cookie的获取方式放在最后。
2. 验证cookie有效性，若无效，需要重新获取cookie。有效时才能抢票。
3. 在放票的时刻点击开始，等待抢票结果。
4. 抢票结果将在日志框展示。
5. 若失败，可点击开始再次抢票。



## 注意事项

- 检查填写的预约日期的格式的正确性，并自行判断预约时间段是否会被体育课占用。
- 经过测试，手机号可以随便填一个真实号码，最好是自己在用的号码，以防被识别出脚本。
- cookie以`_WEU=`开头，复制到脚本输入框时，请注意cookie的开头和结尾不要有空格或回车，否则验证时会失效。
- cookie具有一定的时效性，不要过早获取，建议开抢前5分钟获取并验证。
- 可以预约一个当天的时间段以熟悉整个抢票过程，也可以验证信息有没有填对。

## 如何快速获取cookie

- 进入深圳大学体育场馆预约平台并登录。

  > https://ehall.szu.edu.cn/qljfwapp/sys/lwSzuCgyy/index.do#/sportVenue

- 点击键盘F12进入浏览器开发者模式。

- 进入开发者模式后，按如下步骤进行。

  ![image-20231222201126688](https://raw.githubusercontent.com/drinkwateronly/Image-Host/main/image/image-20231222201126688.png)

  ![image-20231222201244546](https://raw.githubusercontent.com/drinkwateronly/Image-Host/main/image/image-20231222201244546.png)

  ![image-20231226191303980](https://raw.githubusercontent.com/drinkwateronly/Image-Host/main/image/image-20231226191303980.png)

  

