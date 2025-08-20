package main

import (
	"log"

	"github.com/signintech/gopdf"
)

// 基本的参数
var (
	pageWidth   = gopdf.PageSizeA4.W // A4纸的宽度
	pageHeight  = gopdf.PageSizeA4.H // A4纸的高度
	marginLeft  = 13.0
	marginRight = 13.0
	marginTop   = 20.0
	marginBottom = 20.0
)

func main() {
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	pdf.AddPage()

	loadFonts(&pdf)

	// 第一页内容
	addResumeHeader(&pdf, "个人简历")
	addBasicInfoWithPhoto(&pdf, "photo.jpg")
	addWorkExperience(&pdf)
	addSXProjectExperience(&pdf)
	
	// 检查是否需要分页
	checkAndAddNewPage(&pdf)
	
	// 第二页内容
	addSTProjectExperience(&pdf)
	addSkills(&pdf)

	err := pdf.WritePdf("final_resume.pdf")
	if err != nil {
		log.Fatalf("保存 PDF 失败: %v", err)
	}
	log.Println("PDF 文件已生成: final_resume.pdf")
}

// checkAndAddNewPage 检查当前页面剩余空间，如果不足则添加新页面
func checkAndAddNewPage(pdf *gopdf.GoPdf) {
	currentY := pdf.GetY()
	
	// 在东方算芯项目经历后强制分页，确保商汤项目经历在第二页开始
	// 这样可以保证简历的逻辑分布更清晰
	pdf.AddPage()
	pdf.SetY(marginTop)
	log.Printf("分页：第1页结束于Y=%.2f，第2页开始", currentY)
}

func loadFonts(pdf *gopdf.GoPdf) {
	// 尝试使用系统字体，如果失败则创建一个简单的字体映射
	// 在 Windows 系统上尝试使用常见的中文字体
	fontPaths := []string{
		"C:/Windows/Fonts/simhei.ttf",     // 黑体
		"C:/Windows/Fonts/simsun.ttc",     // 宋体
		"C:/Windows/Fonts/msyh.ttc",       // 微软雅黑
		"C:/Windows/Fonts/arial.ttf",      // Arial (英文)
	}
	
	fontLoaded := false
	for _, fontPath := range fontPaths {
		err := pdf.AddTTFFont("default", fontPath)
		if err == nil {
			log.Printf("成功加载字体: %s", fontPath)
			fontLoaded = true
			break
		}
	}
	
	if !fontLoaded {
		log.Println("警告: 无法加载任何字体，PDF 可能无法正常显示中文")
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
	err := pdf.SetFont("default", "", 18)
	if err != nil {
		log.Printf("设置标题字体失败: %v", err)
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
	addSectionContent(pdf, name, 18, 18, contentX)
	pdf.SetTextColor(0, 112, 192)
	content := []string{
		"2012-2016: 南通大学(本,机械工程)       2016-2019: 上海工程技术大学(硕,机械电子工程)",
	}
	addSectionContent(pdf, content, 14, 10, contentX)

	conn := []string{
        "电话: 13764433460   邮箱: nippon1218@163.com    出生年月:1993-12   已婚",
	}
	addSectionContent(pdf, conn, 16, 10, contentX)

	purpose := []string{
		"求职意向: 系统集成交付工程师(主), 测试开发, AI算子开发， golang开发，C/C++开发",
	}
	addSectionContent(pdf, purpose, 18, 13, contentX)
	pdf.Br(30)
}

// addSectionTitle 添加部分标题，并绘制双分割线
func addSectionTitle(pdf *gopdf.GoPdf, title string) {
	err := pdf.SetFont("default", "", 14)
	if err != nil {
		log.Printf("设置标题字体失败: %v", err)
	}
	pdf.SetTextColor(0, 0, 0) // 黑色文字
	pdf.SetX(marginLeft)
	pdf.Cell(nil, title)
	pdf.Br(12)
	drawDoubleLine(pdf, marginLeft, pdf.GetY(), pageWidth-marginRight)
	pdf.Br(12)
}

// addSectionContent 添加部分内容
func addSectionContent(pdf *gopdf.GoPdf, content []string, lineSpacing float64, fontsize int, X float64) {
	err := pdf.SetFont("default", "", fontsize)
	if err != nil {
		log.Printf("设置内容字体失败: %v", err)
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
			Period:  "2025.01 - 至今",
			Company: "上海东方算芯有限公司",
			Role:    "AI系统集成开发工程师",
			Duties: []string{
				"负责国产GPU上AI软件栈的集成与发布, 参与系统debug，优化发布流程",
				"负责大语言模型(llama, gpt)算子实现，使用triton开发和维护常用算子",
				"使用pytorch及megatron分布式训练框架维护日常模型的整网训练",
				"负责国产GPU上llama和GPT模型的训练流程开发",
				"开发精度对比工具sdc-insight, 能dump出模型的权重，输入，输出等中间层信息",
				"参与初创公司流程规范建设, 制定版本发布与回退策略",
				"从0开始搭建AI软件栈的自动化测试, 包括CI/CD流程的优化，集成测试",
				"带领多名外包与实习生完成集成测试任务，包括测试任务分解，任务分配，任务跟踪，任务验收等",
			},
		},
		{
			Period:  "2021.05 - 2024.12",
			Company: "临港绝影智能科技有限公司(商汤科技)",
			Role:    "系统集成开发工程师",
			Duties: []string{
				"负责绝影全栈项目pilot/pap项目的集成与发版，保证内部迭代的稳定, 负责新车的使能, 中间件配置等",
				"保证版本可追溯性，快速识别并回退dirty commit",
				"负责POC任务及商务演示的对接，提供需要的稳定版本与支持",
				"优化CI/CD，优化发版流程, 引入自动化测试流程等",
				"维护与开发monitor, trace平台",
			},
		},
		{
			Company: "霍尼韦尔(中国)有限公司",
			Period:  "2019.05 - 2021.04",
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
		err := pdf.SetFont("default", "", 12)
		if err != nil {
			log.Printf("设置字体失败: %v", err)
		}
		pdf.Cell(nil, exp.Period)

		// 获取当前 Y 位置，确保同一行
		currentY := pdf.GetY()

		// 设置 company 的居中位置和字体
		err = pdf.SetFont("default", "", 12)
		if err != nil {
			log.Printf("设置字体失败: %v", err)
		}
		companyWidth, _ := pdf.MeasureTextWidth(exp.Company)
		companyX := (pageWidth - companyWidth) / 2
		pdf.SetX(companyX)
		pdf.SetY(currentY)
		pdf.Cell(nil, exp.Company)

		// 设置 role 的位置和字体
		err = pdf.SetFont("default", "", 12)
		if err != nil {
			log.Printf("设置字体失败: %v", err)
		}
		roleWidth, _ := pdf.MeasureTextWidth(exp.Role)
		pdf.SetX(pageWidth - 15 - roleWidth)
		pdf.SetY(currentY)
		pdf.Cell(nil, exp.Role)

		// 换行后开始 Duties 部分
		pdf.Br(12)
		dutyX := 50.0
		for _, duty := range exp.Duties {
			// 设置圆点符号
			pdf.SetX(dutyX)
			pdf.Cell(nil, "•") // 使用实心圆点

			// 设置 Duty 的文本
			pdf.SetX(dutyX + 12)
			err := pdf.SetFont("default", "", 11)
			if err != nil {
				log.Printf("设置字体失败: %v", err)
			}
			pdf.Cell(nil, duty)
			pdf.Br(12)
		}

		// 为不同经历之间增加空行
		pdf.Br(12)
	}
}

// addSXProjectExperience 添加算芯项目经历
func addSXProjectExperience(pdf *gopdf.GoPdf) {
	addSectionTitle(pdf, "东方算芯核心项目经历")

	pageWidth := gopdf.PageSizeA4.W // 获取页面宽度
	projects := []struct {
		Title        string
		Description  string
		Period       string
		TechKeywords string
		Duties       []string
	}{
		{
			Title:        "AI软件栈集成与交付",
			Description:  "在国产GPU SDC100上完成AI软件栈的集成与交付",
			Period:       "2025.01 - 至今",
			TechKeywords: "大模型训练，国产GPU，AI软件栈集成",
			Duties: []string{
			    "负责公司整个软件产品的集成与交付，包括编译构建，打包部署，CI/CD流程的优化",
			    "负责软件集成版本的冒烟和集成测试",
			    "负责日常llama2, gpt模型的训练与精度看护",
				"负责从头搭建AI软件栈的自动化集成测试平台",
				"负责CICD集成看护策略的制定与代码实现",
				"带领外包和实习生开发测试用例(精度，性能，稳定性等)与测试执行",
				"负责上下游团队的整合，对软件栈问题初筛，出错模块定位，排查性能，精度， 稳定性等问题",
				"负责对接客户，软件release版本的回溯，问题排查，问题修复等",
			},
		},
		{
			Title:        "大语言模型算子开发",
			Description:  "使用open AI的triton编译器，开发常用算子(gemm, softmax等)，benchmarks性能优于torch原生实现",
			Period:       "2025.01 - 至今",
			TechKeywords: "triton, python, CUDA",
			Duties: []string{
			    "使用triton负责算子开发，包括gemm, softmax等",
			    "开发算子benchmark测试工具",
			    "测试算子的正确性，性能，稳定性等",
				"移植flaggems等开源triton算子库",
			},
		},
		{
            Title:        "精度对比工具开发与维护",
			Description:  "基于华为开源工具msqtt开发SDC100平台的精度对比/dump工具sdc-insight",
			Period:       "2025.01 - 至今",
			TechKeywords: "pytorch, python",
			Duties: []string{
			    "负责需求收集，技术调研和技术选型, 制定排期",
				"移植msqtt到SDC100平台",
                "使用python开发模型训练算子输入输出数据dump功能",
                "使用python开发模型精度对比功能",
                "模型训练可视化",
			},
		},
	}

	for _, project := range projects {
		// 设置项目标题
		err := pdf.SetFont("default", "", 12)
		if err != nil {
			log.Printf("设置字体失败: %v", err)
		}
		pdf.SetTextColor(0, 0, 0)
		pdf.SetX(marginLeft + 5)
		pdf.Cell(nil, project.Title)

		// 设置时间段靠右
		periodWidth, _ := pdf.MeasureTextWidth(project.Period)
		pdf.SetX(pageWidth - 15 - periodWidth)
		pdf.Cell(nil, project.Period)

		// 换行
		pdf.Br(18)

		// 设置项目描述
		err = pdf.SetFont("default", "", 10)
		if err != nil {
			log.Printf("设置字体失败: %v", err)
		}
		pdf.SetX(marginLeft + 20)
		pdf.Cell(nil, "【描述】: "+project.Description)
		pdf.Br(13)

		// 设置技术关键词
		pdf.SetX(marginLeft + 20)
		pdf.SetTextColor(0, 0, 0) // 使用灰色文本突出关键词
		pdf.Cell(nil, "【关键词】: "+project.TechKeywords)
		pdf.Br(13)

		// 添加职责列表
		for _, duty := range project.Duties {
			pdf.SetX(marginLeft + 30)
			pdf.SetTextColor(0, 0, 0) // 黑色文字
			pdf.Cell(nil, "•")        // 实心圆点
			pdf.SetX(marginLeft +  50)
		    pdf.SetTextColor(100, 100, 100) // 使用灰色文本突出关键词
			pdf.Cell(nil, duty)
			pdf.Br(12)
		}

		// 项目之间增加空白间距
		pdf.Br(12)
	}
}

// addSTProjectExperience 添加商汤项目经历
func addSTProjectExperience(pdf *gopdf.GoPdf) {
	addSectionTitle(pdf, "商汤绝影核心项目经历")

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
			Period:       "2024.10 - 2024.11",
			TechKeywords: "性能平台开发，Go, C/C++, Grafana, ebpf",
			Duties: []string{
			    "负责需求收集，技术调研和技术选型, 制定排期",
			    "可观测性平台模块开发，使用ebpf技术和grafana",
				"集成opentemeletry以及prometheus等开源框架",
				"负责Monitor部分的设计与实现Demo，使用go为主",
			},
		},
		{
            Title:        "绝影全栈行泊一体产品Pilot集成交付owner",
			Description:  "基于MDC610/MDC210平台，包括自研中间件和工具链，HNOP, MNOP, CNOP, APA, HPP等功能",
			Period:       "2023.10 - 2024.12",
			TechKeywords: "全栈集成, git/cmake，流程优化",
			Duties: []string{
				"集成AutoSar EM，CM等中间件模块，配置EM与CM, 维护自动驾驶软件栈的通信链路",
				"集成全栈自动驾驶软件功能, 包括行车泊车各功能算法，HMI，车机仪表等，管理软件版本",
                "主导版本问题排查，对问题进行初筛，定位到模块级别，包括功能，性能, 稳定性等",
				"支持商务，对接商务/产品完成对外POC以及其他商务Demo, 如本田，比亚迪，奇瑞，东风",
				"优化迭代流程，引入并执行一些集成/交付新方法",
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
	}

	for _, project := range projects {
		// 设置项目标题
		err := pdf.SetFont("default", "", 12)
		if err != nil {
			log.Printf("设置字体失败: %v", err)
		}
		pdf.SetTextColor(0, 0, 0)
		pdf.SetX(marginLeft + 5)
		pdf.Cell(nil, project.Title)

		// 设置时间段靠右
		periodWidth, _ := pdf.MeasureTextWidth(project.Period)
		pdf.SetX(pageWidth - 15 - periodWidth)
		pdf.Cell(nil, project.Period)

		// 换行
		pdf.Br(18)

		// 设置项目描述
		err = pdf.SetFont("default", "", 10)
		if err != nil {
			log.Printf("设置字体失败: %v", err)
		}
		pdf.SetX(marginLeft + 20)
		pdf.Cell(nil, "【描述】: "+project.Description)
		pdf.Br(13)

		// 设置技术关键词
		pdf.SetX(marginLeft + 20)
		pdf.SetTextColor(0, 0, 0) // 使用灰色文本突出关键词
		pdf.Cell(nil, "【关键词】: "+project.TechKeywords)
		pdf.Br(13)

		// 添加职责列表
		for _, duty := range project.Duties {
			pdf.SetX(marginLeft + 30)
			pdf.SetTextColor(0, 0, 0) // 黑色文字
			pdf.Cell(nil, "•")        // 实心圆点
			pdf.SetX(marginLeft +  50)
		    pdf.SetTextColor(100, 100, 100) // 使用灰色文本突出关键词
			pdf.Cell(nil, duty)
			pdf.Br(11)
		}

		// 项目之间增加空白间距
		pdf.Br(12)
	}
}

// addSkills 添加个人技能
func addSkills(pdf *gopdf.GoPdf) {
	addSectionTitle(pdf, "个人技能与优势")
	skills := []string{
		"1. 熟悉大语言模型AI软件栈的集成, 熟悉大语言模型的分布式训练过程",
		"2. 熟悉pytorch，megatron等分布式训练框架，熟悉triton等AI编译器实现AI算子",
		"3. 熟悉pytest框架, 具备从0-1搭建自动化测试框架能力，熟悉CI/CD流程",
		"4. 熟悉自动驾驶软件全栈的集成, 了解产品上线到交付发版的过程，风险识别等",
		"5. 熟悉linux系统常见的性能问题原因, 擅长定位内存泄漏及cpu问题, 熟悉华为MDC平台域控制器",
		"6. 熟练掌握python, golang, 熟悉C/C++， CMake/Makefile, Git, GDB等常用开发工具或者平台",
		"7. 英语CET-6通过，可进行技术文档阅读与撰写",
	}
	addSectionContent(pdf, skills, 12, 10, 20)
}
