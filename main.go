package main

import (
	"log"

	"github.com/signintech/gopdf"
)

// 基本的参数
var (
	pageWidth   = gopdf.PageSizeA4.W // A4纸的宽度
	marginLeft  = 13.0
	marginRight = 13.0
)

func main() {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	loadFonts(&pdf)

	addResumeHeader(&pdf, "个人简历")
	addBasicInfoWithPhoto(&pdf, "photo.jpg")
	addWorkExperience(&pdf)
	addProjectExperience(&pdf)
	addSkills(&pdf)

	err := pdf.WritePdf("final_resume.pdf")
	if err != nil {
		log.Fatalf("保存 PDF 失败: %v", err)
	}
	log.Println("PDF 文件已生成: final_resume.pdf")
}

func loadFonts(pdf *gopdf.GoPdf) {
	err := pdf.AddTTFFont("NotoSansSC-Regular", "./fonts/static/NotoSansSC-Regular.ttf")
	if err != nil {
		log.Fatalf("加载常规字体失败: %v", err)
	}
	err = pdf.AddTTFFont("NotoSansSC-Bold", "./fonts/static/NotoSansSC-Bold.ttf")
	if err != nil {
		log.Fatalf("加载加粗字体失败: %v", err)
	}
}

// addResumeHeader 添加简历标题，带背景矩形框
func addResumeHeader(pdf *gopdf.GoPdf, title string) {
	blueRectWidth := 20.0 // 蓝色背景的宽度（可以根据需求设置）
	// 设置蓝色矩形（从页面最左边开始，宽度为 blueRectWidth）
	pdf.SetFillColor(0, 112, 192)                                 // 蓝色
	pdf.RectFromUpperLeftWithStyle(0, 10, blueRectWidth, 20, "F") // 从页面最左边开始

	// 设置黑色矩形（从蓝色矩形的右边开始，宽度为剩余部分）
	pdf.SetFillColor(0, 0, 0)                                                           // 黑色
	pdf.RectFromUpperLeftWithStyle(blueRectWidth, 10, pageWidth-blueRectWidth, 20, "F") // 从蓝色矩形右边开始

	// 添加白色标题文字
	err := pdf.SetFont("NotoSansSC-Bold", "", 18)
	if err != nil {
		log.Fatalf("设置标题字体失败: %v", err)
	}
	pdf.SetTextColor(255, 255, 255) // 白色文字
	pdf.SetX(blueRectWidth + 5)     // 设置文字偏移位置
	pdf.SetY(12)                    // 保证文字居中矩形框
	pdf.Cell(nil, title)
	pdf.Br(30) // 适当增加行距，确保标题和内容之间有足够的空隙
}

// addBasicInfoWithPhoto 添加基本信息，并排放照片
func addBasicInfoWithPhoto(pdf *gopdf.GoPdf, imgPath string) {
	// 添加照片
	photoWidth, photoHeight := 66.0, 88.0
	photoX, photoY := 10.0, pdf.GetY()
	err := pdf.Image(imgPath, photoX, photoY, &gopdf.Rect{W: photoWidth, H: photoHeight})
	if err != nil {
		log.Fatalf("添加照片失败: %v", err)
	}

	// 添加基本信息
	contentX := photoX + photoWidth + 10
	pdf.SetX(contentX)
	pdf.SetY(photoY)
	pdf.SetTextColor(0, 0, 0) // 黑色文字
	name := []string{
		"陈文毅",
	}
	addSectionContent(pdf, name, 22, 18, contentX)
	pdf.SetTextColor(0, 112, 192)
	content := []string{
		"2012-2016: 南通大学(本,机械工程)       2016-2019: 上海工程技术大学(硕,机械电子工程)",
	}
	addSectionContent(pdf, content, 16, 10, contentX)

	conn := []string{
        "电话: 13764433460   邮箱: nippon1218@163.com    出生年月:1993-12   已婚",
	}
	addSectionContent(pdf, conn, 20, 10, contentX)

	purpose := []string{
		"求职意向: 集成交付工程师(主)  golang开发，C/C++开发, python自动化测试",
	}
	addSectionContent(pdf, purpose, 25, 13, contentX)
	pdf.Br(10)
}

// addSectionTitle 添加部分标题，并绘制双分割线
func addSectionTitle(pdf *gopdf.GoPdf, title string) {
	err := pdf.SetFont("NotoSansSC-Bold", "", 14)
	if err != nil {
		log.Fatalf("设置标题字体失败: %v", err)
	}
	pdf.SetTextColor(0, 0, 0) // 黑色文字
	pdf.SetX(marginLeft)
	pdf.Cell(nil, title)
	pdf.Br(10)
	drawDoubleLine(pdf, marginLeft, pdf.GetY(), pageWidth-marginRight)
	pdf.Br(10)
}

// addSectionContent 添加部分内容
func addSectionContent(pdf *gopdf.GoPdf, content []string, lineSpacing float64, fontsize int, X float64) {
	err := pdf.SetFont("NotoSansSC-Regular", "", fontsize)
	if err != nil {
		log.Fatalf("设置内容字体失败: %v", err)
	}
	for _, line := range content {
		pdf.SetX(X)
		pdf.Cell(nil, line)
		pdf.Br(lineSpacing)
	}
}

// drawDoubleLine 绘制双分割线
func drawDoubleLine(pdf *gopdf.GoPdf, xStart, y, xEnd float64) {
	// 绘制底层细分割线
	yset := y + 5
	pdf.SetLineWidth(0.5)
	pdf.Line(xStart, yset, xEnd, yset)

	// 绘制上层短粗分割线
	pdf.SetLineWidth(1.5)
	shortLineLength := (xEnd - xStart) / 4
	midXEnd := xStart + shortLineLength
	pdf.Line(xStart, yset, midXEnd, yset)

	// 恢复默认线宽
	pdf.SetLineWidth(0.5)
}

// addWorkExperience 添加工作经历
func addWorkExperience(pdf *gopdf.GoPdf) {
	addSectionTitle(pdf, "工作经历")

	pageWidth := gopdf.PageSizeA4.W // 获取页面宽度
	experiences := []struct {
		Period  string
		Company string
		Role    string
		Duties  []string
	}{
		{
			Period:  "2021.04 - 至今",
			Company: "临港绝影智能科技有限公司(商汤科技)",
			Role:    "系统集成开发工程师",
			Duties: []string{
				"负责交付自动驾驶软件给客户，版本集成及对应的集成测试, 保证交付质量",
				"参与Swc模块的软件开发与性能优化, 与产品，系统和开发人员共同制定接口标准",
				"负责商汤全栈项目的集成与发版，保证内部迭代的稳定, 负责新车的使能, 中间件配置等",
				"保证版本可追溯性，快速识别并回退dirty commit",
				"负责POC任务及商务演示的对接，提供需要的稳定版本与支持",
				"牵头识别解决集成版本各种“奇怪”问题，提供经验性的总结，包括bug定位及性能指标方面",
				"优化CI/CD，优化发版流程, 引入自动化测试流程等",
				"负责团队内部梯度建设，新人mentor引导, lead一支3~5人小交付团队",
				"负责monitor, trace平台从0-1的开发与上线",
			},
		},
		{
			Company: "霍尼韦尔(中国)有限公司",
			Period:  "2019 - 2021.04",
			Role:    "C/C++ 软件开发工程师",
			Duties: []string{
				"Honeywell智能仪表项目功能开发，单元/集成测试与交付",
				"维护Honeywell DCS控制器源码",
				"为制造工厂提供产线标定脚本",
			},
		},
	}

	for _, exp := range experiences {
		// 设置 period 的位置和字体
		pdf.SetX(marginLeft + 5)
		err := pdf.SetFont("NotoSansSC-Regular", "", 12)
		if err != nil {
			log.Fatalf("设置字体失败: %v", err)
		}
		pdf.Cell(nil, exp.Period)

		// 获取当前 Y 位置，确保同一行
		currentY := pdf.GetY()

		// 设置 company 的居中位置和字体
		err = pdf.SetFont("NotoSansSC-Bold", "", 12)
		if err != nil {
			log.Fatalf("设置加粗字体失败: %v", err)
		}
		companyWidth, _ := pdf.MeasureTextWidth(exp.Company)
		companyX := (pageWidth - companyWidth) / 2
		pdf.SetX(companyX)
		pdf.SetY(currentY)
		pdf.Cell(nil, exp.Company)

		// 设置 role 的位置和字体
		err = pdf.SetFont("NotoSansSC-Regular", "", 12)
		if err != nil {
			log.Fatalf("设置字体失败: %v", err)
		}
		roleWidth, _ := pdf.MeasureTextWidth(exp.Role)
		pdf.SetX(pageWidth - 15 - roleWidth)
		pdf.SetY(currentY)
		pdf.Cell(nil, exp.Role)

		// 换行后开始 Duties 部分
		pdf.Br(15)
		dutyX := 50.0
		for _, duty := range exp.Duties {
			// 设置圆点符号
			pdf.SetX(dutyX)
			pdf.Cell(nil, "•") // 使用实心圆点

			// 设置 Duty 的文本
			pdf.SetX(dutyX + 12)
			err := pdf.SetFont("NotoSansSC-Regular", "", 11)
			if err != nil {
				log.Fatalf("设置字体失败: %v", err)
			}
			pdf.Cell(nil, duty)
			pdf.Br(13)
		}

		// 为不同经历之间增加空行
		pdf.Br(10)
	}
}

// addProjectExperience 添加项目经历
func addProjectExperience(pdf *gopdf.GoPdf) {
	addSectionTitle(pdf, "商汤核心项目经历")

	pageWidth := gopdf.PageSizeA4.W // 获取页面宽度
	projects := []struct {
		Title        string
		Description  string
		Period       string
		TechKeywords string
		Duties       []string
	}{
		{
			Title:        "可观测平台SenseFast开发",
			Description:  "负责内部可观测平台的开发与上线，解决自动驾驶性能分析和debug的难题",
			Period:       "2024.06 - 至今",
			TechKeywords: "性能平台开发，Go, C/C++, Grafana, ebpf",
			Duties: []string{
			    "负责需求收集，技术调研和技术选型, 制定排期",
			    "可观测性平台模块开发，使用ebpf技术和grafana",
				"集成opentemeletry以及prometheus等开源框架",
				"负责Monitor部分的设计与实现，使用go为主",
			},
		},
		{
            Title:        "绝影全栈行泊一体产品PilotPro集成交付owner",
			Description:  "MDC610平台为主，J6E平台为辅，包括自研中间件和工具链，HNOP, MNOP, CNOP, APA, HPP等功能",
			Period:       "2023.10 - 至今",
			TechKeywords: "全栈目集成, git/cmake，流程优化",
			Duties: []string{
				"集成EM，CM等中间件模块，HMI，安卓仪表等工具链，常规版本发布",
                "主导版本各种问题排查，包括性能问题，bug等,",
				"支持商务，对接商务/产品完成对外POC以及其他商务Demo, 如本田，比亚迪，奇瑞，东风",
				"优化迭代流程，引入并执行一些集成/交付新方法",
				"管理领导3人交付团队",
			},
		},
		{
			Title:        "Adas/L2/POC等多条项目线集成Owner及交付流程标准制定",
			Description:  "在公司不同的阶段集成交付各种产品软件, 引入并完善CI/CD以及其他集成流程",
			Period:       "2021.04 - 2024.04",
			TechKeywords: "CI/CD, 流程规范化， 对外需求对接，对内规范对齐等",
			Duties: []string{
				"商汤绝影在本田, 广汽， 奇瑞，驭动等项目的POC等项目中的持续集成角色",
				"对接外部需求，确定交付边界，集成规范，部署要求等",
				"与内部算法/工具链开发人员对接，包括设计规范，接口，以及测试内容",
                "主导CI/CD的构建流程，以及集成测试环节的上线",
			},
		},
		{
			Title:        "其他工作及成就",
			Description:  "团队建设出力",
			Period:       "2019 - 至今",
			TechKeywords: "新人Mentor, RoadMap定义等",
			Duties: []string{
				"新入职员工的基础培训",
				"参与制定团队从新人到专家的Roadmap",
				"定期组织团队间的知识分享，践行企业文化",
                "2022，2023年度商汤绝影优秀员工",
			},
		},
	}

	for _, project := range projects {
		// 设置项目标题
		err := pdf.SetFont("NotoSansSC-Bold", "", 12)
		if err != nil {
			log.Fatalf("设置加粗字体失败: %v", err)
		}
		pdf.SetTextColor(0, 0, 0)
		pdf.SetX(marginLeft + 5)
		pdf.Cell(nil, project.Title)

		// 设置时间段靠右
		periodWidth, _ := pdf.MeasureTextWidth(project.Period)
		pdf.SetX(pageWidth - 15 - periodWidth)
		pdf.Cell(nil, project.Period)

		// 换行
		pdf.Br(14)

		// 设置项目描述
		err = pdf.SetFont("NotoSansSC-Regular", "", 10)
		if err != nil {
			log.Fatalf("设置常规字体失败: %v", err)
		}
		pdf.SetX(marginLeft + 20)
		pdf.Cell(nil, "描述: "+project.Description)
		pdf.Br(12)

		// 设置技术关键词
		pdf.SetX(marginLeft + 20)
		pdf.SetTextColor(0, 0, 0) // 使用灰色文本突出关键词
		pdf.Cell(nil, "关键词: "+project.TechKeywords)
		pdf.Br(12)

		// 添加职责列表
		for _, duty := range project.Duties {
			pdf.SetX(marginLeft + 30)
			pdf.SetTextColor(0, 0, 0) // 黑色文字
			pdf.Cell(nil, "•")        // 实心圆点
			pdf.SetX(marginLeft +  50)
		    pdf.SetTextColor(100, 100, 100) // 使用灰色文本突出关键词
			pdf.Cell(nil, duty)
			pdf.Br(10)
		}

		// 项目之间增加空白间距
		pdf.Br(8)
	}
}

// addSkills 添加个人技能
func addSkills(pdf *gopdf.GoPdf) {
	addSectionTitle(pdf, "个人技能与优势")
	skills := []string{
		"1. 熟悉自动驾驶软件全栈的集成, 了解产品上线到交付发版的过程，风险识别等",
		"2. 熟悉linux系统常见的性能问题原因, 擅长定位内存泄漏及cpu问题, 熟悉华为MDC平台域控制器",
		"3. 有系统可观测平台（tracing, metric, log，debug）, ebpf, grafana等开发经验, 可搭建通用监控平台",
		"4. 掌握golang, 熟悉C/C++，python, jenkins CI/CD, Git, GDB等在自动驾驶领域运用广的开发工具或者平台",
		"5. 具备车端debug和工具开发能力，特别是golang适合编写平台无关的通用工具集",
		"6. 英语CET-6通过，可进行技术文档阅读与撰写",
	}
	addSectionContent(pdf, skills, 12, 10, 20)
}
