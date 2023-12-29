// 编译可执行文件命令：go build -ldflags "-s -w -H=windowsgui" .

package main

import (
	"GUIdemo/theme"
	"encoding/json"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// generateInputContainer 生成带标题title的输入框，并带有placeHolder
func generateInputContainer(title, placeHolder string) (*widget.Entry, *fyne.Container) {
	inputEntry := widget.NewEntry()
	inputEntry.SetPlaceHolder(placeHolder)
	inputTitle := canvas.NewText(title, color.White)
	return inputEntry, container.NewGridWithRows(2, inputTitle, inputEntry)
}

// ValidateCookie 验证cookie的有效性
func ValidateCookie(cookie string) (error, bool) {
	openRoomUrl := "https://ehall.szu.edu.cn/qljfwapp/sys/lwSzuCgyy/modules/sportVenue/getOpeningRoom.do"
	reqBody := fmt.Sprintf("XMDM=001&YYRQ=%s&YYLX=1.0&KSSJ=%s%%3A00&JSSJ=%s%%3A00&XQDM=1", time.Now().Format("2006-01-02"), "11", "12")
	req, err := http.NewRequest("POST", openRoomUrl, strings.NewReader(reqBody))
	if err != nil {
		return err, false
	}
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return err, false
	}
	if response.StatusCode != 200 {
		return err, false
	}
	return nil, true
}

// GetOpenRooms 获取所有空闲的场地的WID
func GetOpenRooms(cookie, date, beginHour, endHour string) (error, [][]string) {
	log.Println(cookie, date, beginHour, endHour)
	openRooms := make([][]string, 0)
	openRoomUrl := "https://ehall.szu.edu.cn/qljfwapp/sys/lwSzuCgyy/modules/sportVenue/getOpeningRoom.do"
	reqBody := fmt.Sprintf("XMDM=001&YYRQ=%s&YYLX=1.0&KSSJ=%s%%3A00&JSSJ=%s%%3A00&XQDM=1", date, beginHour, endHour)
	req, _ := http.NewRequest("POST", openRoomUrl, strings.NewReader(reqBody))
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return err, nil
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return err, nil
	}
	// 判断响应码
	if response == nil { // 出现过response为空的情况
		return errors.New("未知错误，请联系软件开发者解决"), nil
	}
	respCode := response.StatusCode
	if respCode == 401 {
		log.Println(err)
		return errors.New("cookie失效"), nil
	}
	if respCode != 200 {
		return errors.New("未知错误，请联系软件开发者解决"), nil
	}
	// 从响应body取出所有的场次信息
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err, nil
	}
	allField := result["datas"].(map[string]interface{})["getOpeningRoom"].(map[string]interface{})["rows"].([]interface{})
	// 筛选出空闲场次的WID + 场次名称，保存为string列表openRooms
	for i := range allField {
		field := allField[i].(map[string]interface{})
		if !field["disabled"].(bool) {
			openRooms = append(openRooms, []string{field["WID"].(string), field["CDMC"].(string)})
		}
	}
	return nil, openRooms
}

// isAnyStringEmpty 判断是否有任意字符串为空，用于校验所有信息是否填写
func isAnyStringEmpty(strings ...string) bool {
	for _, str := range strings {
		if str == "" {
			return true
		}
	}
	return false
}

// checkZeroField 检查BasicInfo结构体有无空字段
func checkZeroField(b BasicInfo) bool {
	if isAnyStringEmpty(b.phone, b.name, b.id, b.cookie, b.sport, b.beginHour, b.endHour, b.date) {
		return true
	}
	return false
}

// AddOrderByInfo 根据填写信息下单，bool为真时下单成功
func AddOrderByInfo(name, id, phone, cookie, date, beginHour, endHour, WID string) error {
	var result map[string]interface{}
	// 请求
	bookUrl := "https://ehall.szu.edu.cn/qljfwapp/sys/lwSzuCgyy/sportVenue/insertVenueBookingInfo.do"
	reqBody := fmt.Sprintf("DHID=&YYRGH=%s&CYRS=&YYRXM=%s&LXFS=%s&CGDM=001&CDWID=%s&XMDM=001&XQWID=1&KYYSJD=%s%%3A00-%s%%3A00&YYRQ=%s&YYLX=1.0&YYKS=%s+%s%%3A00&YYJS=%s+%s%%3A00&PC_OR_PHONE=pc",
		id, name, phone, WID, beginHour, endHour, date, date, beginHour, date, endHour)
	req, _ := http.NewRequest("POST", bookUrl, strings.NewReader(reqBody))
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	// 响应
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return err
	}
	respCode := result["code"].(string)
	msg := result["msg"].(string)
	if respCode == "0" && msg == "成功" {
		return nil
	}
	return errors.New(msg)
}

// AddOrderWorker 根据信息完成预约场地，包括查询空闲场地、根据完整信息下 单、输出结果。
func AddOrderWorker(basicInfo BasicInfo, logOutput *widget.Entry) {
	name := basicInfo.name
	id := basicInfo.id
	phone := basicInfo.phone
	cookie := basicInfo.cookie
	date := basicInfo.date
	beginHour := basicInfo.beginHour
	endHour := basicInfo.endHour
	// 请求获取预约场次的空闲场地
	err, openRooms := GetOpenRooms(cookie, date, beginHour, endHour)
	if err != nil {
		logOutput.Append(NewTextWithPrefix("预约失败：获取空闲场次出错，" + err.Error()))
		return
	}
	// 无空闲场次
	if len(openRooms) == 0 {
		logOutput.Append(NewTextWithPrefix("预约失败：当前时段暂无空闲场次"))
		return
	}
	logOutput.Append(NewTextWithPrefix("当前空闲场次数量:" + strconv.Itoa(len(openRooms))))
	logOutput.Append(NewTextWithPrefix("开始预约"))
	// 随机从空闲场地随机选择一个场地
	var randIndex int
	if len(openRooms)-1 == 0 {
		randIndex = 0
	} else {
		randIndex = rand.Intn(len(openRooms) - 1) // Intn的输入小于等于0时panic
	}
	WID := openRooms[randIndex][0]
	// 根据场地下单
	err = AddOrderByInfo(name, id, phone, cookie, date, beginHour, endHour, WID)
	if err != nil {
		log.Println(err)
		logOutput.Append(NewTextWithPrefix("预约失败: " + err.Error()))
	} else {
		logOutput.Append(NewTextWithPrefix("预约成功，场次信息：" + openRooms[randIndex][1] + "，请及时添加同行人并尽快支付。"))
	}

}

// NewTextWithPrefix 生成带有时间前缀的字符串，用于日志的输出
func NewTextWithPrefix(s string) string {
	return "\n" + time.Now().Format("15:04:05") + "> " + s
}

// isValidDateFormat 验证日期格式
func isValidDateFormat(date string) bool {
	// 该例子中假设日期格式为YYYY-MM-DD
	dateRegex := `^\d{4}-\d{2}-\d{2}$`
	// 编译正则表达式
	regex, err := regexp.Compile(dateRegex)
	if err != nil {
		fmt.Println("正则表达式编译错误:", err)
		return false
	}
	// 使用正则表达式匹配字符串
	match := regex.MatchString(date)

	return match
}

// BasicInfo 描述了下单所需的基本信息
type BasicInfo struct {
	name      string // 姓名
	id        string // 学生卡id
	phone     string // 手机号
	cookie    string //
	date      string // 格式 YYYY-MM-DD
	beginHour string // 预约时间段开始
	endHour   string // 预约时间段结束
	sport     string // 运动场馆类型
}

func main() {
	//ctx, cancel := context.WithCancel(context.Background())

	var basicInfo BasicInfo
	var orderCount int // 已预约的次数

	// 创建App与窗口
	myApp := app.New()
	myWindow := myApp.NewWindow("体育场馆预约")
	myWindow.Resize(fyne.NewSize(500, 450))
	// 设置了中文字体
	myApp.Settings().SetTheme(&theme.MyTheme{})

	// logContent
	logOutput := &widget.Entry{
		MultiLine: true,                         // 多行
		Wrapping:  fyne.TextWrapOff,             //
		Scroll:    container.ScrollVerticalOnly, // 只允许垂直华东
		TextStyle: fyne.TextStyle{ // 文本的格式，由于设置了中文字体，似乎不work
			Bold:   true,
			Italic: true,
		},
	}
	logOutput.SetPlaceHolder(`> 注意事项：
1.日期格式YYYY-MM-DD，例如2023-01-01
2.请提前1~5分钟验证cookie的有效性。
3.保证cookie开头和结尾没有回车或空格。`)
	//logOutput.Disable() // 不希望用户能够修改日志，但暗黑时效果并不好
	logTitle := canvas.NewText("日志", color.White)
	logContent := container.NewBorder(logTitle, nil, nil, nil, container.NewPadded(logOutput))

	// sportSelectContent
	sportSelect := widget.NewSelect([]string{"1-羽毛球-南区东馆羽毛球场"}, func(sportName string) {
		basicInfo.sport = sportName[0:1]
	})
	sportSelect.PlaceHolder = "请选择运动类型与场馆"
	sportSelectTitle := canvas.NewText("运动-场馆", color.White)
	sportSelectContent := container.NewGridWithRows(2, sportSelectTitle, sportSelect)

	// cookieContent
	cookieInput := &widget.Entry{MultiLine: true, Wrapping: fyne.TextWrapWord, Scroll: container.ScrollHorizontalOnly}
	cookieInput.SetPlaceHolder("请输入cookie")
	cookieTitle := canvas.NewText("cookie", color.White)
	valCookieButton := widget.NewButton("验证cookie", func() {
		var dialogWindow dialog.Dialog
		basicInfo.cookie = cookieInput.Text
		err, isPass := ValidateCookie(basicInfo.cookie)
		if err != nil {
			log.Println(err.Error())
			logOutput.Append("\n< 验证cookie出错" + err.Error())
		}
		if isPass {
			dialogWindow = dialog.NewInformation("成功", "cookie验证成功", myWindow)
		} else {
			dialogWindow = dialog.NewInformation("失败", "cookie失效，请重新获取并验证", myWindow)
		}
		dialogWindow.Show()
	})
	cookieContent := container.NewBorder(nil, nil, cookieTitle, valCookieButton, cookieInput)

	// topContent
	idInput, idInputContent := generateInputContainer("学号", "请输入10位学号")
	nameInput, nameInputContent := generateInputContainer("姓名", "请输入姓名")
	phoneInput, phoneInputContent := generateInputContainer("手机号码", "真实手机号")
	dateInput, dateInputContent := generateInputContainer("预约日期", "YYYY-MM-DD")
	inputLine1 := container.NewGridWithRows(1, idInputContent, nameInputContent, phoneInputContent)
	inputLine2 := container.NewGridWithRows(1, sportSelectContent, dateInputContent)
	topContent := container.NewVBox(inputLine1, inputLine2, cookieContent)

	// leftContent
	timeTitle := canvas.NewText("预约时间段", color.White)
	timeChoose := widget.NewRadioGroup([]string{
		"08:00-09:00",
		"09:00-10:00",
		"10:00-11:00",
		"11:00-12:00",
		"12:00-13:00",
		"13:00-14:00",
		"14:00-15:00",
		"15:00-16:00",
		"16:00-17:00",
		"17:00-18:00",
		"18:00-19:00",
		"19:00-20:00",
		"20:00-21:00",
		"21:00-22:00",
	}, func(hourRange string) {
		if hourRange != "" {
			basicInfo.beginHour = hourRange[0:2]
			basicInfo.endHour = hourRange[6:8]
		} else {
			basicInfo.beginHour = ""
			basicInfo.endHour = ""
		}
	})
	leftContent := container.NewVBox(timeTitle, timeChoose)

	// bottomContent
	var startButton *widget.Button
	startButton = widget.NewButton("开始", func() {
		// 关闭预约按钮
		startButton.Disable()
		defer startButton.Enable()
		// 获取组件所输入的基础信息
		orderCount++
		logOutput.Append(fmt.Sprintf("\n\n#############第%d次预约#############", orderCount))
		basicInfo.id = idInput.Text
		basicInfo.name = nameInput.Text
		basicInfo.phone = phoneInput.Text
		basicInfo.cookie = cookieInput.Text
		basicInfo.date = dateInput.Text
		// 检查所需信息的完整性
		if checkZeroField(basicInfo) {
			dialogWindow := dialog.NewInformation("预约失败", "请输入完整信息再确认", myWindow)
			logOutput.Append(NewTextWithPrefix("预约失败: 信息不完整"))
			dialogWindow.Show()

		} else {
			if !isValidDateFormat(basicInfo.date) {
				dialogWindow := dialog.NewInformation("预约失败", "请检查日期格式", myWindow)
				logOutput.Append(NewTextWithPrefix("预约失败: 日期格式有误"))
				dialogWindow.Show()
				return
			}
			logOutput.Append(NewTextWithPrefix("基础信息:"))
			logOutput.Append("\n	学    号: " + basicInfo.id)
			logOutput.Append("\n	姓    名: " + basicInfo.name)
			logOutput.Append("\n	手    机: " + basicInfo.phone)
			logOutput.Append("\n	场馆编号: " + basicInfo.sport)
			logOutput.Append("\n	预约日期: " + basicInfo.date)
			logOutput.Append("\n	预约时段: " + basicInfo.beginHour + ":00-" + basicInfo.endHour + ":00")
			logOutput.Append(NewTextWithPrefix("获取空闲场地"))
			AddOrderWorker(basicInfo, logOutput)
		}
	})
	////停止预约按钮
	//stopButton := widget.NewButton("停止", func() {
	//	cancel()
	//	logOutput.Append(NewTextWithPrefix("您停止了抢票"))
	//})
	// 重置日志按钮
	resetButton := widget.NewButton("重置日志", func() {
		dialog.NewConfirm("注意", "重置将会清空日志信息", func(b bool) {
			if b {
				logOutput.SetText("")
			}
		}, myWindow).Show()
	})

	// 快捷链接入口1
	sportUrlString := "https://ehall.szu.edu.cn/qljfwapp/sys/lwSzuCgyy/index.do#/sportVenue"
	sportUrl, _ := url.Parse(sportUrlString)
	sportLink := widget.NewHyperlink("体育场馆预约快捷入口", sportUrl)

	// 快捷链接入口2
	payUrlString := "https://ehall.szu.edu.cn/qljfwapp/sys/lwSzuCgyy/index.do#/myBooking"
	payUrl, _ := url.Parse(payUrlString)
	payLink := widget.NewHyperlink("支付/添加同行人快捷入口", payUrl)

	// 使用了NewPacer将链接和按钮分别往两头放置
	bottomContent := container.NewHBox(sportLink, payLink, layout.NewSpacer(), resetButton, startButton)

	content := container.NewBorder(topContent, bottomContent, leftContent, nil, logContent)
	// 显示窗口
	myWindow.SetContent(content)
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}
