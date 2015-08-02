# 作用
给定一个url，获取他包含的css，js，img以及css中包含的background图片

# Demo
以下命令可以获取CSS禅意花园的所有素材
```shell
./GetWeb -u http://www.csszengarden.com/ -d csszen
```

注意，可能需要手动修改css、js等文件的名字。因为下载完成的文件是不带版本号的。例如css以这种形式引入: 
```css
<link rel="stylesheet" media="screen" href="/214/214.css?v=8may2013">
```
那么下载下来的css就会存储在`csszen/214/214.css`中，在index.html中需要把以上代码改为

```css
<link rel="stylesheet" media="screen" href="214/214.css">
```

