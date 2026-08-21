# Bug 复现

## Bug 是什么

删除用户后其关注关系仍可被查询。

## 如何触发

创建两个用户建立关注，删除关注者，再按关注者筛选关注列表。

## 错误信息

`expected deleted user's follows removed, got total=1 len=1`
