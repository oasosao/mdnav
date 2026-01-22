package store

import (
	"os"
	"path"

	"mdnav/internal/conf"
	"mdnav/internal/core"

	"github.com/fsnotify/fsnotify"
)

// WatcherFile 函数用于监听指定目录下的文件变化
func WatcherFile(ctx *core.Context, f func()) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		ctx.Logger.Error("监听目录发生错误" + err.Error())
		return
	}
	defer watcher.Close()

	watcherDir := conf.Config().GetString("server.content_dir")

	AddWatcherDir(ctx, watcher, watcherDir)

	// 无限循环，监听文件变化
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok { // 如果监听器关闭，则返回
				return
			}

			fs, err := os.Stat(event.Name)
			if err == nil {
				if fs.IsDir() {
					AddWatcherDir(ctx, watcher, event.Name)
				}
			}

			if path.Ext(event.Name) == ".md" {
				ctx.Logger.Info("文件:" + event.Name + "有变动")
				f()
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				ctx.Logger.Error("监听目录发生错误" + err.Error())
			}

		}
	}
}

func AddWatcherDir(ctx *core.Context, watcher *fsnotify.Watcher, dir string) {
	if err := watcher.Add(dir); err != nil {
		ctx.Logger.Error("监听目录发生错误" + err.Error())
	}
	ctx.Logger.Info("监听目录:" + dir)
}
