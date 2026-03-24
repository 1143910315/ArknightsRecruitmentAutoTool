## Why

当前本地缓存图片依赖 Wails 资产处理链按 GET 请求返回资源，但你当前的运行约束里该处理入口不支持 GET，只支持 POST/PUT。这意味着现有图片加载方式在前端无法稳定工作，需要改成一套与 Wails 调用模型兼容的图片获取机制。

## What Changes

- 移除对 GET 方式本地缓存图片访问的依赖，改为通过 Wails 可支持的请求方式或绑定调用来获取缓存图片内容。
- 调整后端本地图片暴露接口，使其适用于 POST/PUT 约束下的前端调用。
- 调整前端干员数据页的图片展示逻辑，使其通过新的 Wails 兼容流程加载缓存图片。
- 保持现有程序运行目录缓存、本地优先加载和顺序保持能力不变，只替换图片访问方式。

## Capabilities

### New Capabilities
- `wails-compatible-operator-image-fetch`: 通过 Wails 支持的请求或绑定方式获取本地缓存干员图片内容。
- `operator-image-client-rendering`: 前端基于新图片获取接口渲染本地缓存图片，而不是依赖 GET 资产 URL。

### Modified Capabilities
- None.

## Impact

- 影响 Go 后端本地缓存图片读取接口和返回格式。
- 影响前端干员数据页的图片加载逻辑，可能需要使用 blob、base64 或其他可渲染形式。
- 影响 Wails 资源访问设计，需要避免继续依赖 GET 型资产入口。
- 需要验证在缓存存在时，打包后的应用仍可显示本地图片并保持性能可接受。
