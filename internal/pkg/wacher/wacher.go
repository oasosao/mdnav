package wacher

import (
	"os"
	"path"
	"path/filepath"
	"time"

	"mdnav/internal/core"
	"mdnav/internal/pkg/zap"

	"github.com/fsnotify/fsnotify"
)

// WatcherFile 函数用于监听指定目录下的文件变化
func WatcherFile(ctx *core.Context, f func()) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		ctx.Log.Error("创建文件监听器失败", zap.Error(err))
		return
	}
	defer watcher.Close()

	watcherDir := ctx.Conf.GetString("server.content_dir")

	// 递归添加所有子目录到监听列表
	if err := AddWatcherDirRecursive(ctx, watcher, watcherDir); err != nil {
		ctx.Log.Error("添加监听目录失败", zap.String("dir", watcherDir), zap.Error(err))
		return
	}

	// 防抖机制：使用定时器避免短时间内多次触发
	var debounceTimer *time.Timer
	debounceDuration := 500 * time.Millisecond // 500ms防抖

	// 无限循环，监听文件变化
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok { // 如果监听器关闭，则返回
				ctx.Log.Info("文件监听器已关闭")
				return
			}

			// 处理目录创建事件，添加新目录到监听列表
			if event.Op&fsnotify.Create == fsnotify.Create {
				fs, err := os.Stat(event.Name)
				if err == nil && fs.IsDir() {
					if err := AddWatcherDirRecursive(ctx, watcher, event.Name); err != nil {
						ctx.Log.Error("添加新目录到监听列表失败", zap.String("dir", event.Name), zap.Error(err))
					}
				}
			}

			// 只处理Markdown文件的变化
			if path.Ext(event.Name) == ".md" {
				// 取消之前的定时器
				if debounceTimer != nil {
					debounceTimer.Stop()
				}

				// 创建新的定时器
				debounceTimer = time.AfterFunc(debounceDuration, func() {
					ctx.Log.Info("文件变化，触发重新加载", zap.String("file", event.Name), zap.String("op", event.Op.String()))
					f()
				})
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				ctx.Log.Error("文件监听器通道关闭", zap.Error(err))
				return
			}
			ctx.Log.Error("文件监听错误", zap.Error(err))
		}
	}
}

// AddWatcherDirRecursive 递归添加目录及其子目录到监听列表
func AddWatcherDirRecursive(ctx *core.Context, watcher *fsnotify.Watcher, dir string) error {
	// 添加当前目录到监听列表
	if err := watcher.Add(dir); err != nil {
		return err
	}

	ctx.Log.Info("添加监听目录", zap.String("dir", dir))

	// 递归遍历子目录
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && path != dir {
			if err := watcher.Add(path); err != nil {
				ctx.Log.Error("添加子目录到监听列表失败", zap.String("dir", path), zap.Error(err))
				return nil // 继续处理其他目录
			}

			ctx.Log.Info("添加监听子目录", zap.String("dir", path))
		}
		return nil
	})
}
