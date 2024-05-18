package shine

import (
	"github.com/phuhao00/shine/pkg/log"
	"github.com/phuhao00/shine/servers/game"
	"github.com/phuhao00/shine/servers/game/module"
	"os"
	"os/signal"

	"github.com/phuhao00/shine/conf"
)

func Run(mods ...module.Module) {
	// logger
	if conf.LogLevel != "" {
		logger, err := log.New(conf.LogLevel, conf.LogPath, conf.LogFlag)
		if err != nil {
			panic(err)
		}
		log.Export(logger)
		defer logger.Close()
	}

	log.Release("Leaf %v starting up", version)

	// module
	for i := 0; i < len(mods); i++ {
		module.Register(mods[i])
	}
	module.Init()

	// cluster
	game.Init()

	// close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	log.Release("Leaf closing down (signal: %v)", sig)
	game.Destroy()
	module.Destroy()
}
