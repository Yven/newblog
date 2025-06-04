CREATE TABLE IF NOT EXISTS article(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	slug TEXT NOT NULL,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	cid INTEGER DEFAULT 0,
	create_time INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
	update_time INTEGER NOT NULL DEFAULT (strftime('%s', 'now')),
	delete_time INTEGER
);

CREATE TABLE IF NOT EXISTS category(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS tag(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS article_tag(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	aid INTEGER NOT NULL,
	tid INTEGER NOT NULL
);

-- 插入测试数据，仅在表为空时插入
INSERT INTO category (name)
SELECT '测试'
WHERE NOT EXISTS (SELECT 1 FROM category);

INSERT INTO article (slug, title, content, cid)
SELECT 'test', 'Markdown 入门指南', '# Markdown 入门指南

## 什么是 Markdown？

Markdown 是一种轻量级标记语言，创始人为约翰·格鲁伯（John Gruber）。它允许人们使用易读易写的纯文本格式编写文档，然后转换成有效的 HTML 文档。

## 为什么选择 Markdown？

1. 简单易学
2. 纯文本编写
3. 可转换为多种格式
4. 广泛支持
5. 专注于内容

## 基本语法示例

### 1. 文本样式

普通文本就直接写就好了。

**这是加粗的文字**
*这是斜体文字*
***这是加粗斜体文字***
~~这是删除线文字~~

### 2. 列表

无序列表：
- 苹果
- 香蕉
- 橙子

有序列表：
1. 第一步
2. 第二步
3. 第三步

### 3. 引用

> 这是一段引用
> 可以有多行
>> 这是嵌套引用

### 4. 代码

行内代码：`print("Hello World")`

代码块：
```go
def hello():
print("Hello World")
```

### 5. 链接和图片

[百度](https://www.baidu.com)
![百度](https://www.baidu.com/favicon.ico)

## 高级语法

### 1. 表格

| 姓名 | 年龄 | 性别 |
| ---- | ---- | ---- |
| 张三 | 20   | 男   |
| 李四 | 21   | 女   |

### 2. 数学公式

$E=mc^2$

## 结语
Markdown 是一种非常流行的标记语言，它的语法简单，易于学习和使用。如果你想快速编写文档，Markdown 是一个很好的选择。', 1
WHERE NOT EXISTS (SELECT 1 FROM article);
