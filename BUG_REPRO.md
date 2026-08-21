# Bug 复现

## Bug 是什么

对待开播房间执行下播失败后，房间状态仍被改成 ended。

## 如何触发

创建 pending 房间，直接执行下播，再查询房间状态。

## 错误信息

`expected pending status, got ended`
