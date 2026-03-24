package main

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"sync"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"golang.org/x/sys/windows"
)

// 常量定义
const (
	LWA_COLORKEY = 0x00000001
	LWA_ALPHA    = 0x00000002

	GWL_EXSTYLE   = -20
	WS_EX_LAYERED = 0x00080000
)

var (
	user32 = syscall.NewLazyDLL("user32.dll")
	gdi32  = syscall.NewLazyDLL("gdi32.dll")

	getWindowTextW             = user32.NewProc("GetWindowTextW")
	getWindowTextLengthW       = user32.NewProc("GetWindowTextLengthW")
	setLayeredWindowAttributes = user32.NewProc("SetLayeredWindowAttributes")
	fillRect                   = user32.NewProc("FillRect")
	windowFromPoint            = user32.NewProc("WindowFromPoint")

	createSolidBrush = gdi32.NewProc("CreateSolidBrush")
)

// GetWindowText 返回指定窗口句柄的文本
func GetWindowText(hwnd uintptr) (string, error) {
	// 1. 获取文本长度（字符数，不含结尾的 null）
	length, _, _ := getWindowTextLengthW.Call(uintptr(hwnd))
	if length == 0 {
		// 窗口没有文本，返回空字符串
		return "", nil
	}

	// 2. 分配缓冲区（UTF-16 编码，包含结尾 null）
	buf := make([]uint16, length+1)

	// 3. 获取文本
	ret, _, err := getWindowTextW.Call(uintptr(hwnd), uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)))
	if ret == 0 {
		// 调用失败，返回错误
		return "", fmt.Errorf("GetWindowText failed: %v", err)
	}

	// 4. 将 UTF-16 切片转换为 Go 字符串
	return syscall.UTF16ToString(buf), nil
}

// SetLayeredWindowAttributes 设置分层窗口属性
func SetLayeredWindowAttributes(hwnd win.HWND, crKey uint32, bAlpha byte, dwFlags uint32) (bool, error) {
	ret, _, err := setLayeredWindowAttributes.Call(
		uintptr(hwnd),
		uintptr(crKey),
		uintptr(bAlpha),
		uintptr(dwFlags),
	)
	if ret == 0 { // 返回 0 表示失败
		return false, err
	}
	return true, nil
}

// CreateSolidBrush 创建实心画刷
func CreateSolidBrush(crColor win.COLORREF) (win.HGDIOBJ, error) {
	ret, _, err := createSolidBrush.Call(uintptr(crColor))
	if ret == 0 {
		return 0, err
	}
	return win.HGDIOBJ(ret), nil
}

// FillRect 用画刷填充矩形
func FillRect(hdc win.HDC, rect *win.RECT, hbr win.HGDIOBJ) error {
	ret, _, err := fillRect.Call(uintptr(hdc), uintptr(unsafe.Pointer(rect)), uintptr(hbr))
	if ret == 0 {
		return err
	}
	return nil
}

//go:embed frontend/dist/*
var assets embed.FS

func main() {
	app := NewApp()
	err := wails.Run(&options.App{
		Title:  "Window Tools",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
			Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
				operatorAssetHandler{}.ServeHTTP(rw, req)
			}),
		},
		Bind: []any{app},
	})
	if err != nil {
		log.Fatal("应用启动失败:", err)
	}
}

type App struct{}

type WindowInfo struct {
	Hwnd      uintptr `json:"hwnd"`
	Title     string  `json:"title"`
	ClassName string  `json:"className"`
}

// GetWindowInfoResult 获取窗口信息的结果
type GetWindowInfoResult struct {
	Title     string `json:"title"`
	ClassName string `json:"className"`
}

// HighlightWindowResult 高亮窗口的结果
type HighlightWindowResult struct {
	Error string `json:"error,omitempty"`
}

func NewApp() *App {
	return &App{}
}

type MousePosition struct {
	X int `json:"x"`
	Y int `json:"y"`
}

func (a *App) GetMousePosition() (MousePosition, error) {
	var pt win.POINT
	ok := win.GetCursorPos(&pt)
	if !ok {
		return MousePosition{}, errors.New("failed to get cursor position")
	}
	return MousePosition{X: int(pt.X), Y: int(pt.Y)}, nil
}

// GetTopWindows 获取所有顶级窗口
func (a *App) GetTopWindows() ([]WindowInfo, error) {
	var windowList []WindowInfo
	callback := syscall.NewCallback(func(hwnd syscall.Handle, lParam uintptr) uintptr {
		hw := win.HWND(hwnd)
		if win.IsWindowVisible(hw) {
			title, className := getWindowTextAndClass(uintptr(hwnd))
			windowList = append(windowList, WindowInfo{
				Hwnd:      uintptr(hwnd),
				Title:     title,
				ClassName: className,
			})
		}
		return 1
	})
	if err := windows.EnumWindows(callback, unsafe.Pointer(nil)); err != nil {
		return nil, err
	}
	return windowList, nil
}

func (a *App) GetWindowUnderMouse() uintptr {
	var pt win.POINT
	win.GetCursorPos(&pt)
	return uintptr(WindowFromPoint(pt))
}

func (a *App) GetParentWindow(hwnd uintptr) uintptr {
	return uintptr(win.GetParent(win.HWND(hwnd)))
}

func (a *App) GetChildWindows(hwnd uintptr) []uintptr {
	var children []uintptr
	callback := syscall.NewCallback(func(hwndChild syscall.Handle, lParam uintptr) uintptr {
		children = append(children, uintptr(hwndChild))
		return 1
	})
	win.EnumChildWindows(win.HWND(hwnd), callback, 0)
	return children
}

func (a *App) HighlightWindow(hwnd uintptr) error {
	if hwnd == 0 {
		return errors.New("无效窗口句柄")
	}
	var rect win.RECT
	if !win.GetWindowRect(win.HWND(hwnd), &rect) {
		return errors.New("获取窗口位置失败")
	}
	width := rect.Right - rect.Left
	height := rect.Bottom - rect.Top
	if width <= 0 || height <= 0 {
		return errors.New("窗口尺寸无效")
	}
	go createHighlightWindow(rect.Left, rect.Top, width, height)
	return nil
}

func (a *App) GetWindowInfo(hwnd uintptr) (result GetWindowInfoResult, err error) {
	if hwnd == 0 {
		return GetWindowInfoResult{
			Title:     "",
			ClassName: "",
		}, errors.New("无效窗口句柄")
	}
	title, className := getWindowTextAndClass(hwnd)
	return GetWindowInfoResult{
		Title:     title,
		ClassName: className,
	}, nil
}

func WindowFromPoint(point win.POINT) win.HWND {
	val := uint64(point.X) | (uint64(point.Y) << 32)
	ret, _, _ := windowFromPoint.Call(uintptr(val))

	return win.HWND(ret)
}

// ---------- 辅助函数 ----------

func getWindowTextAndClass(hwnd uintptr) (string, string) {
	hw := win.HWND(hwnd)
	// 获取标题
	title, _ := GetWindowText(hwnd)
	// 获取类名
	classBuf := make([]uint16, 256)
	classLen, _ := win.GetClassName(hw, &classBuf[0], int(len(classBuf)))
	if classLen > 0 {
		classBuf = classBuf[:classLen]
	}
	className := syscall.UTF16ToString(classBuf)
	return title, className
}

var highlightMutex sync.Mutex

// createHighlightWindow 创建红色闪烁边框窗口
func createHighlightWindow(x, y, width, height int32) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	instance := win.GetModuleHandle(nil)
	className, _ := syscall.UTF16PtrFromString("HighlightWindowClass")

	// 注册窗口类
	wndClass := win.WNDCLASSEX{
		CbSize:        uint32(unsafe.Sizeof(win.WNDCLASSEX{})),
		Style:         0,
		LpfnWndProc:   syscall.NewCallback(wndProc),
		CbClsExtra:    0,
		CbWndExtra:    0,
		HInstance:     instance,
		HIcon:         0,
		HCursor:       0,
		HbrBackground: 0,
		LpszMenuName:  nil,
		LpszClassName: className,
		HIconSm:       0,
	}
	if win.RegisterClassEx(&wndClass) == 0 {
		// 类可能已注册，忽略错误
	}

	// 创建窗口：分层、穿透、置顶
	hwnd := win.CreateWindowEx(
		win.WS_EX_LAYERED|win.WS_EX_TRANSPARENT|win.WS_EX_TOPMOST,
		className,
		nil,
		win.WS_POPUP,
		x, y, width, height,
		0, 0, instance, nil,
	)
	if hwnd == 0 {
		return
	}
	defer win.DestroyWindow(hwnd)

	// 设置半透明红色（整体透明度 128）
	SetLayeredWindowAttributes(hwnd, 0, 128, LWA_ALPHA)

	// 存储状态到窗口用户数据
	type state struct {
		visible    bool
		flashCount int
	}
	s := &state{visible: true, flashCount: 0}
	win.SetWindowLongPtr(hwnd, win.GWLP_USERDATA, uintptr(unsafe.Pointer(s)))

	win.ShowWindow(hwnd, win.SW_SHOW)
	win.InvalidateRect(hwnd, nil, true)
	win.UpdateWindow(hwnd)

	// 启动定时器（间隔 500ms）
	const timerID = 1
	const interval = 500
	win.SetTimer(hwnd, timerID, uint32(interval), uintptr(0))

	// 消息循环
	var msg win.MSG
	for {
		ret := win.GetMessage(&msg, 0, 0, 0)
		if ret == 0 {
			break
		}
		win.TranslateMessage(&msg)
		win.DispatchMessage(&msg)
	}
}

// wndProc 窗口过程
func wndProc(hwnd win.HWND, msg uint32, wParam, lParam uintptr) uintptr {
	switch msg {
	case win.WM_PAINT:
		var ps win.PAINTSTRUCT
		hdc := win.BeginPaint(hwnd, &ps)
		if hdc != 0 {
			var rect win.RECT
			win.GetClientRect(hwnd, &rect)
			brush, _ := CreateSolidBrush(win.RGB(255, 0, 0))
			defer win.DeleteObject(brush)
			FillRect(hdc, &rect, brush)
			win.EndPaint(hwnd, &ps)
		}
		return 0

	case win.WM_TIMER:
		if wParam == 1 {
			ptr := win.GetWindowLongPtr(hwnd, win.GWLP_USERDATA)
			if ptr == 0 {
				return 0
			}
			//nolint:unsafeptr
			s := (*state)(unsafe.Pointer(ptr))
			if s.flashCount >= 4 { // 闪烁两次（显示-隐藏-显示-隐藏）
				win.KillTimer(hwnd, 1)
				win.DestroyWindow(hwnd)
			} else {
				if s.visible {
					win.ShowWindow(hwnd, win.SW_HIDE)
				} else {
					win.ShowWindow(hwnd, win.SW_SHOW)
				}
				s.visible = !s.visible
				s.flashCount++
			}
		}
		return 0

	case win.WM_DESTROY:
		win.PostQuitMessage(0)
		return 0
	}
	return win.DefWindowProc(hwnd, msg, wParam, lParam)
}

// state 类型定义（必须与 createHighlightWindow 中的一致）
type state struct {
	visible    bool
	flashCount int
}
