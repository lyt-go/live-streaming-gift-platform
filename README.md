# 直播平台（live-streaming）

一个纯 Go 标准库（`net/http` + 标准库，零第三方依赖）实现的直播平台后端服务。支持用户、直播间、弹幕、礼物、关注关系与送礼记录的管理，含直播间开播/下播状态机、弹幕审核状态机、送礼跨实体校验与打赏榜统计。

## 技术栈

- 语言：Go 1.22+
- 依赖：仅标准库
- 存储：内存（进程内，线程安全）

## 运行

```bash
cd origin
go run ./cmd/server
# 或
go build -o server ./cmd/server && ./server
```

服务默认监听 `:8080`，可通过环境变量 `PORT` 或 `ADDR` 覆盖。

## 统一响应

```json
{"code":0,"message":"ok","data":{}}
```

错误映射：`ValidationError→400`、`ErrNotFound→404`、`ErrConflict→409`、其他→500。

## 金额单位

礼物价格与送礼金额均以「分」为单位（`int64`），避免浮点精度问题。

## API 一览

### 用户 User

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/users` | 创建用户 |
| GET | `/api/users` | 用户列表（分页，筛选 role/status/keyword） |
| GET | `/api/users/{id}` | 查询用户 |
| PUT | `/api/users/{id}` | 更新用户 |
| DELETE | `/api/users/{id}` | 删除用户 |
| POST | `/api/users/{id}/ban` | 封禁用户 |

### 直播间 Room

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/rooms` | 创建直播间 |
| GET | `/api/rooms` | 直播间列表（筛选 category/status/streamer_id/keyword） |
| GET | `/api/rooms/{id}` | 查询直播间 |
| PUT | `/api/rooms/{id}` | 更新直播间 |
| DELETE | `/api/rooms/{id}` | 删除直播间 |
| POST | `/api/rooms/{id}/start` | 开播（pending→live） |
| POST | `/api/rooms/{id}/end` | 下播（live→ended） |

### 弹幕 Danmaku

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/danmakus` | 发送弹幕 |
| GET | `/api/danmakus` | 弹幕列表（筛选 room_id/user_id/status/keyword） |
| GET | `/api/danmakus/{id}` | 查询弹幕 |
| POST | `/api/danmakus/{id}/approve` | 审核通过 |
| POST | `/api/danmakus/{id}/block` | 拦截 |
| POST | `/api/danmakus/batch-approve` | 批量审核通过 |

### 礼物 Gift

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/gifts` | 创建礼物 |
| GET | `/api/gifts` | 礼物列表（筛选 category/status/keyword） |
| GET | `/api/gifts/{id}` | 查询礼物 |
| PUT | `/api/gifts/{id}` | 更新礼物 |
| DELETE | `/api/gifts/{id}` | 删除礼物 |
| POST | `/api/gifts/batch-status` | 批量上下架 |

### 关注 Follow

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/follows` | 关注 |
| DELETE | `/api/follows` | 取关（body 传 follower_id/followee_id） |
| GET | `/api/follows` | 关注列表（筛选 follower_id/followee_id） |
| GET | `/api/follows/{id}` | 查询关注关系 |

### 送礼记录 GiftRecord

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/gift-records` | 送礼（房间需直播中） |
| GET | `/api/gift-records` | 送礼记录列表（筛选 room_id/user_id/gift_id） |
| GET | `/api/gift-records/{id}` | 查询送礼记录 |

### 统计 Stats

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/stats/overview` | 平台总览 |
| GET | `/api/stats/top-streamers` | 主播打赏榜（?limit=N） |
| GET | `/api/stats/room-danmakus` | 房间弹幕统计 |
| GET | `/api/stats/gift-income` | 按天礼物收入 |

## 项目结构

```
origin/
├── cmd/server/main.go        # 入口：配置、装配、优雅关闭
├── internal/
│   ├── app/app.go            # 依赖装配
│   ├── config/config.go      # 环境变量配置
│   ├── model/                # 领域模型 + 状态机 + 校验
│   ├── store/                # Store 接口 + 内存实现
│   ├── service/              # 业务逻辑 + 统计
│   └── handler/              # 路由 + 中间件 + 处理器
└── pkg/                      # httpx / idgen / logger 通用包
```
