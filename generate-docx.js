const { Document, Packer, Paragraph, TextRun, HeadingLevel } = require('docx');
const fs = require('fs');
const path = require('path');

// Markdown源文件目录
const sourceDir = 'c:\\Users\\yanglei\\AppData\\Roaming\\WorkBuddy\\User\\globalStorage\\tencent-cloud.coding-copilot\\brain\\d3f08bf6d9a949e3a7d391a047b6bee2';

// 读取所有Markdown文件
const files = [
    '01-框架概述.md',
    '02-核心概念.md',
    '03-HTTP服务开发.md',
    '04-RPC服务开发.md',
    '05-定时任务开发.md',
    '06-消息队列消费.md',
    '07-服务注册与发现.md',
    '08-组件使用.md',
    '09-15-实战案例.md',
];

const docChildren = [];

// 添加标题
docChildren.push(
    new Paragraph({
        heading: HeadingLevel.TITLE,
        children: [new TextRun({ text: 'Hydra微服务框架使用手册', bold: true, size: 48 })],
    }),
    new Paragraph({ text: '' })
);

// 解析Markdown并转换为Word文档元素
function parseMarkdownToDocx(content) {
    const lines = content.split('\n');
    let elements = [];
    let inCodeBlock = false;
    let codeContent = [];
    let codeLanguage = '';
    
    for (let i = 0; i < lines.length; i++) {
        const line = lines[i];
        
        // 代码块开始/结束
        if (line.startsWith('```')) {
            if (!inCodeBlock) {
                inCodeBlock = true;
                codeLanguage = line.slice(3).trim();
                codeContent = [];
            } else {
                inCodeBlock = false;
                if (codeContent.length > 0) {
                    // 添加代码块内容，每行作为单独的段落，使用等宽字体
                    codeContent.forEach((c, idx) => {
                        elements.push(new Paragraph({
                            children: [new TextRun({ 
                                text: c, 
                                font: 'Consolas',
                                size: 16
                            })],
                        }));
                    });
                }
            }
            continue;
        }
        
        if (inCodeBlock) {
            codeContent.push(line);
            continue;
        }
        
        // 处理标题
        if (line.startsWith('# ')) {
            elements.push(new Paragraph({
                heading: HeadingLevel.HEADING_1,
                children: [new TextRun({ text: line.slice(2), bold: true, size: 32 })],
            }));
        } else if (line.startsWith('## ')) {
            elements.push(new Paragraph({
                heading: HeadingLevel.HEADING_2,
                children: [new TextRun({ text: line.slice(3), bold: true, size: 28 })],
            }));
        } else if (line.startsWith('### ')) {
            elements.push(new Paragraph({
                heading: HeadingLevel.HEADING_3,
                children: [new TextRun({ text: line.slice(4), bold: true, size: 24 })],
            }));
        } else if (line.startsWith('#### ')) {
            elements.push(new Paragraph({
                heading: HeadingLevel.HEADING_4,
                children: [new TextRun({ text: line.slice(5), bold: true, size: 22 })],
            }));
        } else if (line.trim()) {
            // 普通段落，处理行内代码
            let text = line;
            let children = [];
            
            if (text.includes('`')) {
                const parts = text.split('`');
                parts.forEach((part, idx) => {
                    if (idx % 2 === 1) {
                        children.push(new TextRun({ 
                            text: part, 
                            font: 'Consolas',
                            size: 18
                        }));
                    } else if (part) {
                        children.push(new TextRun({ text: part, size: 20 }));
                    }
                });
            } else {
                children.push(new TextRun({ text: text, size: 20 }));
            }
            
            if (children.length > 0) {
                elements.push(new Paragraph({ children: children }));
            }
        } else {
            elements.push(new Paragraph({ text: '' }));
        }
    }
    
    return elements;
}

// 读取并添加各章节内容
files.forEach(fileName => {
    const filePath = path.join(sourceDir, fileName);
    try {
        const content = fs.readFileSync(filePath, 'utf8');
        console.log(`Processing: ${fileName}`);
        const elements = parseMarkdownToDocx(content);
        docChildren.push(...elements);
        docChildren.push(new Paragraph({ text: '' }));
        docChildren.push(new Paragraph({ text: '' }));
    } catch (e) {
        console.error(`Error reading ${fileName}:`, e.message);
    }
});

// 创建文档
const doc = new Document({
    sections: [{
        properties: {
            page: {
                size: { width: 12240, height: 15840 },
                margin: { top: 1440, right: 1440, bottom: 1440, left: 1440 }
            }
        },
        children: docChildren
    }]
});

// 保存文档到当前目录
const outputPath = 'D:\\yanglei\\work\\src\\github.com\\micro-plat\\hydra\\Hydra微服务框架使用手册.docx';
Packer.toBuffer(doc).then(buffer => {
    fs.writeFileSync(outputPath, buffer);
    console.log(`文档生成成功！保存到: ${outputPath}`);
}).catch(err => {
    console.error('生成文档失败:', err);
});
